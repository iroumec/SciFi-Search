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

func ValidarConstancia(pdf multipart.File) bool {

	texto := leerPDF(pdf)

	// Se normalizan espacios y mayúsculas para facilitar el matching.
	textoNormalizado := strings.ReplaceAll(texto, "\n", " ")
	textoNormalizado = strings.TrimSpace(textoNormalizado)

	dni := hallarDNI(textoNormalizado)
	fmt.Println(dni)
	codigo := hallarCodigoValidacion(textoNormalizado)
	fmt.Println(codigo)

	fmt.Printf("DNI   : %s\n", dni)
	fmt.Printf("Código: %s\n", codigo)

	valido := realizarValidacion(dni, codigo)

	fmt.Println("Mensaje:")
	if valido {
		fmt.Println("Certificado válido")
		return true
	} else {
		fmt.Println("Certificado inválido")
		return false
	}
}

func leerPDF(file multipart.File) string {

	// Se lee todo el contenido del archivo
	data, err := io.ReadAll(file)
	if err != nil {
		// Error al leer el archivo.
		return ""
	}

	cmd := exec.Command("pdftotext", "-", "-") // stdin -> stdout
	cmd.Stdin = bytes.NewReader(data)
	out, err := cmd.Output()
	if err != nil {
		panic(err)
	}
	return string(out)
}

func hallarDNI(text string) string {
	// DNI en Argentina suele ser 7 u 8 dígitos; puede venir con puntos.
	// Buscamos 8 o 7 dígitos o formato con puntos: 12.345.678
	re := regexp.MustCompile(`\b(\d{1,2}\.\d{3}\.\d{3}|\d{7,8})\b`)
	m := re.FindStringSubmatch(text)
	if len(m) > 0 {
		return m[0]
	}
	return ""
}

func hallarCodigoValidacion(text string) string {
	// Buscamos la línea que contiene "CÓDIGO DE VALIDACIÓN" y se captura lo que sigue.
	re := regexp.MustCompile(`(?i)CÓDIGO\s+DE\s+VALIDA(?:C|C)I[ÓO]N[:\s\-]*([A-Z0-9\-]+)`)
	m := re.FindStringSubmatch(text)
	if len(m) >= 2 {
		return m[1]
	}
	// alternativa: buscar cualquier token alfanumérico de 4-12 chars cercano a la palabra CÓDIGO
	re2 := regexp.MustCompile(`(?i)CÓDIGO[\s\S]{0,40}([A-Z0-9]{4,12})`)
	m2 := re2.FindStringSubmatch(text)
	if len(m2) >= 2 {
		return m2[1]
	}
	return ""
}

func realizarValidacion(dni, codigo string) bool {

	// Se crea un contexto de Chrome.
	ctx, cancel := chromedp.NewContext(context.Background())
	defer cancel()

	// Timeout general.
	ctx, cancel = context.WithTimeout(ctx, 15*time.Second)
	defer cancel()

	var resultado string

	// Se ejecutan las tareas.
	err := chromedp.Run(ctx,
		// 1. Abrir la página del validador.
		chromedp.Navigate(`https://guarani.unicen.edu.ar/autogestion/exactas/validador_certificados`),

		// 2. Rellenar los campos del formulario.
		chromedp.SetValue(`#documento`, dni),
		chromedp.SetValue(`#codigo_valid`, codigo),

		// 3. Hacer click en el botón Validar.
		chromedp.Click(`#validar`, chromedp.NodeVisible),

		// 4. Esperar a que aparezca el resultado y obtener el texto.
		chromedp.Text(`div.hero-unit h1`, &resultado, chromedp.NodeVisible, chromedp.ByQuery),
	)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Resultado:", resultado)
	if resultado == "Certificado Válido" || resultado == "Certificado Valido" {
		fmt.Println("Certificado válido")
		return true
	} else {
		fmt.Println("Certificado inválido")
		return false
	}
}
