package dotenv

import (
	"fmt"
	"os"
	"reflect"
	"strconv"
)

func LoadStruct(data interface{}) error {
	dataType := reflect.TypeOf(data)
	dataValue := reflect.ValueOf(data)

	if dataType.Kind() != reflect.Ptr || dataValue.IsNil() {
		return fmt.Errorf("data must be a pointer to a struct")
	}

	dataType = dataType.Elem()
	dataValue = dataValue.Elem()

	return parseFields(dataType, dataValue)
}

func parseFields(dataType reflect.Type, dataValue reflect.Value) error {
	for i := 0; i < dataType.NumField(); i++ {
		field := dataType.Field(i)
		value := dataValue.Field(i)

		if value.Kind() == reflect.Struct {
			if err := parseFields(field.Type, value); err != nil {
				return err
			}
			continue
		}

		envTag := field.Tag.Get("env")
		if envTag == "" {
			continue
		}

		envValue, found := os.LookupEnv(envTag)
		if !found {
			continue
		}

		switch value.Kind() {
		case reflect.String:
			value.SetString(envValue)
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			intVal, err := strconv.ParseInt(envValue, 10, 64)
			if err != nil {
				return fmt.Errorf("failed to parse int field %s: %s", field.Name, err)
			}
			value.SetInt(intVal)
		case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
			uintVal, err := strconv.ParseUint(envValue, 10, 64)
			if err != nil {
				return fmt.Errorf("failed to parse uint field %s: %s", field.Name, err)
			}
			value.SetUint(uintVal)
		case reflect.Bool:
			boolVal, err := strconv.ParseBool(envValue)
			if err != nil {
				return fmt.Errorf("failed to parse bool field %s: %s", field.Name, err)
			}
			value.SetBool(boolVal)
		case reflect.Float32, reflect.Float64:
			floatVal, err := strconv.ParseFloat(envValue, 64)
			if err != nil {
				return fmt.Errorf("failed to parse float field %s: %s", field.Name, err)
			}
			value.SetFloat(floatVal)
		default:
			return fmt.Errorf("unsupported type for field %s", field.Name)
		}
	}

	return nil
}
