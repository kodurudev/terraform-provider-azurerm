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

func TestDecode_TopLevelFieldsRequired(t *testing.T) {
	type SimpleType struct {
		String        string    `hcl:"string"`
		Number        int       `hcl:"number"`
		Price         float64   `hcl:"price"`
		Enabled       bool      `hcl:"enabled"`
		ListOfFloats  []float64 `hcl:"list_of_floats"`
		ListOfNumbers []int     `hcl:"list_of_numbers"`
		ListOfStrings []string  `hcl:"list_of_strings"`
		//MapOfBools    map[string]bool   `hcl:"map_of_bools"`   // TODO: fixme
		//MapOfNumbers  map[string]int    `hcl:"map_of_numbers"` // TODO: fixme
		//MapOfStrings  map[string]string `hcl:"map_of_strings"` // TODO: fixme
	}
	decodeTestData{
		State: map[string]interface{}{
			"number":  int64(42),
			"price":   float64(129.99),
			"string":  "world",
			"enabled": true,
			"list_of_floats": []float64{
				1.0,
				2.0,
				3.0,
				1.234567890,
			},
			"list_of_numbers": []int{1, 2, 3},
			"list_of_strings": []string{
				"have",
				"you",
				"heard",
			},
			//"map_of_bools": map[string]interface{}{
			//	"awesome_feature": true,
			//},
			//"map_of_numbers": map[string]interface{}{
			//	"hello": 1,
			//	"there": 3,
			//},
			//"map_of_strings": map[string]interface{}{
			//	"hello":   "there",
			//	"salut":   "tous les monde",
			//	"guten":   "tag",
			//	"morning": "alvaro",
			//},
		},
		Input: &SimpleType{},
		Expected: &SimpleType{
			String:  "world",
			Number:  42,
			Price:   129.99,
			Enabled: true,
			ListOfFloats: []float64{
				1.0,
				2.0,
				3.0,
				1.234567890},
			ListOfNumbers: []int{1, 2, 3},
			ListOfStrings: []string{
				"have",
				"you",
				"heard",
			},
			//MapOfBools: map[string]bool{
			//	"awesome_feature": true,
			//},
			//MapOfNumbers: map[string]int{
			//	"hello": 1,
			//	"there": 3,
			//},
			//MapOfStrings: map[string]string{
			//	"hello":   "there",
			//	"salut":   "tous les monde",
			//	"guten":   "tag",
			//	"morning": "alvaro",
			//},
		},
		ExpectError: false,
	}.test(t)
}

func TestDecode_TopLevelFieldsOptional(t *testing.T) {
	type SimpleType struct {
		String        string    `hcl:"string"`
		Number        int       `hcl:"number"`
		Price         float64   `hcl:"price"`
		Enabled       bool      `hcl:"enabled"`
		ListOfFloats  []float64 `hcl:"list_of_floats"`
		ListOfNumbers []int     `hcl:"list_of_numbers"`
		ListOfStrings []string  `hcl:"list_of_strings"`
		//MapOfBools    map[string]bool   `hcl:"map_of_bools"` // TODO: fix me
		//MapOfNumbers  map[string]int    `hcl:"map_of_numbers"` // TODO: fix me
		//MapOfStrings  map[string]string `hcl:"map_of_strings"` // TODO: fix me
	}
	decodeTestData{
		State: map[string]interface{}{
			"number":          int64(0),
			"price":           float64(0),
			"string":          "",
			"enabled":         false,
			"list_of_floats":  []float64{},
			"list_of_numbers": []int{},
			"list_of_strings": []string{},
			//"map_of_bools":    map[string]interface{}{},
			//"map_of_numbers":  map[string]interface{}{},
			//"map_of_strings":  map[string]interface{}{},
		},
		Input:       &SimpleType{},
		Expected:    &SimpleType{},
		ExpectError: false,
	}.test(t)
}

