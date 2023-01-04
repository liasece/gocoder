package ast

import (
	"fmt"
	"testing"

	"github.com/liasece/gocoder"
)

func TestGetTypeFromSource(t *testing.T) {
	res, err := GetTypeFromSource("../test/source/struct.go", "BigStruct")
	if err != nil {
		t.Error(err)
	}
	str, err := gocoder.WriteToFileStr(gocoder.NewCode().C(res))
	if err != nil {
		t.Error(err)
		return
	}
	fmt.Println(str)
	t.Error("finish")
	nextType := res.FieldByName("Next").GetType().Elem()
	nextType.SetInReference(false)
	fmt.Printf("finish: %+v\n", nextType)
}
