package sdk

import (
	"fmt"
	"log"
	"reflect"
)

func (rmd *ResourceMetaData) Encode(input interface{}) error {
	objType := reflect.TypeOf(input).Elem()
	objVal := reflect.ValueOf(input).Elem()

	serialized, err := recurse(objType, objVal)
	if err != nil {
		return err
	}

	for k, v := range *serialized {
		if err := rmd.ResourceData.Set(k, v); err != nil {
			return fmt.Errorf("settting %q: %+v", k, err)
		}
	}
	return nil
}

func recurse(objType reflect.Type, objVal reflect.Value) (*map[string]interface{}, error) {
	output := make(map[string]interface{}, 0)
	for i := 0; i < objType.NumField(); i++ {
		field := objType.Field(i)
		fieldVal := objVal.Field(i)
		if hclTag, exists := field.Tag.Lookup("hcl"); exists {
			switch field.Type.Kind() {
			case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
				iv := fieldVal.Int()
				log.Printf("[TOMTOM] Setting %q to %d", hclTag, iv)

				output[hclTag] = iv

			case reflect.Float32, reflect.Float64:
				fv := fieldVal.Float()
				log.Printf("[TOMTOM] Setting %q to %f", hclTag, fv)

				output[hclTag] = fv

			case reflect.String:
				sv := fieldVal.String()
				log.Printf("[TOMTOM] Setting %q to %q", hclTag, sv)
				output[hclTag] = sv

			case reflect.Bool:
				bv := fieldVal.Bool()
				log.Printf("[BOOL] Setting %q to %t", hclTag, bv)
				output[hclTag] = bv

			case reflect.Map:
				iter := fieldVal.MapRange()
				attr := make(map[string]interface{})
				for iter.Next() {
					attr[iter.Key().String()] = iter.Value().Interface()
				}
				output[hclTag] = attr

			case reflect.Slice:
				sv := fieldVal.Slice(0, fieldVal.Len())
				attr := make([]interface{}, sv.Len())
				switch sv.Type() {
				case reflect.TypeOf([]string{}), reflect.TypeOf([]int{}), reflect.TypeOf([]float64{}), reflect.TypeOf([]bool{}):
					log.Printf("[SLICE] Setting %q to %q", hclTag, sv)
					output[hclTag] = sv.Interface()

				default:
					for i := 0; i < sv.Len(); i++ {
						log.Printf("[SLICE] Index %d is %q", i, sv.Index(i).Interface())
						log.Printf("[SLICE] Type %+v", sv.Type())
						nestedType := sv.Index(i).Type()
						nestedValue := sv.Index(i)
						serialized, err := recurse(nestedType, nestedValue)
						if err != nil {
							panic(err)
						}
						attr[i] = serialized
					}
					log.Printf("[SLICE] Setting %q to %+v", hclTag, attr)
					output[hclTag] = attr
				}
			default:
				return &output, fmt.Errorf("unknown type %+v for key %q", field.Type.Kind(), hclTag)
			}
		}
	}

	return &output, nil
}
