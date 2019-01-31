package heck

import (
	"strconv"
)

type CastableMap struct {
	data map[string]string
}

func NewCastableMap(data map[string]string) *CastableMap {
	return &CastableMap{
		data: data,
	}
}

func (self *CastableMap) Get(key string) (string, bool) {
	v, ok := self.data[key]
	return v, ok
}

func (self *CastableMap) Int(key string) (int, bool, bool) {
	v, ok := self.data[key]
	if !ok {
		return 0, false, false
	}
	ret, err := strconv.Atoi(v)
	if err != nil {
		return 0, false, true
	}
	return ret, true, true
}

func (self *CastableMap) Double(key string) (float64, bool, bool) {
	v, ok := self.data[key]
	if !ok {
		return 0, false, false
	}
	ret, err := strconv.ParseFloat(v, 64)
	if err != nil {
		return 0, false, true
	}
	return ret, true, true
}

func (self *CastableMap) Bool(key string) (bool, bool, bool) {
	v, ok := self.data[key]
	if !ok {
		return false, false, false
	}
	ret, err := strconv.ParseBool(v)
	if err != nil {
		return false, false, true
	}
	return ret, true, true
}

// TODO datetime

// TODO list
