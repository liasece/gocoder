package ast

import (
	"fmt"
	"testing"

	"github.com/liasece/gocoder"
)

func TestGetInterfaceFromSource(t *testing.T) {
	res, err := GetInterfaceFromSource("../test/source/struct.go", "IBigStruct")
	if err != nil {
		t.Error(err)
	}
	str, err := gocoder.WriteToFileStr(gocoder.NewCode().C(res))
	if err != nil {
		t.Error(err)
		return
	}
	fmt.Println(str)
	t.Error("")
}
