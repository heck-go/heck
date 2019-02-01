package heck

import (
	"fmt"
	"testing"
)

func TestVariabler_GetVariableById(t *testing.T) {
	v := NewVariabler()
	v.SetVariable("number", 5)
	var num int
	var err error
	err = v.GetVariableById("number", &num)
	if err != nil {
		t.FailNow()
	}
	fmt.Println(num)
	
	num = 0
	err = v.GetVariableByType(&num, "number")
	if err != nil {
		t.FailNow()
	}
	fmt.Println(num)
	
	num = 0
	err = v.GetVariableByType(&num, "")
	if err != nil {
		t.FailNow()
	}
	fmt.Println(num)
}
