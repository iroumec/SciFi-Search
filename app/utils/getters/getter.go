package getters

// ------------------------------------------------------------------------------------------------
// Servicios
// ------------------------------------------------------------------------------------------------

// Retorna el valor, si no está vacío. Si no, el valor por defecto.
func GetOrDefault(value, defaultValue string) string {
	if value != "" {
		return value
	}
	return defaultValue
}

// ------------------------------------------------------------------------------------------------
