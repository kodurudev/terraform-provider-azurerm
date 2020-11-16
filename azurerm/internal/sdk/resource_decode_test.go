package sdk

import (
	"reflect"
	"testing"
)

type decodeTestData struct {
	State       map[string]interface{}
	Input       interface{}
	Expected    interface{}
	ExpectError bool
}

func TestDecode(t *testing.T) {
	testCases := []struct {
		Name        string
		State       map[string]interface{}
		Expected    *ExampleObj
		ExpectError bool
	}{ /*
			{
				Name: "top level - name",
				State: map[string]interface{}{
					"name": "bingo bango",
				},
				Expected: &ExampleObj{
					Name: "bingo bango",
				},
				ExpectError: false,
			},
			{
				Name: "top level - everything",
				State: map[string]interface{}{
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
				State: map[string]interface{}{
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
				State: map[string]interface{}{
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
				State: map[string]interface{}{
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
				State: map[string]interface{}{
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
				State: map[string]interface{}{
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
				State: map[string]interface{}{
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
				State: map[string]interface{}{
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
			State: map[string]interface{}{
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
		test := decodeTestData{
			State:       v.State,
			Input:       &ExampleObj{},
			Expected:    v.Expected,
			ExpectError: v.ExpectError,
		}
		test.test(t)
	}
}

func (testData *decodeTestData) test(t *testing.T) {
	debugLogger := ConsoleLogger{}
	state := testData.stateWrapper()
	if err := decodeReflectedType(testData.Input, state, debugLogger); err != nil {
		if testData.ExpectError {
			// we're good
			return
		}

		t.Fatalf("unexpected error: %+v", err)
	}
	if testData.ExpectError {
		t.Fatalf("expected an error but didn't get one!")
	}

	if !reflect.DeepEqual(testData.Input, testData.Expected) {
		t.Fatalf("Expected: %+v\n\n Received %+v\n\n", testData.Input, testData.Expected)
	}
}

func (testData decodeTestData) stateWrapper() testDataGetter {
	return testDataGetter{
		values: testData.State,
	}
}

type testDataGetter struct {
	values map[string]interface{}
}

func (td testDataGetter) Get(key string) interface{} {
	return td.values[key]
}
func (td testDataGetter) GetOk(key string) (interface{}, bool) {
	val, ok := td.values[key]
	return val, ok
}
