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
		String        string            `hcl:"string"`
		Number        int               `hcl:"number"`
		Price         float64           `hcl:"price"`
		Enabled       bool              `hcl:"enabled"`
		ListOfFloats  []float64         `hcl:"list_of_floats"`
		ListOfNumbers []int             `hcl:"list_of_numbers"`
		ListOfStrings []string          `hcl:"list_of_strings"`
		MapOfBools    map[string]bool   `hcl:"map_of_bools"`   // TODO: fixme
		MapOfNumbers  map[string]int    `hcl:"map_of_numbers"` // TODO: fixme
		MapOfStrings  map[string]string `hcl:"map_of_strings"` // TODO: fixme
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
			MapOfBools: map[string]bool{
				"awesome_feature": true,
			},
			MapOfNumbers: map[string]int{
				"hello": 1,
				"there": 3,
			},
			MapOfStrings: map[string]string{
				"hello":   "there",
				"salut":   "tous les monde",
				"guten":   "tag",
				"morning": "alvaro",
			},
		},
		ExpectError: false,
	}.test(t)
}

func TestDecode_TopLevelFieldsOptional(t *testing.T) {
	type SimpleType struct {
		String        string            `hcl:"string"`
		Number        int               `hcl:"number"`
		Price         float64           `hcl:"price"`
		Enabled       bool              `hcl:"enabled"`
		ListOfFloats  []float64         `hcl:"list_of_floats"`
		ListOfNumbers []int             `hcl:"list_of_numbers"`
		ListOfStrings []string          `hcl:"list_of_strings"`
		MapOfBools    map[string]bool   `hcl:"map_of_bools"`
		MapOfNumbers  map[string]int    `hcl:"map_of_numbers"`
		MapOfStrings  map[string]string `hcl:"map_of_strings"`
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
			"map_of_bools":    map[string]interface{}{},
			"map_of_numbers":  map[string]interface{}{},
			"map_of_strings":  map[string]interface{}{},
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
