package utils

// ------------------------------------------------------------------------------------------------

import "os"

// ------------------------------------------------------------------------------------------------

/*
Permite obtener una variable de ambiente o
un valor por defecto, en caso de no hallar la primera.
*/
func GetEnv(key, fallback string) string {
	val := os.Getenv(key)
	if val == "" {
		return fallback
	}
	return val
}

// ------------------------------------------------------------------------------------------------
