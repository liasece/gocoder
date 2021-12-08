package cde

import (
	"testing"

	"github.com/liasece/gocoder"
)

func TestLoadTypeFromSource(t *testing.T) {
	resType, err := LoadTypeFromSource("../test/source/struct.go", "BigStruct")
	if err != nil {
		t.Error(err)
	}
	t.Error(gocoder.WriteToFile("testOut_gen.go", gocoder.NewCode().C(Type(resType))))
	t.Error("Done")
}
