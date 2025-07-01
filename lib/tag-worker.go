package cfgldrlib

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/hashicorp/vault/api"
)

func getVaultValue(client *api.Client, mount, secretPath, key string) (any, error) {
	vaultPath := fmt.Sprintf("%s/data/%s", mount, secretPath)
	secret, err := client.Logical().Read(vaultPath)
	if err != nil {
		return nil, err
	}
	if secret == nil || secret.Data == nil {
		return nil, fmt.Errorf("vault: secret not found at path %s", vaultPath)
	}

	data, ok := secret.Data["data"].(map[string]any)
	if !ok {
		return nil, fmt.Errorf("vault: failed to get data from secret %s", vaultPath)
	}

	if key == "" {
		return data, nil
	}
	val, ok := data[key]
	if !ok {
		return nil, fmt.Errorf("vault: key %s not found in secret %s", key, vaultPath)
	}
	return val, nil
}

func fillStructFromVault(v any, client *api.Client, basePath string) error {
	val := reflect.ValueOf(v)
	if val.Kind() != reflect.Ptr || val.IsNil() {
		return fmt.Errorf("expected pointer to struct, got %T", v)
	}
	val = val.Elem()
	typ := val.Type()
	if typ.Kind() != reflect.Struct {
		return fmt.Errorf("expected struct, got %T", v)
	}

	for i, field := range reflect.VisibleFields(typ) {
		fieldVal := val.Field(i)
		cfgldrTag := field.Tag.Get(TagName)
		if cfgldrTag == "-" {
			continue
		}
		vaultKey := field.Name
		if cfgldrTag != "" {
			for part := range strings.SplitSeq(cfgldrTag, ",") {
				if strings.HasPrefix(part, ValTagParam) {
					vaultKey = strings.TrimPrefix(part, ValTagParam)
				}
			}
		}
		if fieldVal.Kind() == reflect.Struct {
			nextPath := vaultKey
			if basePath != "" {
				nextPath = basePath + "/" + vaultKey
			}
			err := fillStructFromVault(fieldVal.Addr().Interface(), client, nextPath)
			if err != nil {
				fmt.Printf("⚠️ Error processing nested struct at path %s: %v\n", nextPath, err)
			}
			continue
		}

		secretPath := basePath
		if secretPath != "" {
			secretPath = secretPath + "/" + vaultKey
		} else {
			secretPath = vaultKey
		}
		valFromVault, err := getVaultValue(client, basePath, secretPath, vaultKey)
		if err != nil {
			fmt.Printf("⚠️ Value not found for key: %s (field: %s) — %v\n", secretPath, vaultKey, err)
			continue
		}
		if fieldVal.CanSet() {
			err = setValue(fieldVal, valFromVault)
			if err != nil {
				fmt.Printf("⚠️ Failed to set value for %s (field: %s): %v\n", secretPath, vaultKey, err)
			}
		}
	}
	return nil
}

