package sdk

import (
	"reflect"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestEncode(t *testing.T) {
	testCases := []struct {
		Name        string
		Input       *ExampleObj
		Expected    map[string]interface{}
		ExpectError bool
	}{
		{
			Name: "top level - name",
			Input: &ExampleObj{
				Name: "bingo bango",
			},
			Expected: map[string]interface{}{
				"name":         "bingo bango",
				"enabled":      false,
				"float":        float64(0),
				"list":         []interface{}{},
				"networks":     []string(nil),
				"networks_set": []string(nil),
				"int_list":     []int(nil),
				"int_set":      []int(nil),
				"float_list":   []float64(nil),
				"float_set":    []float64(nil),
				"bool_list":    []bool(nil),
				"bool_set":     []bool(nil),
				"number":       int64(0),
				"output":       "",
				"set":          []interface{}{},
				"map":          map[string]interface{}{},
			},
			ExpectError: false,
		},
		{
			Name: "top level - everything",
			Input: &ExampleObj{
				Name:    "bingo bango",
				Float:   123.4,
				Number:  123,
				Enabled: false,
			},
			Expected: map[string]interface{}{
				"name":         "bingo bango",
				"enabled":      false,
				"float":        123.4,
				"list":         []interface{}{},
				"networks":     []string(nil),
				"networks_set": []string(nil),
				"int_list":     []int(nil),
				"int_set":      []int(nil),
				"float_list":   []float64(nil),
				"float_set":    []float64(nil),
				"bool_list":    []bool(nil),
				"bool_set":     []bool(nil),
				"number":       int64(123),
				"output":       "",
				"set":          []interface{}{},
				"map":          map[string]interface{}{},
			},
			ExpectError: false,
		},
		{
			Name: "top level - list",
			Input: &ExampleObj{
				Name:    "bingo bango",
				Float:   123.4,
				Number:  123,
				Enabled: false,
				List: []NetworkList{{
					Name: "first",
				}},
			},
			Expected: map[string]interface{}{
				"name":    "bingo bango",
				"enabled": false,
				"float":   123.4,
				"list": []interface{}{
					&map[string]interface{}{
						"name":  "first",
						"inner": []interface{}{},
					},
				},
				"networks":     []string(nil),
				"networks_set": []string(nil),
				"int_list":     []int(nil),
				"int_set":      []int(nil),
				"float_list":   []float64(nil),
				"float_set":    []float64(nil),
				"bool_list":    []bool(nil),
				"bool_set":     []bool(nil),
				"number":       int64(123),
				"output":       "",
				"set":          []interface{}{},
				"map":          map[string]interface{}{},
			},
			ExpectError: false,
		},
		{
			Name: "top level - list in lists",
			Input: &ExampleObj{
				Name:    "bingo bango",
				Float:   123.4,
				Number:  123,
				Enabled: false,
				List: []NetworkList{{
					Name: "first",
					Inner: []NetworkInner{{
						Name: "get-a-mac",
					}},
				}},
			},
			Expected: map[string]interface{}{
				"name":    "bingo bango",
				"enabled": false,
				"float":   123.4,
				"list": []interface{}{
					&map[string]interface{}{
						"name": "first",
						"inner": []interface{}{
							&map[string]interface{}{
								"name":  "get-a-mac",
								"inner": []interface{}{},
								"set":   []interface{}{},
							},
						},
					},
				},
				"networks":     []string(nil),
				"networks_set": []string(nil),
				"int_list":     []int(nil),
				"int_set":      []int(nil),
				"float_list":   []float64(nil),
				"float_set":    []float64(nil),
				"bool_list":    []bool(nil),
				"bool_set":     []bool(nil),
				"number":       int64(123),
				"output":       "",
				"set":          []interface{}{},
				"map":          map[string]interface{}{},
			},
			ExpectError: false,
		},
		{
			Name: "top level - everything",
			Input: &ExampleObj{
				Name:        "bingo bango",
				Float:       123.4,
				Number:      123,
				Enabled:     false,
				Networks:    []string{"network1", "network2", "network3"},
				NetworksSet: []string{"networkset1", "networkset2", "networkset3"},
				List: []NetworkList{{
					Name: "first",
					Inner: []NetworkInner{{
						Name: "get-a-mac",
					}},
				}},
				Set: []NetworkSet{{
					Name: "setname",
				}},
			},
			Expected: map[string]interface{}{
				"name":    "bingo bango",
				"enabled": false,
				"float":   123.4,
				"list": []interface{}{
					&map[string]interface{}{
						"name": "first",
						"inner": []interface{}{
							&map[string]interface{}{
								"name":  "get-a-mac",
								"inner": []interface{}{},
								"set":   []interface{}{},
							},
						},
					},
				},
				"networks":     []string{"network1", "network2", "network3"},
				"networks_set": []string{"networkset1", "networkset2", "networkset3"},
				"int_list":     []int(nil),
				"int_set":      []int(nil),
				"float_list":   []float64(nil),
				"float_set":    []float64(nil),
				"bool_list":    []bool(nil),
				"bool_set":     []bool(nil),
				"number":       int64(123),
				"output":       "",
				"map":          map[string]interface{}{},
				"set": []interface{}{
					&map[string]interface{}{
						"name":  "setname",
						"inner": []interface{}{},
					},
				},
			},
			ExpectError: false,
		},
		{
			Name: "nests",
			Input: &ExampleObj{
				Name:    "bingo bango",
				Float:   123.4,
				Number:  123,
				Enabled: false,
				List: []NetworkList{
					{
						Name: "first",
						Inner: []NetworkInner{{
							Name: "get-a-mac",
							Inner: []InnerInner{{
								Name:         "innerinner",
								ShouldBeFine: true,
							}},
							Set: []NetworkListSet{{
								Name: "nestedsetname",
							}},
						}},
					},
					{
						Name: "second",
						Inner: []NetworkInner{{
							Name: "get-a-mac2",
							Inner: []InnerInner{{
								Name:         "innerinner2",
								ShouldBeFine: true,
							}},
							Set: []NetworkListSet{{
								Name: "nestedsetname2",
							}},
						}},
					},
				},
				Set: []NetworkSet{{
					Name: "setname",
				}},
			},
			Expected: map[string]interface{}{
				"name":    "bingo bango",
				"enabled": false,
				"float":   123.4,
				"list": []interface{}{
					&map[string]interface{}{
						"name": "first",
						"inner": []interface{}{
							&map[string]interface{}{
								"name": "get-a-mac",
								"inner": []interface{}{
									&map[string]interface{}{
										"name":           "innerinner",
										"should_be_fine": true,
									},
								},
								"set": []interface{}{
									&map[string]interface{}{
										"name": "nestedsetname",
									},
								},
							},
						},
					},
					&map[string]interface{}{
						"name": "second",
						"inner": []interface{}{
							&map[string]interface{}{
								"name": "get-a-mac2",
								"inner": []interface{}{
									&map[string]interface{}{
										"name":           "innerinner2",
										"should_be_fine": true,
									},
								},
								"set": []interface{}{
									&map[string]interface{}{
										"name": "nestedsetname2",
									},
								},
							},
						},
					},
				},

				"networks":     []string(nil),
				"networks_set": []string(nil),
				"int_list":     []int(nil),
				"int_set":      []int(nil),
				"float_list":   []float64(nil),
				"float_set":    []float64(nil),
				"bool_list":    []bool(nil),
				"bool_set":     []bool(nil),
				"number":       int64(123),
				"output":       "",
				"map":          map[string]interface{}{},
				"set": []interface{}{
					&map[string]interface{}{
						"name":  "setname",
						"inner": []interface{}{},
					},
				},
			},
			ExpectError: false,
		},
		{
			Name: "top level - int lists/sets",
			Input: &ExampleObj{
				IntList: []int{1, 2, 3},
				IntSet:  []int{3, 4, 5},
			},
			Expected: map[string]interface{}{
				"name":         "",
				"enabled":      false,
				"float":        float64(0),
				"list":         []interface{}{},
				"networks":     []string(nil),
				"networks_set": []string(nil),
				"int_list":     []int{1, 2, 3},
				"int_set":      []int{3, 4, 5},
				"float_list":   []float64(nil),
				"float_set":    []float64(nil),
				"bool_list":    []bool(nil),
				"bool_set":     []bool(nil),
				"number":       int64(0),
				"output":       "",
				"set":          []interface{}{},
				"map":          map[string]interface{}{},
			},
			ExpectError: false,
		},
		{
			Name: "top level - float lists/sets",
			Input: &ExampleObj{
				FloatList: []float64{1, 2, 3},
				FloatSet:  []float64{3, 4, 5},
			},
			Expected: map[string]interface{}{
				"name":         "",
				"enabled":      false,
				"float":        float64(0),
				"list":         []interface{}{},
				"networks":     []string(nil),
				"networks_set": []string(nil),
				"int_list":     []int(nil),
				"int_set":      []int(nil),
				"float_list":   []float64{1, 2, 3},
				"float_set":    []float64{3, 4, 5},
				"bool_list":    []bool(nil),
				"bool_set":     []bool(nil),
				"number":       int64(0),
				"output":       "",
				"set":          []interface{}{},
				"map":          map[string]interface{}{},
			},
			ExpectError: false,
		},
		{
			Name: "top level - bool lists/sets",
			Input: &ExampleObj{
				BoolList: []bool{true, false, true},
				BoolSet:  []bool{false, true, false},
			},
			Expected: map[string]interface{}{
				"name":         "",
				"enabled":      false,
				"float":        float64(0),
				"list":         []interface{}{},
				"networks":     []string(nil),
				"networks_set": []string(nil),
				"int_list":     []int(nil),
				"int_set":      []int(nil),
				"float_list":   []float64(nil),
				"float_set":    []float64(nil),
				"bool_list":    []bool{true, false, true},
				"bool_set":     []bool{false, true, false},
				"number":       int64(0),
				"output":       "",
				"set":          []interface{}{},
				"map":          map[string]interface{}{},
			},
			ExpectError: false,
		},
		{
			Name: "top level - map",
			Input: &ExampleObj{
				Map: map[string]string{
					"bingo": "bango",
				},
			},
			Expected: map[string]interface{}{
				"name":         "",
				"enabled":      false,
				"float":        float64(0),
				"list":         []interface{}{},
				"networks":     []string(nil),
				"networks_set": []string(nil),
				"int_list":     []int(nil),
				"int_set":      []int(nil),
				"float_list":   []float64(nil),
				"float_set":    []float64(nil),
				"bool_list":    []bool(nil),
				"bool_set":     []bool(nil),
				"number":       int64(0),
				"output":       "",
				"set":          []interface{}{},
				"map": map[string]interface{}{
					"bingo": "bango",
				},
			},
			ExpectError: false,
		},
	}
	for _, v := range testCases {
		output, err := encodeHelper(v.Input)
		if err != nil {
			t.Fatalf("encoding error: %+v", err)
		}

		if !cmp.Equal(output, v.Expected) {
			t.Fatalf("Test Failed %q: output mismatch\n\n Expected: %+v\n\n Received: %+v\n\n", v.Name, v.Expected, output)
		}
	}
}

func encodeHelper(input interface{}) (map[string]interface{}, error) {
	objType := reflect.TypeOf(input).Elem()
	objVal := reflect.ValueOf(input).Elem()

	serialized, err := recurse(objType, objVal)
	if err != nil {
		return nil, err
	}
	return *serialized, nil
}
