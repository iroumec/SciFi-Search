package utils

import (
	"fmt"
	"os/exec"
	"regexp"
	"strings"
)

func ValidarConstancia() {

	texto := leerPDF("static/Comprobante.pdf")

	// Se normalizan espacios y mayúsculas para facilitar el matching.
	textoNormalizado := strings.ReplaceAll(texto, "\n", " ")
	textoNormalizado = strings.TrimSpace(textoNormalizado)

	fmt.Println("Entré")

	dni := hallarDNI(textoNormalizado)
	fmt.Println(dni)
	codigo := hallarCodigoValidacion(textoNormalizado)
	fmt.Println(codigo)

	fmt.Printf("DNI   : %s\n", dni)
	fmt.Printf("Código: %s\n", codigo)

	if dni != "" && codigo != "" {
		valURL := "https://guarani.unicen.edu.ar/autogestion/exactas/validador_certificados"
		fmt.Printf("\nValidador: %s\n", valURL)
		fmt.Println("Podés usar DNI y Código para completar el formulario del validador.")
	}
}

func leerPDF(path string) string {
	out, err := exec.Command("pdftotext", path, "-").Output()
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
