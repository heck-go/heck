package heck

import (
	"errors"
	"reflect"
)

type Variabler struct {
	variablesById map[string]interface{}

	// TODO by type and id
	variablesByType map[reflect.Type]map[string]interface{}
}

func NewVariabler() *Variabler {
	return &Variabler{
		variablesById: map[string]interface{}{},
		variablesByType: map[reflect.Type]map[string]interface{}{},
	}
}

func (self *Variabler) SetVariable(key string, variable interface{}) {
	self.variablesById[key] = variable
	m, ok := self.variablesByType[reflect.TypeOf(variable)]
	if !ok {
		m = map[string]interface{}{}
		self.variablesByType[reflect.TypeOf(variable)] = m
	}
	m[key] = variable
}

func (self *Variabler) GetVariableById(key string, variable interface{}) error {
	value, ok := self.variablesById[key]
	if !ok {
		return errors.New("variable with given id not found")
	}
	
	if reflect.TypeOf(value) == reflect.TypeOf(variable).Elem() {
		reflect.ValueOf(variable).Elem().Set(reflect.ValueOf(value))
	} else if reflect.TypeOf(value).Implements(reflect.TypeOf(variable).Elem()) {
		reflect.ValueOf(variable).Elem().Set(reflect.ValueOf(value))
	} else {
		return errors.New("mismatch of variable type")
	}

	return nil
}

func (self *Variabler) GetVariableByType(variable interface{}, key string) error {
	if reflect.TypeOf(variable).Kind() != reflect.Ptr {
		return errors.New("variable parameter must be a pointer")
	}

	t := reflect.TypeOf(variable).Elem()

	m, ok := self.variablesByType[t]
	if !ok {
		return errors.New("variable with given type not found")
	}
	
	if len(key) == 0 {
		// Use any (in this case) value
		for _, first := range m {
			reflect.ValueOf(variable).Elem().Set(reflect.ValueOf(first))
			return nil
		}
	}
	
	v, ok1 := m[key]
	if !ok1 {
		return errors.New("variable with given type and id not found")
	}
	reflect.ValueOf(variable).Elem().Set(reflect.ValueOf(v))
	return nil
}
