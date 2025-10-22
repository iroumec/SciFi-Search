package utils

import "fmt"

func ToString(v any) string {
	if v == nil {
		return ""
	}
	return fmt.Sprintf("%v", v)
}
