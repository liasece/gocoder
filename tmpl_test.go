package gocoder

import (
	"reflect"
	"testing"
	"text/template"
)

func TestTemplate(t *testing.T) {
	type args struct {
		tmplContent string
		env         interface{}
		fn          template.FuncMap
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{
			name: "OK",
			args: args{
				tmplContent: `testing: {{ .Name }}`,
				env: struct {
					Name string
				}{
					"test",
				},
			},
			want: `testing: test`,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Template(tt.args.tmplContent, tt.args.env, nil)
			if (err != nil) != tt.wantErr {
				t.Errorf("Template() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			gotStr := ToCode(got)
			// if err != nil {
			// 	t.Errorf("Template() error = %v, wantErr %v", err, tt.wantErr)
			// 	return
			// }
			if !reflect.DeepEqual(gotStr, tt.want) {
				t.Errorf("Template() = %v, want %v", gotStr, tt.want)
			}
		})
	}
}
