package lan

import (
	"github.com/nicksnyder/go-i18n/v2/i18n"
	"golang.org/x/text/language"
	"reflect"
	"testing"
)

func TestNewTranslate(t *testing.T) {
	tests := []struct {
		name    string
		want    *Translate
		wantErr bool
	}{
		{
			"test",
			&Translate{
				&i18n.Bundle{},
				make(map[string]*i18n.Localizer),
				make([]string, 0),
			},
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := NewTranslate()
			if (err != nil) != tt.wantErr {
				t.Errorf("NewTranslate() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			//if !reflect.DeepEqual(got, tt.want) {
			//	t.Errorf("NewTranslate() got = %v, want %v", got, tt.want)
			//}
		})
	}
}

func TestTranslate_NewTranslator(t1 *testing.T) {
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
			t.NewTranslator()
		})
	}
}

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
	obj, _ := NewTranslate()
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
				args: []interface{}{2},
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