func TestDecode_TopLevelFieldsComputed(t *testing.T) {
	type SimpleType struct {
		ComputedString        string   `hcl:"computed_string" computed:"true"`
		ComputedNumber        int      `hcl:"computed_number" computed:"true"`
		ComputedBool          bool     `hcl:"computed_bool" computed:"true"`
		ComputedListOfNumbers []int    `hcl:"computed_list_of_numbers" computed:"true"`
		ComputedListOfStrings []string `hcl:"computed_list_of_strings" computed:"true"`
		// TODO: computed maps
	}
	decodeTestData{
		State: map[string]interface{}{
			"computed_string":          "je suis computed",
			"computed_number":          int64(732),
			"computed_bool":            true,
			"computed_list_of_numbers": []int{1, 2, 3},
			"computed_list_of_strings": []string{
				"have",
				"you",
				"heard",
			},
		},
		Input: &SimpleType{},
		Expected: &SimpleType{
			ComputedString:        "je suis computed",
			ComputedNumber:        732,
			ComputedBool:          true,
			ComputedListOfNumbers: []int{1, 2, 3},
			ComputedListOfStrings: []string{
				"have",
				"you",
				"heard",
			},
		},
		ExpectError: false,
	}.test(t)
}

func TestResourceDecode_NestedOneLevelDeepEmpty(t *testing.T) {
	type Inner struct {
		Value string `hcl:"value"`
	}
	type Type struct {
		NestedObject []Inner `hcl:"inner"`
	}
	decodeTestData{
		State: map[string]interface{}{
			"inner": []interface{}{},
		},
		Input: &Type{},
		Expected: &Type{
			NestedObject: []Inner{}, // TODO: has to be `[]Inner(nil)` right now
		},
	}.test(t)
}

func TestResourceDecode_NestedOneLevelDeepSingle(t *testing.T) {
	type Inner struct {
		String        string    `hcl:"string"`
		Number        int       `hcl:"number"`
		Price         float64   `hcl:"price"`
		Enabled       bool      `hcl:"enabled"`
		ListOfFloats  []float64 `hcl:"list_of_floats"`
		ListOfNumbers []int     `hcl:"list_of_numbers"`
		ListOfStrings []string  `hcl:"list_of_strings"`
		//MapOfBools    map[string]bool   `hcl:"map_of_bools"` // TODO: fixme
		//MapOfNumbers  map[string]int    `hcl:"map_of_numbers"` // TODO: fixme
		//MapOfStrings  map[string]string `hcl:"map_of_strings"` // TODO: fixme
	}
	type Type struct {
		NestedObject []Inner `hcl:"inner"`
	}
	decodeTestData{
		State: map[string]interface{}{
			"inner": []interface{}{
				map[string]interface{}{
					"number":  int64(42),
					"price":   float64(129.99),
					"string":  "world",
					"enabled": true,
					"list_of_floats": []float64{
						1.0,
						2.0,
						3.0,
						1.234567890,
					},
					"list_of_numbers": []int{1, 2, 3},
					"list_of_strings": []string{
						"have",
						"you",
						"heard",
					},
					//"map_of_bools": map[string]interface{}{
					//	"awesome_feature": true,
					//},
					//"map_of_numbers": map[string]interface{}{
					//	"hello": 1,
					//	"there": 3,
					//},
					//"map_of_strings": map[string]interface{}{
					//	"hello":   "there",
					//	"salut":   "tous les monde",
					//	"guten":   "tag",
					//	"morning": "alvaro",
					//},
				},
			},
		},
		Input: &Type{},
		Expected: &Type{
			NestedObject: []Inner{
				{
					String:  "world",
					Number:  42,
					Price:   129.99,
					Enabled: true,
					ListOfFloats: []float64{
						1.0,
						2.0,
						3.0,
						1.234567890},
					ListOfNumbers: []int{1, 2, 3},
					ListOfStrings: []string{
						"have",
						"you",
						"heard",
					},
					//MapOfBools: map[string]bool{
					//	"awesome_feature": true,
					//},
					//MapOfNumbers: map[string]int{
					//	"hello": 1,
					//	"there": 3,
					//},
					//MapOfStrings: map[string]string{
					//	"hello":   "there",
					//	"salut":   "tous les monde",
					//	"guten":   "tag",
					//	"morning": "alvaro",
					//},
				},
			},
		},
	}.test(t)
}

