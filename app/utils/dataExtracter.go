package utils

import (
	"bytes"
	"context"
	"fmt"
	"os/exec"
)

func main() {
	// See "man pdftotext" for more options.
	args := []string{
		"-layout",         // Maintain (as best as possible) the original physical layout of the text.
		"-nopgbrk",        // Don't insert page breaks (form feed characters) between pages.
		"Comprobante.pdf", // The input file.
		"-",               // Send the output to stdout.
	}
	cmd := exec.CommandContext(context.Background(), "pdftotext", args...)

	var buf bytes.Buffer
	cmd.Stdout = &buf

	if err := cmd.Run(); err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println(buf.String())
}
