package response

import (
	"reflect"
	"testing"
)

func TestIsCodeError(t *testing.T) {
	type args struct {
		errMsg string
	}
	tests := []struct {
		name  string
		args  args
		want  bool
		want1 CodeError
	}{
		{
			name:  "非正常",
			args:  args{errMsg: "xxx"},
			want:  false,
			want1: CodeError{},
		},
		{
			name:  "正常",
			args:  args{errMsg: "ErrCode:101，ErrMsg:session expired"},
			want:  true,
			want1: CodeError{IRet: 101, SMsg: "session expired"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1 := IsCodeError(tt.args.errMsg)
			if got != tt.want {
				t.Errorf("IsCodeError() got = %v, want %v", got, tt.want)
			}
			if !reflect.DeepEqual(got1, tt.want1) {
				t.Errorf("IsCodeError() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}
