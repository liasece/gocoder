package ast

import (
	"fmt"
	"testing"

	"github.com/liasece/gocoder"
)

func TestGetMethods(t *testing.T) {
	res, err := GetMethodsFromSource("../test/source/struct.go", "BigStruct")
	if err != nil {
		t.Error(err)
	}
	list := make([]gocoder.Codeable, 0, len(res))
	for _, v := range res {
		list = append(list, v)
	}
	str, err := gocoder.WriteToFileStr(gocoder.NewCode().C(list...))
	if err != nil {
		t.Error(err)
		return
	}
	fmt.Println(str)
	t.Error("")
}
