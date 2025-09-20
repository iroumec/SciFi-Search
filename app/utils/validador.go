package utils

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"log"
	"mime/multipart"
	"os/exec"
	"regexp"
	"strings"
	"time"

	"github.com/chromedp/chromedp"
)

func ValidarConstancia(pdf multipart.File) (bool, error) {

	texto, err := leerPDF(pdf)
	if err != nil {
		return false, err
	}

	// Se normalizan espacios y mayúsculas para facilitar el matching.
	textoNormalizado := strings.ReplaceAll(texto, "\n", " ")
	textoNormalizado = strings.TrimSpace(textoNormalizado)

	// Se buscan el DNI y el codigo de validación en el texto.
	dni := hallarDNI(textoNormalizado)
	codigo := hallarCodigoValidacion(textoNormalizado)

	fmt.Printf("DNI   : %s\n", dni)
	fmt.Printf("Código: %s\n", codigo)

	// Con los datos obtenidos, se realiza la validación en la página.
	valido := realizarValidacion(dni, codigo)

	if valido {
		fmt.Println("Certificado válido")
	} else {
		fmt.Println("Certificado inválido")
	}

	return valido, nil
}

// ------------------------------------------------------------------------------------------------
// Lectura del PDF
// ------------------------------------------------------------------------------------------------

func leerPDF(file multipart.File) (string, error) {

	// Lectura del contendio del archivo.
	data, err := io.ReadAll(file)
	if err != nil {
		return "", fmt.Errorf("Error leyendo archivo: %w", err)
	}

	// Se mira la cabecera para verificar que le archvio sea un PDF.
	if !bytes.HasPrefix(data, []byte("%PDF-")) {
		return "", fmt.Errorf("El archivo no es un PDF válido.")
	}

	// Se convierte el PDF a texto, para su posterior análisis, usando
	// la herramienta "pdftotext".
	cmd := exec.Command("pdftotext", "-", "-")

	// Se setea como stdin el contenido del archivo.
	cmd.Stdin = bytes.NewReader(data)

	// Se setea como stdout una variable "out".
	out, err := cmd.Output()
	if err != nil {
		return "", fmt.Errorf("error ejecutando pdftotext: %w", err)
	}

	// Se retorna un string con el contenido del PDF.
	return string(out), nil
}

// ------------------------------------------------------------------------------------------------
// Extracción de datos de texto
// ------------------------------------------------------------------------------------------------

func hallarDNI(text string) string {

	// El DNI en Argentina tiene entre 7 u 8 dígitos y puede venir con puntos.
	// Se buscan entonces 8 o 7 dígitos o formato con puntos.

	// Se define la expresión regular.
	re := regexp.MustCompile(`\b(\d{1,2}\.\d{3}\.\d{3}|\d{7,8})\b`)

	// Se buscan matchings.
	m := re.FindStringSubmatch(text)

	// Si se halló al menos un match, se retorna el primero hallado.
	if len(m) > 0 {
		return m[0]
	}

	// Cadena vacía en caso de no haber hallado nada.
	return ""
}

func hallarCodigoValidacion(text string) string {

	// Se busca una línea que contenga "CÓDIGO DE VALIDACIÓN" y se captura lo que sigue.
	// Se define, para ello, una expresión regular.
	// (?i) -> Hace que no distinga entre minúsculas/mayúsculas.
	// \s+ -> Indica al menos un espacio en blanco.
	// [ÓO] -> Indica "O" con tilde o sin tilde.
	// [:\s\-]* -> Dos puntos, espacios o guiones repetidos cero o más veces.
	// ([0-9])+ -> Secuencia de al menos un número entre 0 y 9.
	re := regexp.MustCompile(`(?i)CÓDIGO\s+DE\s+VALIDACI[ÓO]N[:\s\-]*([0-9]+)`)

	// Se busca algún matching.
	// (m) es un slice de strings con esta estructura:
	// m[0] → el match completo (toda la cadena que coincidió con el regex).
	// m[1] → el contenido del primer grupo de captura (...).
	// m[2] → el contenido del segundo grupo, y así sucesivamente.
	// En la expresión regular definida arriba, solo hay un grupo de captura: ([0-9]+).
	// Los grupos de captura van entre paréntesis.
	m := re.FindStringSubmatch(text)
	if len(m) >= 2 {
		fmt.Println("Primer regex")
		return m[1]
	}

	// Cadena vacía en caso de no haber hallado nada.
	return ""
}

// ------------------------------------------------------------------------------------------------
// Obtención del Resultado de la Validación
// ------------------------------------------------------------------------------------------------

func realizarValidacion(dni, codigo string) bool {

	// Configuración del ExecAllocator con flags para Docker.
	// chromedp.DefaultExecAllocatorOptions[:] --> toma las opciones por defecto de chromedp para iniciar Chromium.
	// chromedp.ExecPath() --> indica dónde está el ejecutable de Chromium (en nuestro caso, en Alpine).
	// chromedp.Flag("no-sandbox", true) --> desactiva el sandbox de Chromium. Es necesario porque el contenedor se corre en modo usuario (sin privilegios).
	// chromedp.Flag("disable-gpu", true) → desactiva aceleración por GPU. Mejora el rendimiento ya que no hay entorno gráfico.
	opts := append(chromedp.DefaultExecAllocatorOptions[:],
		chromedp.ExecPath("/usr/bin/chromium-browser"), // o "chromium" según tu apk
		chromedp.Flag("no-sandbox", true),
		chromedp.Flag("disable-gpu", true),
	)

	// Inicialización del "motor" que ejecuta Chronium.
	// Primer contexto: allocator con opciones
	allocCtx, cancel := chromedp.NewExecAllocator(context.Background(), opts...)
	defer cancel()

	// Segundo contexto: Chrome real (usando allocCtx, no uno nuevo)
	ctx, cancel := chromedp.NewContext(allocCtx)
	defer cancel()

	// Se limita toda la operación a 15 segundos.
	// De esta forma, sse evita que la ejecución quede colgada indefinidamente.
	ctx, cancel = context.WithTimeout(ctx, 15*time.Second)

	// Los recursos serán liberados independientemente de lo que pase.
	defer cancel()

	// Variable en la que se almacenará el resultado.
	var resultado string

	// Ejecución de tareas.
	err := chromedp.Run(ctx,

		// 1. Se abre la página del validador.
		chromedp.Navigate(`https://guarani.unicen.edu.ar/autogestion/exactas/validador_certificados`),

		// 2. Se rellenan los campos del formulario con los parámetros.
		chromedp.SetValue(`#documento`, dni),
		chromedp.SetValue(`#codigo_valid`, codigo),

		// 3. Se hace click en el botón "Validar".
		chromedp.Click(`#validar`, chromedp.NodeVisible),

		// 4. Se espera a que aparezca el resultado y se obtiene el texto.
		chromedp.Text(`div.hero-unit h1`, &resultado, chromedp.NodeVisible, chromedp.ByQuery),
	)
	if err != nil {
		log.Fatal(err)
	}

	return resultado == "Certificado Válido" || resultado == "Certificado Valido"
}
