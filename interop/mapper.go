package interop

import (
	"fmt"
	"reflect"
)

// Mapper is a struct that will map any struct to any other struct with matching field names;

type Mapper struct {
	customMappings map[string]map[string]func(interface{}, interface{}) interface{}
}

// Map will map any struct to any other struct with matching field names;

func (m *Mapper) Map(from interface{}, to interface{}) error {
	//get value of from
	fromValue := m.extractReflectedValue(from)
	//get type of from
	fromType := fromValue.Type()
	//get value of to
	toValue := m.extractReflectedValue(to)
	//get type of to
	toType := toValue.Type()
	//get kind of from
	// fromKind := fromType.Kind()
	//get kind of to
	// toKind := toType.Kind()
	//check if from is a struct

	// handledSuccessfully  := tryHandleIterable(fromKind, toKind, fromValue, toValue, m);

	//get number of fields in from
	fromNumField := fromType.NumField()
	//get number of fields in to
	toNumField := toType.NumField()
	//loop through fields in from
	for i := 0; i < fromNumField; i++ {
		//get field from from
		fromField := fromType.Field(i)
		//get field value from from
		fromFieldValue := fromValue.Field(i)
		//loop through fields in to
		for j := 0; j < toNumField; j++ {
			//get field from to
			toField := toType.Field(j)
			//get field value from to
			toFieldValue := toValue.Field(j)
			//check if field names match
			if m.fieldsMatch(fromField, toField) {
				m.mapFields(fromField, fromFieldValue, toField, toFieldValue, from, to)
			}

		}
	}

	return nil
}

func (m *Mapper) extractReflectedValue(value interface{}) reflect.Value {
	//get value of from
	valueOf := reflect.ValueOf(value)
	//check if value is a pointer
	if valueOf.Kind() == reflect.Ptr {
		//get value of pointer
		valueOf = valueOf.Elem()
	}
	return valueOf
}

// fix mapFields

func (m *Mapper) mapFields(fromField reflect.StructField, fromFieldValue reflect.Value, toField reflect.StructField, toFieldValue reflect.Value, from interface{}, to interface{}) {
	//check if fields are mappable
	if m.fieldsAreMappable(fromField.Type, toField.Type) {
		m.Map(fromFieldValue.Interface(), toFieldValue.Interface())
	} else {

		if m.typesAreConvertible(fromField, toField) {
			// convert value of from field to value of to field
			fromFieldValue = fromFieldValue.Convert(toField.Type)
		}

		if m.matchesCustomMapping(fromField.Name, toField.Name) {

			customMappedValue := m.customMappings[fromField.Name][toField.Name](from, to)

			reflectedMappedValue := reflect.ValueOf(customMappedValue)

			//set value of to field to value of from field
			toFieldValue.Set(reflectedMappedValue)
		} else {

			//skip if toField has a custom mapping
			if m.hasCustomMapping(fromField.Name) {
				return
			}

			fmt.Printf("fromFieldValue: %v \n\n", fromFieldValue.Interface())
			//set value of to field to value of from field
			toFieldValue.Set(fromFieldValue)
		}
	}
}

func (m *Mapper) hasCustomMapping(fieldName string) bool {
	fmt.Printf("customMappings: %v; field name: %v \n\n", m.customMappings, fieldName)
	if m.customMappings != nil {
		// if m.customMappings[fieldName] != nil {
		// 	return true
		// } 
		// else {
			for _, customMapping := range m.customMappings {
				if customMapping[fieldName] != nil {
					return true
				}
			}
		// }
	}

	return false
}

// fieldsMatch will check if the field names match

func (m *Mapper) fieldsMatch(fromField reflect.StructField, toField reflect.StructField) bool {

	return m.fieldTypesMatch(fromField, toField) &&
		m.fieldNamesMatch(fromField, toField)
}

// fieldNamesMatch will check if the field names match

func (m *Mapper) fieldNamesMatch(fromField reflect.StructField, toField reflect.StructField) bool {
	return fromField.Name == toField.Name || m.matchesCustomMapping(fromField.Name, toField.Name)
}

// fieldTypesMatch will check if the field types match

func (m *Mapper) fieldTypesMatch(fromField reflect.StructField, toField reflect.StructField) bool {
	return m.typesAreEqual(fromField, toField) || m.typesAreConvertible(fromField, toField)
}

// typesAreEqual will check if the types are equal

func (m *Mapper) typesAreEqual(fromField reflect.StructField, toField reflect.StructField) bool {
	return fromField.Type == toField.Type
}

// typesAreConvertible will check if the types are convertible

func (m *Mapper) typesAreConvertible(fromField reflect.StructField, toField reflect.StructField) bool {
	return fromField.Type.ConvertibleTo(toField.Type)
}

// matchesCustomMapping will check if a custom mapping exists for the given field names

func (m *Mapper) matchesCustomMapping(fromField string, toField string) bool {
	if m.customMappings != nil {
		if m.customMappings[fromField] != nil {
			if m.customMappings[fromField][toField] != nil {
				return true
			}
		}
	}

	return false
}

func (m *Mapper) fieldsAreMappable(type1, type2 reflect.Type) bool {
	return type1.Kind() == reflect.Struct && type2.Kind() == reflect.Struct ||
		type1.Kind() == reflect.Slice && type2.Kind() == reflect.Slice ||
		type1.Kind() == reflect.Ptr && type2.Kind() == reflect.Ptr ||
		type1.Kind() == reflect.Map && type2.Kind() == reflect.Map ||
		type1.Kind() == reflect.Array && type2.Kind() == reflect.Array
}

// AddCustomMapping will add a custom mapping to the mapper

func (m *Mapper) AddCustomMapping(fromField string, toField string, mapping func(interface{}, interface{}) interface{}) {

	if m.customMappings[fromField] == nil {
		m.customMappings[fromField] = make(map[string]func(interface{}, interface{}) interface{})
	}

	m.customMappings[fromField][toField] = mapping
}

// NewMapper will create a new mapper

func NewMapper() *Mapper {
	return &Mapper{
		customMappings: make(map[string]map[string]func(interface{}, interface{}) interface{}),
	}
}