func setValue(fieldVal reflect.Value, val any) error {
	switch fieldVal.Kind() {
	case reflect.String:
		str, ok := val.(string)
		if !ok {
			str = fmt.Sprintf("%v", val)
		}
		fieldVal.SetString(str)
	case reflect.Int, reflect.Int64, reflect.Int32, reflect.Int16, reflect.Int8:
		var intVal int64
		switch v := val.(type) {
		case int:
			intVal = int64(v)
		case int64:
			intVal = v
		case float64:
			intVal = int64(v)
		case string:
			var err error
			intVal, err = parseStringToInt64(v)
			if err != nil {
				return err
			}
		default:
			str := fmt.Sprintf("%v", v)
			var err error
			intVal, err = parseStringToInt64(str)
			if err != nil {
				return err
			}
		}
		fieldVal.SetInt(intVal)
	case reflect.Float32, reflect.Float64:
		var floatVal float64
		switch v := val.(type) {
		case float64:
			floatVal = v
		case float32:
			floatVal = float64(v)
		case int:
			floatVal = float64(v)
		case int64:
			floatVal = float64(v)
		case string:
			var err error
			floatVal, err = parseStringToFloat64(v)
			if err != nil {
				return err
			}
		default:
			str := fmt.Sprintf("%v", v)
			var err error
			floatVal, err = parseStringToFloat64(str)
			if err != nil {
				return err
			}
		}
		fieldVal.SetFloat(floatVal)
	case reflect.Bool:
		var boolVal bool
		switch v := val.(type) {
		case bool:
			boolVal = v
		case string:
			var err error
			boolVal, err = parseStringToBool(v)
			if err != nil {
				return err
			}
		default:
			str := fmt.Sprintf("%v", v)
			var err error
			boolVal, err = parseStringToBool(str)
			if err != nil {
				return err
			}
		}
		fieldVal.SetBool(boolVal)
	case reflect.Slice:
		return setSliceValue(fieldVal, val)
	default:
		return fmt.Errorf("unsupported field type: %s", fieldVal.Kind())
	}
	return nil
}

func setSliceValue(fieldVal reflect.Value, val any) error {
	elemType := fieldVal.Type().Elem()

	if slice, ok := val.([]any); ok {
		result := reflect.MakeSlice(fieldVal.Type(), len(slice), len(slice))

		for i, item := range slice {
			elem := reflect.New(elemType).Elem()

			err := setValue(elem, item)
			if err != nil {
				return fmt.Errorf("failed to set slice element %d: %v", i, err)
			}

			result.Index(i).Set(elem)
		}

		fieldVal.Set(result)
		return nil
	}

	if str, ok := val.(string); ok {
		if strings.HasPrefix(strings.TrimSpace(str), "[") && strings.HasSuffix(strings.TrimSpace(str), "]") {
			if elemType.Kind() == reflect.String {
				content := strings.TrimSpace(str[1 : len(str)-1])
				if content == "" {
					fieldVal.Set(reflect.MakeSlice(fieldVal.Type(), 0, 0))
					return nil
				}

				items := strings.Split(content, ",")
				result := reflect.MakeSlice(fieldVal.Type(), len(items), len(items))

				for i, item := range items {
					cleanItem := strings.Trim(strings.TrimSpace(item), `"'`)
					result.Index(i).SetString(cleanItem)
				}

				fieldVal.Set(result)
				return nil
			}
		}

		if elemType.Kind() == reflect.String {
			items := strings.Split(str, ",")
			result := reflect.MakeSlice(fieldVal.Type(), len(items), len(items))

			for i, item := range items {
				cleanItem := strings.TrimSpace(item)
				result.Index(i).SetString(cleanItem)
			}

			fieldVal.Set(result)
			return nil
		}
	}

	elem := reflect.New(elemType).Elem()
	err := setValue(elem, val)
	if err != nil {
		return fmt.Errorf("failed to set single element in slice: %v", err)
	}

	result := reflect.MakeSlice(fieldVal.Type(), 1, 1)
	result.Index(0).Set(elem)
	fieldVal.Set(result)

	return nil
}

func parseStringToInt64(s string) (int64, error) {
	var i int64
	_, err := fmt.Sscanf(s, "%d", &i)
	if err != nil {
		return 0, fmt.Errorf("failed to parse int from string: %s", s)
	}
	return i, nil
}

func parseStringToFloat64(s string) (float64, error) {
	var f float64
	_, err := fmt.Sscanf(s, "%f", &f)
	if err != nil {
		return 0, fmt.Errorf("failed to parse float from string: %s", s)
	}
	return f, nil
}

func parseStringToBool(s string) (bool, error) {
	s = strings.ToLower(strings.TrimSpace(s))
	if s == "true" || s == "1" {
		return true, nil
	}
	if s == "false" || s == "0" {
		return false, nil
	}
	return false, fmt.Errorf("failed to parse bool from string: %s", s)
}
