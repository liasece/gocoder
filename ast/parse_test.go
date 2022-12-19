package ast

import (
	"testing"

	source_test "github.com/liasece/gocoder/test/source"
)

func TestGetGoFileFullPackage(t *testing.T) {
	type args struct {
		filePath string
	}
	tests := []struct {
		name      string
		args      args
		wantPkg   string
		wantAlias string
	}{
		{
			name: "test1",
			args: args{
				filePath: "./astCoder.go",
			},
			wantPkg:   "github.com/liasece/gocoder/ast",
			wantAlias: "ast",
		},
		{
			name: "test1",
			args: args{
				filePath: "../test/source/struct.go",
			},
			wantPkg:   "github.com/liasece/gocoder/test/source",
			wantAlias: "source_test",
		},
	}
	var _ source_test.BigStruct
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotPkg, gotAlias := GetGoFileFullPackage(tt.args.filePath)
			if gotPkg != tt.wantPkg {
				t.Errorf("GetGoFileFullPackage() gotPkg = %v, want %v", gotPkg, tt.wantPkg)
			}
			if gotAlias != tt.wantAlias {
				t.Errorf("GetGoFileFullPackage() gotAlias = %v, want %v", gotAlias, tt.wantAlias)
			}
		})
	}
}
