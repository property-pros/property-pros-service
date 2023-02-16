package interop

import (
	"testing"
)

func TestMapUnaryTypes(t *testing.T) {
	type TestStruct struct {
		FirstName string
		LastName  string
		Age       int
	}

	type TestStruct2 struct {
		FirstName string
		LastName  string
		Age       int
	}

	mapper := NewMapper()

	testStruct := TestStruct{
		FirstName: "John",
		LastName:  "Doe",
		Age:       42,
	}

	var testStruct2 TestStruct2

	err := mapper.Map(testStruct, &testStruct2)

	if err != nil {
		t.Error(err)
	}

	if testStruct.FirstName != testStruct2.FirstName {
		t.Error("FirstName does not match")
	}

	if testStruct.LastName != testStruct2.LastName {
		t.Error("LastName does not match")
	}

	if testStruct.Age != testStruct2.Age {
		t.Error("Age does not match")
	}
}

func TestMapCustomMappingUnaryTypes(t *testing.T) {
	type TestStruct struct {
		FirstName string
		LastName  string
		Age       int
	}

	type TestStruct2 struct {
		FirstName string
		LastName  string
		Age       int
	}

	mapper := NewMapper()

	mapper.AddCustomMapping("FirstName", "LastName", func(from interface{}, to interface{}) interface{} {
		t.Logf("from: %+#v \n\n to: %+#v \n\n", from, to)
		fromVal := from.(*TestStruct)

		return fromVal.FirstName
	})

	testStruct := &TestStruct{
		FirstName: "John",
		LastName:  "Doe",
		Age:       42,
	}

	testStruct2 := &TestStruct2{}

	err := mapper.Map(testStruct, testStruct2)

	t.Logf("testStruct: %+#v; testStruct2: %+#v \n\n", testStruct, testStruct2)
	if err != nil {
		t.Error(err)
	}

	if testStruct.FirstName != testStruct2.FirstName {
		t.Error("FirstName does not match")
	}

	if testStruct.FirstName != testStruct2.LastName {
		t.Error("LastName does not match")
	}

	if testStruct.Age != testStruct2.Age {
		t.Error("Age does not match")
	}
}

func TestMapSlice(t *testing.T) {
	type TestStruct struct {
		FirstName string
		LastName  string
		Age       int
	}

	type TestStruct2 struct {
		FirstName string
		LastName  string
		Age       int
	}

	mapper := NewMapper()

	testStruct := []TestStruct{
		{
			FirstName: "John",
			LastName:  "Doe",
			Age:       42,
		},
		{
			FirstName: "Jane",
			LastName:  "Doe",
			Age:       42,
		},
	}

	var testStruct2 []TestStruct2

	err := mapper.Map(testStruct, &testStruct2)

	if err != nil {
		t.Error(err)
	}

	if len(testStruct) != len(testStruct2) {
		t.Error("Length of slices does not match")
	}

	for i := 0; i < len(testStruct); i++ {
		if testStruct[i].FirstName != testStruct2[i].FirstName {
			t.Error("FirstName does not match")
		}

		if testStruct[i].LastName != testStruct2[i].LastName {
			t.Error("LastName does not match")
		}

		if testStruct[i].Age != testStruct2[i].Age {
			t.Error("Age does not match")
		}
	}
}

func TestMapArray(t *testing.T) {
	type TestStruct struct {
		FirstName string
		LastName  string
		Age       int
	}

	type TestStruct2 struct {
		FirstName string
		LastName  string
		Age       int
	}

	mapper := NewMapper()

	testStruct := [2]TestStruct{
		{
			FirstName: "John",
			LastName:  "Doe",
			Age:       42,
		},
		{
			FirstName: "Jane",
			LastName:  "Doe",
			Age:       42,
		},
	}

	var testStruct2 [2]TestStruct2

	err := mapper.Map(testStruct, &testStruct2)

	if err != nil {
		t.Error(err)
	}

	for i := 0; i < len(testStruct); i++ {
		if testStruct[i].FirstName != testStruct2[i].FirstName {
			t.Error("FirstName does not match")
		}

		if testStruct[i].LastName != testStruct2[i].LastName {
			t.Error("LastName does not match")
		}

		if testStruct[i].Age != testStruct2[i].Age {
			t.Error("Age does not match")
		}
	}
}

func TestMap2StructsWithStructFields(t *testing.T) {
	type TestStruct struct {
		FirstName string
		LastName  string
		Age       int
	}

	type TestStruct2 struct {
		FirstName string
		LastName  string
		Age       int
	}

	type TestStruct3 struct {
		FirstName string
		LastName  string
		Age       int
		TestStruct
	}

	type TestStruct4 struct {
		FirstName string
		LastName  string
		Age       int
		TestStruct2
	}

	mapper := NewMapper()

	testStruct := TestStruct3{
		FirstName: "John",
		LastName:  "Doe",
		Age:       42,
		TestStruct: TestStruct{
			FirstName: "Jane",
			LastName:  "Doe",
			Age:       42,
		},
	}

	var testStruct2 TestStruct4

	err := mapper.Map(testStruct, &testStruct2)

	if err != nil {
		t.Error(err)
	}

	if testStruct.FirstName != testStruct2.FirstName {
		t.Error("FirstName does not match")
	}

	if testStruct.LastName != testStruct2.LastName {
		t.Error("LastName does not match")
	}

	if testStruct.Age != testStruct2.Age {
		t.Error("Age does not match")
	}

	if testStruct.TestStruct.FirstName != testStruct2.TestStruct2.FirstName {
		t.Error("FirstName does not match")
	}

	if testStruct.TestStruct.LastName != testStruct2.TestStruct2.LastName {
		t.Error("LastName does not match")
	}

	if testStruct.TestStruct.Age != testStruct2.TestStruct2.Age {
		t.Error("Age does not match")
	}
}