func TestResourceDecode_NestedOneLevelDeepSingleOmittedValues(t *testing.T) {
	type Inner struct {
		String        string    `hcl:"string"`
		Number        int       `hcl:"number"`
		Price         float64   `hcl:"price"`
		Enabled       bool      `hcl:"enabled"`
		ListOfFloats  []float64 `hcl:"list_of_floats"`
		ListOfNumbers []int     `hcl:"list_of_numbers"`
		ListOfStrings []string  `hcl:"list_of_strings"`
		//MapOfBools    map[string]bool   `hcl:"map_of_bools"`
		//MapOfNumbers  map[string]int    `hcl:"map_of_numbers"`
		//MapOfStrings  map[string]string `hcl:"map_of_strings"`
	}
	type Type struct {
		NestedObject []Inner `hcl:"inner"`
	}
	decodeTestData{
		State: map[string]interface{}{
			"inner": []interface{}{
				map[string]interface{}{
					"number":          int64(0),
					"price":           float64(0),
					"string":          "",
					"enabled":         false,
					"list_of_floats":  []float64{},
					"list_of_numbers": []int{},
					"list_of_strings": []string{},
					//"map_of_bools":    map[string]interface{}{},
					//"map_of_numbers":  map[string]interface{}{},
					//"map_of_strings":  map[string]interface{}{},
				},
			},
		},
		Input: &Type{},
		Expected: &Type{
			NestedObject: []Inner{
				{},
			},
		},
	}.test(t)
}

func TestResourceDecode_NestedOneLevelDeepSingleMultiple(t *testing.T) {
	type Inner struct {
		Value string `hcl:"value"`
	}
	type Type struct {
		NestedObject []Inner `hcl:"inner"`
	}
	decodeTestData{
		State: map[string]interface{}{
			"inner": []interface{}{
				map[string]interface{}{
					"value": "first",
				},
				map[string]interface{}{
					"value": "second",
				},
				map[string]interface{}{
					"value": "third",
				},
			},
		},
		Input: &Type{},
		Expected: &Type{
			NestedObject: []Inner{
				{
					Value: "first",
				},
				{
					Value: "second",
				},
				{
					Value: "third",
				},
			},
		},
	}.test(t)
}

func TestResourceDecode_NestedThreeLevelsDeepEmpty(t *testing.T) {
	type ThirdInner struct {
		Value string `hcl:"value"`
	}
	type SecondInner struct {
		Third []ThirdInner `hcl:"third"`
	}
	type FirstInner struct {
		Second []SecondInner `hcl:"second"`
	}
	type Type struct {
		First []FirstInner `hcl:"first"`
	}

	t.Log("Top Level Empty")
	decodeTestData{
		State: map[string]interface{}{
			"first": []interface{}{},
		},
		Input: &Type{},
		Expected: &Type{
			First: []FirstInner{},
		},
	}.test(t)

	t.Log("Second Level Empty")
	decodeTestData{
		State: map[string]interface{}{
			"first": []interface{}{
				map[string]interface{}{
					"second": []interface{}{},
				},
			},
		},
		Input: &Type{},
		Expected: &Type{
			First: []FirstInner{
				{
					Second: []SecondInner{},
				},
			},
		},
	}.test(t)

	t.Log("Third Level Empty")
	decodeTestData{
		State: map[string]interface{}{
			"first": []interface{}{
				map[string]interface{}{
					"second": []interface{}{
						map[string]interface{}{
							"third": []interface{}{},
						},
					},
				},
			},
		},
		Input: &Type{},
		Expected: Type{
			First: []FirstInner{
				{
					Second: []SecondInner{
						{
							Third: []ThirdInner{},
						},
					},
				},
			},
		},
	}.test(t)
}

