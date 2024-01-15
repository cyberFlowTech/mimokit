package lan

import (
	"github.com/nicksnyder/go-i18n/v2/i18n"
	"golang.org/x/text/language"
	"reflect"
	"testing"
)

func TestTranslate_Trans(t1 *testing.T) {
	type fields struct {
		bundle    *i18n.Bundle
		localizer map[string]*i18n.Localizer
		lanList   []string
	}
	type args struct {
		lan  string
		key  string
		args []interface{}
	}
	InitTranslate()
	obj := GetTranslate()
	tests := []struct {
		name   string
		fields fields
		args   args
		want   string
	}{
		{
			"test",
			fields{
				bundle:    obj.bundle,
				localizer: obj.localizer,
				lanList:   obj.lanList,
			},
			args{
				lan:  "en",
				key:  "1000049",
				args: []interface{}{1},
			},
			"Please enter an introduction within 1 words",
		},
	}
	for _, tt := range tests {
		t1.Run(tt.name, func(t1 *testing.T) {
			t := &Translate{
				bundle:    tt.fields.bundle,
				localizer: tt.fields.localizer,
				lanList:   tt.fields.lanList,
			}
			if got := t.Trans(tt.args.lan, tt.args.key, tt.args.args...); got != tt.want {
				t1.Errorf("Trans() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestTranslate_getLanTag(t1 *testing.T) {
	type fields struct {
		bundle    *i18n.Bundle
		localizer map[string]*i18n.Localizer
		lanList   []string
	}
	type args struct {
		lan string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   language.Tag
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t1.Run(tt.name, func(t1 *testing.T) {
			t := &Translate{
				bundle:    tt.fields.bundle,
				localizer: tt.fields.localizer,
				lanList:   tt.fields.lanList,
			}
			if got := t.getLanTag(tt.args.lan); !reflect.DeepEqual(got, tt.want) {
				t1.Errorf("getLanTag() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestTranslate_loadLan(t1 *testing.T) {
	type fields struct {
		bundle    *i18n.Bundle
		localizer map[string]*i18n.Localizer
		lanList   []string
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t1.Run(tt.name, func(t1 *testing.T) {
			t := &Translate{
				bundle:    tt.fields.bundle,
				localizer: tt.fields.localizer,
				lanList:   tt.fields.lanList,
			}
			if err := t.loadLan(); (err != nil) != tt.wantErr {
				t1.Errorf("loadLan() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestGetTranslate(t *testing.T) {
	tests := []struct {
		name string
		want *Translate
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := GetTranslate(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetTranslate() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestInitTranslate(t *testing.T) {
	tests := []struct {
		name string
	}{
		{
			"test",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			InitTranslate()
		})
	}
}

func TestTranslate_Trans1(t1 *testing.T) {
	type fields struct {
		bundle    *i18n.Bundle
		localizer map[string]*i18n.Localizer
		lanList   []string
	}
	type args struct {
		lan  string
		key  string
		args []interface{}
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t1.Run(tt.name, func(t1 *testing.T) {
			t := &Translate{
				bundle:    tt.fields.bundle,
				localizer: tt.fields.localizer,
				lanList:   tt.fields.lanList,
			}
			if got := t.Trans(tt.args.lan, tt.args.key, tt.args.args...); got != tt.want {
				t1.Errorf("Trans() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestTranslate_getLanTag1(t1 *testing.T) {
	type fields struct {
		bundle    *i18n.Bundle
		localizer map[string]*i18n.Localizer
		lanList   []string
	}
	type args struct {
		lan string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   language.Tag
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t1.Run(tt.name, func(t1 *testing.T) {
			t := &Translate{
				bundle:    tt.fields.bundle,
				localizer: tt.fields.localizer,
				lanList:   tt.fields.lanList,
			}
			if got := t.getLanTag(tt.args.lan); !reflect.DeepEqual(got, tt.want) {
				t1.Errorf("getLanTag() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestTranslate_getLocalizer(t1 *testing.T) {
	type fields struct {
		bundle    *i18n.Bundle
		localizer map[string]*i18n.Localizer
		lanList   []string
	}
	type args struct {
		lan string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   *i18n.Localizer
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t1.Run(tt.name, func(t1 *testing.T) {
			t := &Translate{
				bundle:    tt.fields.bundle,
				localizer: tt.fields.localizer,
				lanList:   tt.fields.lanList,
			}
			if got := t.getLocalizer(tt.args.lan); !reflect.DeepEqual(got, tt.want) {
				t1.Errorf("getLocalizer() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestTranslate_loadLan1(t1 *testing.T) {
	type fields struct {
		bundle    *i18n.Bundle
		localizer map[string]*i18n.Localizer
		lanList   []string
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t1.Run(tt.name, func(t1 *testing.T) {
			t := &Translate{
				bundle:    tt.fields.bundle,
				localizer: tt.fields.localizer,
				lanList:   tt.fields.lanList,
			}
			if err := t.loadLan(); (err != nil) != tt.wantErr {
				t1.Errorf("loadLan() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestTranslate_newTranslator(t1 *testing.T) {
	type fields struct {
		bundle    *i18n.Bundle
		localizer map[string]*i18n.Localizer
		lanList   []string
	}
	tests := []struct {
		name   string
		fields fields
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t1.Run(tt.name, func(t1 *testing.T) {
			t := &Translate{
				bundle:    tt.fields.bundle,
				localizer: tt.fields.localizer,
				lanList:   tt.fields.lanList,
			}
			t.newTranslator()
		})
	}
}
