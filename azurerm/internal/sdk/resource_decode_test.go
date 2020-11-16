package sdk

import (
	"reflect"
	"testing"
)

func TestDecode(t *testing.T) {
	testCases := []struct {
		Name        string
		Input       map[string]interface{}
		Expected    *ExampleObj
		ExpectError bool
	}{ /*
			{
				Name: "top level - name",
				Input: map[string]interface{}{
					"name": "bingo bango",
				},
				Expected: &ExampleObj{
					Name: "bingo bango",
				},
				ExpectError: false,
			},
			{
				Name: "top level - everything",
				Input: map[string]interface{}{
					"name": "bingo bango",
					"float": 123.4,
					"number": 123,
					"enabled": false,
				},
				Expected: &ExampleObj{
					Name: "bingo bango",
					Float: 123.4,
					Number: 123,
					Enabled: false,
				},
				ExpectError: false,
			},
			{
				Name: "top level - list",
				Input: map[string]interface{}{
					"name": "bingo bango",
					"float": 123.4,
					"number": 123,
					"enabled": false,
					"list": []interface{}{
						map[string]interface{}{
							"name": "first",
						},
					},
				},
				Expected: &ExampleObj{
					Name: "bingo bango",
					Float: 123.4,
					Number: 123,
					Enabled: false,
					List: []NetworkList{{
						Name: "first",
					}},
				},
				ExpectError: false,
			},
			{
				Name: "top level - list in lists",
				Input: map[string]interface{}{
					"name": "bingo bango",
					"float": 123.4,
					"number": 123,
					"enabled": false,
					"list": []interface{}{
						map[string]interface{}{
							"name": "first",
							"inner": []interface{}{
								map[string]interface{}{
									"name": "get-a-mac",
								},
							},
						},
					},
				},
				Expected: &ExampleObj{
					Name: "bingo bango",
					Float: 123.4,
					Number: 123,
					Enabled: false,
					List: []NetworkList{{
						Name: "first",
						Inner: []NetworkInner{{
							Name: "get-a-mac",
						}},
					}},
				},
				ExpectError: false,
			},
			{
				Name: "top level - everything",
				Input: map[string]interface{}{
					"name": "bingo bango",
					"float": 123.4,
					"number": 123,
					"enabled": false,
					"networks": []interface{}{"network1", "network2", "network3"},
					"networks_set": []interface{}{"networkset1", "networkset2", "networkset3"},
					"list": []interface{}{
						map[string]interface{}{
							"name": "first",
							"inner": []interface{}{
								map[string]interface{}{
									"name": "get-a-mac",
								},
							},
						},
					},
					"set": schema.NewSet(FakeHashSchema(),
						[]interface{}{
							map[string]interface{}{
								"name": "setname",
							},
						}),
				},
				Expected: &ExampleObj{
					Name: "bingo bango",
					Float: 123.4,
					Number: 123,
					Enabled: false,
					Networks: []string{"network1", "network2", "network3"},
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
				ExpectError: false,
			},
			{
				Name: "nests",
				Input: map[string]interface{}{
					"name": "bingo bango",
					"float": 123.4,
					"number": 123,
					"enabled": false,
					"networks": []interface{}{"network1", "network2", "network3"},
					"networks_set": []interface{}{"networkset1", "networkset2", "networkset3"},
					"list": []interface{}{
						map[string]interface{}{
							"name": "first",
							"inner": []interface{}{
								map[string]interface{}{
									"name": "get-a-mac",
									"inner": []interface{}{
										map[string]interface{}{
											"name": "innerinner",
											"should_be_fine": true,
										},
									},
									"set": schema.NewSet(FakeHashSchema(),
										[]interface{}{
											map[string]interface{}{
												"name": "nestedsetname",
											},
										}),
								},
							},
						},
						map[string]interface{}{
							"name": "second",
							"inner": []interface{}{
								map[string]interface{}{
									"name": "get-a-mac2",
									"inner": []interface{}{
										map[string]interface{}{
											"name": "innerinner2",
											"should_be_fine": true,
										},
									},
									"set": schema.NewSet(FakeHashSchema(),
										[]interface{}{
											map[string]interface{}{
												"name": "nestedsetname2",
											},
										}),
								},
							},
						},
					},
					"set": schema.NewSet(FakeHashSchema(),
						[]interface{}{
							map[string]interface{}{
								"name": "setname",
							},
						}),
				},
				Expected: &ExampleObj{
					Name: "bingo bango",
					Float: 123.4,
					Number: 123,
					Enabled: false,
					Networks: []string{"network1", "network2", "network3"},
					NetworksSet: []string{"networkset1", "networkset2", "networkset3"},
					List: []NetworkList{
						{
							Name: "first",
							Inner: []NetworkInner{{
								Name: "get-a-mac",
								Inner: []InnerInner{{
									Name: "innerinner",
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
									Name: "innerinner2",
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
				ExpectError: false,
			},
			{
				Name: "top level - int lists/sets",
				Input: map[string]interface{}{
					"int_list": []interface{}{1,2,3},
					"int_set": []interface{}{3,4,5},
				},
				Expected: &ExampleObj{
					IntList: []int{1,2,3},
					IntSet: []int{3,4,5},
				},
				ExpectError: false,
			},
			{
				Name: "top level - float lists/sets",
				Input: map[string]interface{}{
					"float_list": []interface{}{1.1,2.2,3.3},
					"float_set": []interface{}{3.3,4.4,5.5},
				},
				Expected: &ExampleObj{
					FloatList: []float64{1.1,2.2,3.3},
					FloatSet: []float64{3.3,4.4,5.5},
				},
				ExpectError: false,
			},
			{
				Name: "top level - bool lists/sets",
				Input: map[string]interface{}{
					"bool_list": []interface{}{true,false,true},
					"bool_set": []interface{}{false,true,false},
				},
				Expected: &ExampleObj{
					BoolList: []bool{true,false,true},
					BoolSet: []bool{false,true,false},
				},
				ExpectError: false,
			},*/
		{
			Name: "top level - map",
			Input: map[string]interface{}{
				"map": map[string]interface{}{
					"bingo": "bango",
				},
			},
			Expected: &ExampleObj{
				Map: map[string]string{
					"bingo": "bango",
				},
			},
			ExpectError: false,
		},
	}

	for _, v := range testCases {
		obj := &ExampleObj{}
		if err := decodeHelper(obj, v.Input); err != nil && !v.ExpectError {
			t.Fatalf("error decoding into obj: %+v", err)
		}

		if !reflect.DeepEqual(obj, v.Expected) {
			t.Fatalf("ExampleObj mismatch\n\n Expected: %+v\n\n Received %+v\n\n", v.Expected, obj)
		}
	}
}

func decodeHelper(input interface{}, config map[string]interface{}) error {
	objType := reflect.TypeOf(input).Elem()
	for i := 0; i < objType.NumField(); i++ {
		field := objType.Field(i)

		if val, exists := field.Tag.Lookup("computed"); exists {
			if val == "true" {
				continue
			}
		}

		if val, exists := field.Tag.Lookup("hcl"); exists {
			hclValue := config[val]

			//TODO Actually check error
			if err := setValue(input, hclValue, i); err != nil {
				return err
			}
		}
	}
	return nil
}
