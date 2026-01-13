package export

import (
	"database/sql"
	"fmt"
	"reflect"
)

type fieldMeta struct {
	Name  string
	Index int
}

func getFields[T any]() []fieldMeta {
	var t T
	typ := reflect.TypeOf(t)

	fields := make([]fieldMeta, 0, typ.NumField())

	for i := 0; i < typ.NumField(); i++ {
		f := typ.Field(i)

		name := f.Tag.Get("json")
		if name == "" || name == "-" {
			continue
		}

		fields = append(fields, fieldMeta{
			Name:  name,
			Index: i,
		})
	}

	return fields
}

func valueToString(v reflect.Value) string {
	if !v.IsValid() {
		return ""
	}

	switch v.Interface().(type) {

	case sql.NullString:
		ns := v.Interface().(sql.NullString)
		if ns.Valid {
			return ns.String
		}
		return ""

	case sql.NullInt32:
		ni := v.Interface().(sql.NullInt32)
		if ni.Valid {
			return fmt.Sprint(ni.Int32)
		}
		return ""
	}

	return fmt.Sprint(v.Interface())
}
