package sdk

import (
	"fmt"
	"reflect"
)

func (rmd *ResourceMetaData) Encode(input interface{}) error {
	objType := reflect.TypeOf(input).Elem()
	objVal := reflect.ValueOf(input).Elem()

	serialized, err := recurse(objType, objVal, rmd.serializationDebugLogger)
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

func recurse(objType reflect.Type, objVal reflect.Value, debugLogger Logger) (*map[string]interface{}, error) {
	output := make(map[string]interface{}, 0)
	for i := 0; i < objType.NumField(); i++ {
		field := objType.Field(i)
		fieldVal := objVal.Field(i)
		if hclTag, exists := field.Tag.Lookup("hcl"); exists {
			switch field.Type.Kind() {
			case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
				iv := fieldVal.Int()
				debugLogger.Infof("Setting %q to %d", hclTag, iv)

				output[hclTag] = iv

			case reflect.Float32, reflect.Float64:
				fv := fieldVal.Float()
				debugLogger.Infof("Setting %q to %f", hclTag, fv)

				output[hclTag] = fv

			case reflect.String:
				sv := fieldVal.String()
				debugLogger.Infof("Setting %q to %q", hclTag, sv)
				output[hclTag] = sv

			case reflect.Bool:
				bv := fieldVal.Bool()
				debugLogger.Infof("Setting %q to %t", hclTag, bv)
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
					debugLogger.Infof("Setting %q to %q", hclTag, sv)
					output[hclTag] = sv.Interface()

				default:
					for i := 0; i < sv.Len(); i++ {
						debugLogger.Infof("[SLICE] Index %d is %q", i, sv.Index(i).Interface())
						debugLogger.Infof("[SLICE] Type %+v", sv.Type())
						nestedType := sv.Index(i).Type()
						nestedValue := sv.Index(i)
						serialized, err := recurse(nestedType, nestedValue, debugLogger)
						if err != nil {
							return nil, fmt.Errorf("serializing nested object %q: %+v", sv.Type(), exists)
						}
						attr[i] = serialized
					}
					debugLogger.Infof("[SLICE] Setting %q to %+v", hclTag, attr)
					output[hclTag] = attr
				}
			default:
				return &output, fmt.Errorf("unknown type %+v for key %q", field.Type.Kind(), hclTag)
			}
		}
	}

	return &output, nil
}
