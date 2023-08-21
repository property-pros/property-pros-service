package interop

import (
	"fmt"
	"reflect"
)

type Mapper struct {
	customMappings map[string]map[string]interface{}
}

func (m *Mapper) Map(from interface{}, to interface{}) error {
	fromValue := m.extractReflectedValue(from)
	toValue := m.extractReflectedValue(to)

	if fromValue.Kind() == reflect.Ptr {
		fromValue = fromValue.Elem()
	}

	if toValue.Kind() == reflect.Ptr {
		toValue = toValue.Elem()
	}

	fromType := fromValue.Type()
	toType := toValue.Type()

	if fromType.Kind() == reflect.Struct && toType.Kind() == reflect.Struct {
		return m.mapStruct(fromValue, toValue, fromType, toType)
	}

	if fromType.Kind() == reflect.Slice && toType.Kind() == reflect.Slice {
		return m.mapSlice(fromValue, toValue, fromType.Elem(), toType.Elem())
	}

	if fromType.Kind() == reflect.Array && toType.Kind() == reflect.Array {
		return m.mapArray(fromValue, toValue, fromType.Elem(), toType.Elem())
	}

	return fmt.Errorf("unsupported type %s", fromValue.Type().Kind())
}
func (m *Mapper) mapStruct(fromValue reflect.Value, toValue reflect.Value, fromType reflect.Type, toType reflect.Type) error {
	for i := 0; i < fromType.NumField(); i++ {
		fromField := fromType.Field(i)
		fromFieldValue := fromValue.Field(i)

		if fromField.Type.Kind() == reflect.Struct {
			err := m.Map(fromFieldValue.Interface(), toValue.Field(i).Addr().Interface())
			if err != nil {
				return err
			}
			continue
		}

		toField, found := toValue.Type().FieldByName(fromField.Name)
		if found {
			toFieldValue := toValue.FieldByName(toField.Name)

			fmt.Printf("fromField: %s, toField: %s\n", fromField.Name, toField.Name)
			fmt.Printf("custom from mapping: %v\n", m.customMappings[fromField.Name])
			fmt.Printf("custom from/to mapping: %v\n", m.customMappings[fromField.Name][toField.Name])
			
			if customMapping, ok := m.customMappings[fromField.Name][toField.Name]; ok {
				switch customFunc := customMapping.(type) {
				case func(from interface{}, to interface{}) interface{}:
					mappedValue := customFunc(fromFieldValue.Interface(), toFieldValue.Interface())
					if toFieldValue.CanSet() {
						toFieldValue.Set(reflect.ValueOf(mappedValue))
					}
				case string:
					mappedField := fromFieldValue.FieldByName(customFunc)
					if mappedField.IsValid() && toFieldValue.CanSet() {
						toFieldValue.Set(mappedField)
					}
				default:
					return fmt.Errorf("unsupported custom mapping type: %T", customMapping)
				}
			} else if fromFieldValue.Type().AssignableTo(toFieldValue.Type()) && toFieldValue.CanSet() {
				toFieldValue.Set(fromFieldValue)
			}
		}
	}
	return nil
}

func (m *Mapper) mapWithCustomMapping(customMappingFuncs map[string]interface{}, fromField, toField reflect.StructField, fromFieldValue, toFieldValue reflect.Value) error {
	if customMapping, ok := customMappingFuncs[toField.Name]; ok {
		switch customFunc := customMapping.(type) {
		case func(from interface{}, to interface{}) interface{}:
			mappedValue := customFunc(fromFieldValue.Interface(), toFieldValue.Interface())
			if toFieldValue.CanSet() {
				toFieldValue.Set(reflect.ValueOf(mappedValue))
			}
		case string:
			mappedField := fromFieldValue.FieldByName(customFunc)
			if mappedField.IsValid() && toFieldValue.CanSet() {
				toFieldValue.Set(mappedField)
			}
		default:
			return fmt.Errorf("unsupported custom mapping type: %T", customMapping)
		}
	}
	return nil
}

