package utils

// ------------------------------------------------------------------------------------------------

import "os"

// ------------------------------------------------------------------------------------------------

// Permite obtener el valor de una variable de ambiente
// o un valorp or defecto, en caso de no hallar la primera.
func GetEnv(key, fallback string) string {
	val := os.Getenv(key)
	if val == "" {
		return fallback
	}
	return val
}

// ------------------------------------------------------------------------------------------------