func TestResourceDecode_NestedThreeLevelsDeepSingleItem(t *testing.T) {
	type ThirdInner struct {
		Value string `hcl:"value"`
	}
	type SecondInner struct {
		Third []ThirdInner `hcl:"third"`
	}
	type FirstInner struct {
		Second []SecondInner `hcl:"second"`
	}
	type Type struct {
		First []FirstInner `hcl:"first"`
	}

	decodeTestData{
		State: map[string]interface{}{
			"first": []interface{}{
				map[string]interface{}{
					"second": []interface{}{
						map[string]interface{}{
							"third": []interface{}{
								map[string]interface{}{
									"value": "salut",
								},
							},
						},
					},
				},
			},
		},
		Input: &Type{},
		Expected: &Type{
			First: []FirstInner{
				{
					Second: []SecondInner{
						{
							Third: []ThirdInner{
								{
									Value: "salut",
								},
							},
						},
					},
				},
			},
		},
	}.test(t)
}

func TestResourceDecode_NestedThreeLevelsDeepMultipleItems(t *testing.T) {
	type ThirdInner struct {
		Value string `hcl:"value"`
	}
	type SecondInner struct {
		Value string       `hcl:"value"`
		Third []ThirdInner `hcl:"third"`
	}
	type FirstInner struct {
		Value  string        `hcl:"value"`
		Second []SecondInner `hcl:"second"`
	}
	type Type struct {
		First []FirstInner `hcl:"first"`
	}

	decodeTestData{
		State: map[string]interface{}{
			"first": []interface{}{
				map[string]interface{}{
					"value": "first - 1",
					"second": []interface{}{
						map[string]interface{}{
							"value": "second - 1",
							"third": []interface{}{
								map[string]interface{}{
									"value": "third - 1",
								},
								map[string]interface{}{
									"value": "third - 2",
								},
								map[string]interface{}{
									"value": "third - 3",
								},
							},
						},
						map[string]interface{}{
							"value": "second - 2",
							"third": []interface{}{
								map[string]interface{}{
									"value": "third - 4",
								},
								map[string]interface{}{
									"value": "third - 5",
								},
								map[string]interface{}{
									"value": "third - 6",
								},
							},
						},
					},
				},
				map[string]interface{}{
					"value": "first - 2",
					"second": []interface{}{
						map[string]interface{}{
							"value": "second - 3",
							"third": []interface{}{
								map[string]interface{}{
									"value": "third - 7",
								},
								map[string]interface{}{
									"value": "third - 8",
								},
							},
						},
						map[string]interface{}{
							"value": "second - 4",
							"third": []interface{}{
								map[string]interface{}{
									"value": "third - 9",
								},
							},
						},
					},
				},
			},
		},
		Input: &Type{},
		Expected: &Type{
			First: []FirstInner{
				{
					Value: "first - 1",
					Second: []SecondInner{
						{
							Value: "second - 1",
							Third: []ThirdInner{
								{
									Value: "third - 1",
								},
								{
									Value: "third - 2",
								},
								{
									Value: "third - 3",
								},
							},
						},
						{
							Value: "second - 2",
							Third: []ThirdInner{
								{
									Value: "third - 4",
								},
								{
									Value: "third - 5",
								},
								{
									Value: "third - 6",
								},
							},
						},
					},
				},
				{
					Value: "first - 2",
					Second: []SecondInner{
						{
							Value: "second - 3",
							Third: []ThirdInner{
								{
									Value: "third - 7",
								},
								{
									Value: "third - 8",
								},
							},
						},
						{
							Value: "second - 4",
							Third: []ThirdInner{
								{
									Value: "third - 9",
								},
							},
						},
					},
				},
			},
		},
	}.test(t)
}

func TestDecode(t *testing.T) {
	decodeTestData{
		State: map[string]interface{}{
			"map": map[string]interface{}{
				"bingo": "bango",
			},
		},
		Input: &ExampleObj{},
		Expected: &ExampleObj{
			Map: map[string]string{
				"bingo": "bango",
			},
		},
		ExpectError: false,
	}.test(t)
}

func (testData decodeTestData) test(t *testing.T) {
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
		t.Fatalf("Expected: %+v\n\n Received %+v\n\n", testData.Expected, testData.Input)
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
