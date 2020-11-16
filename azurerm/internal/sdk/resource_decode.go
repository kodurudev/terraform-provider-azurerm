package sdk

import (
	"reflect"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func (rmd ResourceMetaData) Decode(input interface{}) error {
	objType := reflect.TypeOf(input).Elem()
	for i := 0; i < objType.NumField(); i++ {
		field := objType.Field(i)
		rmd.serializationDebugLogger.Infof("Field", field)

		if val, exists := field.Tag.Lookup("computed"); exists {
			if val == "true" {
				continue
			}
		}

		if val, exists := field.Tag.Lookup("hcl"); exists {
			hclValue := rmd.ResourceData.Get(val)

			rmd.serializationDebugLogger.Infof("HCLValue: ", hclValue)
			rmd.serializationDebugLogger.Infof("Input Type: ", reflect.ValueOf(input).Elem().Field(i).Type())

			if err := setValue(input, hclValue, i, rmd.serializationDebugLogger); err != nil {
				return err
			}
		}
	}
	return nil
}

func setValue(input, hclValue interface{}, index int, debugLogger Logger) error {
	if v, ok := hclValue.(string); ok {
		debugLogger.Infof("[String] Decode %+v", v)
		debugLogger.Infof("Input %+v", reflect.ValueOf(input))
		debugLogger.Infof("Input Elem %+v", reflect.ValueOf(input).Elem())
		reflect.ValueOf(input).Elem().Field(index).SetString(v)
		return nil
	}

	if v, ok := hclValue.(int); ok {
		debugLogger.Infof("[INT] Decode %+v", v)
		reflect.ValueOf(input).Elem().Field(index).SetInt(int64(v))
		return nil
	}

	if v, ok := hclValue.(float64); ok {
		debugLogger.Infof("[Float] Decode %+v", v)
		reflect.ValueOf(input).Elem().Field(index).SetFloat(v)
		return nil
	}

	// Doesn't work for empty bools?
	if v, ok := hclValue.(bool); ok {
		debugLogger.Infof("[BOOL] Decode %+v", v)

		reflect.ValueOf(input).Elem().Field(index).SetBool(v)
		return nil
	}

	if v, ok := hclValue.(*schema.Set); ok {
		setListValue(input, index, v.List(), debugLogger)
		return nil
	}

	if mapConfig, ok := hclValue.(map[string]interface{}); ok {
		mapOutput := reflect.MakeMap(reflect.TypeOf(map[string]string{}))
		for key, val := range mapConfig {
			mapOutput.SetMapIndex(reflect.ValueOf(key), reflect.ValueOf(val))
		}

		reflect.ValueOf(input).Elem().Field(index).Set(mapOutput)
		return nil
	}

	if v, ok := hclValue.([]interface{}); ok {
		setListValue(input, index, v, debugLogger)
		return nil
	}

	return nil
}

func setListValue(input interface{}, index int, v []interface{}, debugLogger Logger) {
	switch fieldType := reflect.ValueOf(input).Elem().Field(index).Type(); fieldType {
	// TODO do I have to do it this way for the rest of the types?
	case reflect.TypeOf([]string{}):
		stringSlice := reflect.MakeSlice(reflect.TypeOf([]string{}), len(v), len(v))
		for i, stringVal := range v {
			stringSlice.Index(i).SetString(stringVal.(string))
		}
		reflect.ValueOf(input).Elem().Field(index).Set(stringSlice)

	case reflect.TypeOf([]int{}):
		iSlice := reflect.MakeSlice(reflect.TypeOf([]int{}), len(v), len(v))
		for i, iVal := range v {
			iSlice.Index(i).SetInt(int64(iVal.(int)))
		}
		reflect.ValueOf(input).Elem().Field(index).Set(iSlice)

	case reflect.TypeOf([]float64{}):
		fSlice := reflect.MakeSlice(reflect.TypeOf([]float64{}), len(v), len(v))
		for i, fVal := range v {
			fSlice.Index(i).SetFloat(fVal.(float64))
		}
		reflect.ValueOf(input).Elem().Field(index).Set(fSlice)

	case reflect.TypeOf([]bool{}):
		bSlice := reflect.MakeSlice(reflect.TypeOf([]bool{}), len(v), len(v))
		for i, bVal := range v {
			bSlice.Index(i).SetBool(bVal.(bool))
		}
		reflect.ValueOf(input).Elem().Field(index).Set(bSlice)

	default:
		valueToSet := reflect.New(reflect.ValueOf(input).Elem().Field(index).Type())
		debugLogger.Infof("List Type", valueToSet.Type())

		for _, mapVal := range v {
			if test, ok := mapVal.(map[string]interface{}); ok && test != nil {
				elem := reflect.New(fieldType.Elem())
				debugLogger.Infof("element ", elem)
				for j := 0; j < elem.Type().Elem().NumField(); j++ {
					nestedField := elem.Type().Elem().Field(j)
					debugLogger.Infof("nestedField ", nestedField)
					if val, exists := nestedField.Tag.Lookup("computed"); exists {
						if val == "true" {
							continue
						}
					}

					if val, exists := nestedField.Tag.Lookup("hcl"); exists {
						nestedHCLValue := test[val]
						setValue(elem.Interface(), nestedHCLValue, j, debugLogger)
					}
				}

				if !elem.CanSet() {
					elem = elem.Elem()
				}

				if valueToSet.Kind() == reflect.Ptr {
					valueToSet.Elem().Set(reflect.Append(valueToSet.Elem(), elem))
				} else {
					valueToSet = reflect.Append(valueToSet, elem)
				}

				debugLogger.Infof("value to set type after changes", valueToSet.Type())
			}
		}
		fieldToSet := reflect.ValueOf(input).Elem().Field(index)

		if valueToSet.Kind() != reflect.Ptr {
			fieldToSet.Set(valueToSet)
		} else {
			fieldToSet.Set(valueToSet.Elem())
		}
	}
}
