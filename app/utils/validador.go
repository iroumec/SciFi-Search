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

	dni := hallarDNI(textoNormalizado)
	codigo := hallarCodigoValidacion(textoNormalizado)

	fmt.Printf("DNI   : %s\n", dni)
	fmt.Printf("Código: %s\n", codigo)

	valido := realizarValidacion(dni, codigo)

	if valido {
		fmt.Println("Certificado válido")
	} else {
		fmt.Println("Certificado inválido")
	}

	return valido, nil
}

func leerPDF(file multipart.File) (string, error) {
	// Se lee todo el contenido del archivo
	data, err := io.ReadAll(file)
	if err != nil {
		return "", fmt.Errorf("Error leyendo archivo: %w", err)
	}

	// Verificación de cabecera PDF
	if !bytes.HasPrefix(data, []byte("%PDF-")) {
		return "", fmt.Errorf("El archivo no es un PDF válido.")
	}

	// Convertir a texto usando pdftotext.
	cmd := exec.Command("pdftotext", "-", "-")
	cmd.Stdin = bytes.NewReader(data)
	out, err := cmd.Output()
	if err != nil {
		return "", fmt.Errorf("error ejecutando pdftotext: %w", err)
	}

	return string(out), nil
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

	// Configuración del ExecAllocator con flags para Docker
	opts := append(chromedp.DefaultExecAllocatorOptions[:],
		chromedp.ExecPath("/usr/bin/chromium-browser"), // o "chromium" según tu apk
		chromedp.Flag("no-sandbox", true),
		chromedp.Flag("disable-gpu", true),
	)

	// Primer contexto: allocator con opciones
	allocCtx, cancel := chromedp.NewExecAllocator(context.Background(), opts...)
	defer cancel()

	// Segundo contexto: Chrome real (usando allocCtx, no uno nuevo)
	ctx, cancel := chromedp.NewContext(allocCtx)
	defer cancel()

	// Timeout general
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

	return resultado == "Certificado Válido" || resultado == "Certificado Valido"
}