func (m *Mapper) AddCustomMapping(fromFieldName string, toFieldName string, customFuncOrField interface{}) {
	if m.customMappings[fromFieldName] == nil {
		m.customMappings[fromFieldName] = make(map[string]interface{})
	}
	m.customMappings[fromFieldName][toFieldName] = customFuncOrField
}

func (m *Mapper) mapField(fromValue, toValue reflect.Value, i int) error {
	fromField := fromValue.Type().Field(i)
	fromFieldValue := fromValue.Field(i)

	if fromField.Type.Kind() == reflect.Struct {
		return m.mapEmbeddedStruct(fromField, fromFieldValue, toValue)
	}

	toField, found := toValue.Type().FieldByName(fromField.Name)
	if !found {
		return nil
	}

	toFieldValue := toValue.FieldByName(toField.Name)

	if customMappingFuncs, ok := m.customMappings[fromField.Name]; ok {
		return m.mapWithCustomMapping(customMappingFuncs, fromField, toField, fromFieldValue, toFieldValue)
	} else {
		return m.mapWithDefaultMapping(fromFieldValue, toFieldValue)
	}
}

func (m *Mapper) mapEmbeddedStruct(fromField reflect.StructField, fromFieldValue, toValue reflect.Value) error {
	toField, found := toValue.Type().FieldByName(fromField.Name)
	if found {
		toFieldValue := toValue.FieldByName(toField.Name)
		return m.Map(fromFieldValue.Interface(), toFieldValue.Addr().Interface())
	}

	for j := 0; j < fromFieldValue.NumField(); j++ {
		innerFromField := fromFieldValue.Type().Field(j)
		innerFromFieldValue := fromFieldValue.Field(j)
		innerToField, found := toValue.Type().FieldByName(innerFromField.Name)
		if found {
			innerToFieldValue := toValue.FieldByName(innerToField.Name)
			if innerFromFieldValue.Type().AssignableTo(innerToFieldValue.Type()) && innerToFieldValue.CanSet() {
				innerToFieldValue.Set(innerFromFieldValue)
			}
		}
	}
	return nil
}

func (m *Mapper) mapWithDefaultMapping(fromFieldValue, toFieldValue reflect.Value) error {
	if fromFieldValue.Type().AssignableTo(toFieldValue.Type()) && toFieldValue.CanSet() {
		toFieldValue.Set(fromFieldValue)
	}
	return nil
}
func (m *Mapper) mapSlice(fromValue reflect.Value, toValue reflect.Value, fromElemType reflect.Type, toElemType reflect.Type) error {
	fromSlice := fromValue.Interface()
	toSlice := reflect.MakeSlice(toValue.Type(), fromValue.Len(), fromValue.Cap())
	for i := 0; i < fromValue.Len(); i++ {
		fromElemValue := reflect.ValueOf(fromSlice).Index(i)
		toElemValue := reflect.New(toElemType).Elem()

		err := m.Map(fromElemValue.Interface(), toElemValue.Addr().Interface())
		if err != nil {
			return err
		}

		toSlice.Index(i).Set(toElemValue)
	}

	toValue.Set(toSlice)

	return nil
}

func (m *Mapper) mapArray(fromValue reflect.Value, toValue reflect.Value, fromElemType reflect.Type, toElemType reflect.Type) error {
	if fromValue.Len() != toValue.Len() {
		return fmt.Errorf("cannot map array of different length")
	}

	for i := 0; i < fromValue.Len(); i++ {
		fromElemValue := fromValue.Index(i)
		toElemValue := toValue.Index(i)

		err := m.Map(fromElemValue.Interface(), toElemValue.Addr().Interface())
		if err != nil {
			return err
		}
	}

	return nil
}

func (m *Mapper) extractReflectedValue(i interface{}) reflect.Value {
	val := reflect.ValueOf(i)

	// If pointer get the underlying element≤
	if val.Kind() == reflect.Ptr {
		val = val.Elem()
	}

	return val
}

func NewMapper() *Mapper {
	return &Mapper{customMappings: make(map[string]map[string]interface{})}
}
