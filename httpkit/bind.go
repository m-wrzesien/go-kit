package httpkit

import (
	"fmt"
	"net/http"
	"reflect"
	"strconv"
)

const (
	tagForm = "form"
	tagPath = "path"
)

// Bind parses request Path Wildcards and Form Data from FormValue and stores the result in the value pointed by v. Currently only structs are supported.
// Each field in the struct can use one of the tags - "form" or "path" - which indicates it's source.
func Bind(r *http.Request, v any) error {
	tags := []string{tagForm, tagPath}
	err := r.ParseForm()
	if err != nil {
		return fmt.Errorf("error while parsing form in bind: %w", err)
	}

	typ := reflect.TypeOf(v).Elem()
	val := reflect.ValueOf(v).Elem()

	if typ.Kind() != reflect.Struct {
		return nil
	}

	for i := range typ.NumField() {
		typeField := typ.Field(i)
		structField := val.Field(i)
		if typeField.Anonymous {
			if structField.Kind() == reflect.Pointer {
				structField = structField.Elem()
			}
		}
		if !structField.CanSet() {
			continue
		}
		for _, tag := range tags {
			inputFieldName := typeField.Tag.Get(tag)

			if inputFieldName == "" {
				continue
			}

			fv := getBindValue(r, tag, inputFieldName)

			// value is empty so we don't need to parse it - moving to next field
			if fv == "" {
				break
			}

			switch typeField.Type.Kind() {
			case reflect.Float64:
				v, err := strconv.ParseFloat(fv, 64)
				if err != nil {
					return fmt.Errorf("error while parsing %q field as float64 in bind: %w", inputFieldName, err)
				}
				structField.SetFloat(v)
			case reflect.Int:
				v, err := strconv.ParseInt(fv, 10, 0)
				if err != nil {
					return fmt.Errorf("error while parsing %q field as int in bind: %w", inputFieldName, err)
				}
				structField.SetInt(v)
			case reflect.String:
				structField.SetString(fv)
			case reflect.Bool:
				v, err := strconv.ParseBool(fv)
				if err != nil {
					return fmt.Errorf("error while parsing %q field as bool in bind: %w", inputFieldName, err)
				}
				structField.SetBool(v)
			}

			break
		}
	}

	return nil
}

// getBindValue return correct value from *http.Request from Form or Path, based on tag value
func getBindValue(r *http.Request, tag, field string) string {
	switch tag {
	case tagForm:
		return r.FormValue(field)
	case tagPath:
		return r.PathValue(field)
	}
	return ""
}
