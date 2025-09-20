package utils

import (
	"fmt"
	"time"
)

func countdownHandler() {
	// Define la fecha y hora de finalización
	finalDate := time.Date(2025, time.December, 31, 23, 59, 59, 0, time.UTC)

	// Bucle para la cuenta regresiva
	for time.Now().Before(finalDate) {
		// Calcula la diferencia de tiempo
		remaining := finalDate.Sub(time.Now())

		// Extrae días, horas, minutos y segundos
		days := int(remaining.Hours() / 24)
		hours := int(remaining.Hours()) % 24
		minutes := int(remaining.Minutes()) % 60
		seconds := int(remaining.Seconds()) % 60

		// Imprime el tiempo restante
		fmt.Printf("Días: %d, Horas: %d, Minutos: %d, Segundos: %d\r", days, hours, minutes, seconds)

		// Pausa por un segundo
		time.Sleep(1 * time.Second)
	}

	fmt.Println("\n¡La cuenta regresiva ha terminado!")
}
