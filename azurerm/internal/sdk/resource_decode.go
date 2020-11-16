package sdk

import (
	"log"
	"reflect"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func (rmd ResourceMetaData) Decode(input interface{}) error {
	objType := reflect.TypeOf(input).Elem()
	for i := 0; i < objType.NumField(); i++ {
		field := objType.Field(i)
		log.Print("[MATTHEWMATTHEW] Field", field)

		if val, exists := field.Tag.Lookup("computed"); exists {
			if val == "true" {
				continue
			}
		}

		if val, exists := field.Tag.Lookup("hcl"); exists {
			hclValue := rmd.ResourceData.Get(val)

			log.Print("[MATTHEWMATTHEW] HCLValue: ", hclValue)
			log.Print("[MATTHEWMATTHEW] Input Type: ", reflect.ValueOf(input).Elem().Field(i).Type())

			if err := setValue(input, hclValue, i); err != nil {
				return err
			}
		}
	}
	return nil
}

func setValue(input, hclValue interface{}, index int) error {
	if v, ok := hclValue.(string); ok {
		log.Printf("[String] Decode %+v", v)
		log.Printf("Input %+v", reflect.ValueOf(input))
		log.Printf("Input Elem %+v", reflect.ValueOf(input).Elem())
		reflect.ValueOf(input).Elem().Field(index).SetString(v)
		return nil
	}

	if v, ok := hclValue.(int); ok {
		log.Printf("[INT] Decode %+v", v)
		reflect.ValueOf(input).Elem().Field(index).SetInt(int64(v))
		return nil
	}

	if v, ok := hclValue.(float64); ok {
		log.Printf("[Float] Decode %+v", v)
		reflect.ValueOf(input).Elem().Field(index).SetFloat(v)
		return nil
	}

	// Doesn't work for empty bools?
	if v, ok := hclValue.(bool); ok {
		log.Printf("[BOOL] Decode %+v", v)

		reflect.ValueOf(input).Elem().Field(index).SetBool(v)
		return nil
	}

	if v, ok := hclValue.(*schema.Set); ok {
		setListValue(input, index, v.List())
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
		setListValue(input, index, v)
		return nil
	}

	return nil
}

func setListValue(input interface{}, index int, v []interface{}) {
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
		log.Print("[MATTHEWMATTHEW] List Type", valueToSet.Type())

		for _, mapVal := range v {
			if test, ok := mapVal.(map[string]interface{}); ok && test != nil {
				elem := reflect.New(fieldType.Elem())
				log.Print("[MATTHEWMATTHEW] element ", elem)
				for j := 0; j < elem.Type().Elem().NumField(); j++ {
					nestedField := elem.Type().Elem().Field(j)
					log.Print("[MATTHEWMATTHEW] nestedField ", nestedField)
					if val, exists := nestedField.Tag.Lookup("computed"); exists {
						if val == "true" {
							continue
						}
					}

					if val, exists := nestedField.Tag.Lookup("hcl"); exists {
						nestedHCLValue := test[val]
						setValue(elem.Interface(), nestedHCLValue, j)
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

				log.Print("value to set type after changes", valueToSet.Type())
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
