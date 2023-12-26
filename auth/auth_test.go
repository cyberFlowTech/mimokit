package auth

import (
	"bytes"
	"io"
	"net/http"
	"testing"
)

func TestAuth(t *testing.T) {
	type args struct {
		r *http.Request
	}
	tests := []struct {
		name    string
		args    args
		want    bool
		wantErr bool
	}{
		{
			name: "title不存在空字符参数",
			args: args{
				r: &http.Request{
					Form: map[string][]string{
						"title":        {"gggg"},
						"api":          {"i_1648893994"},
						"lan":          {"zh_TW"},
						"sessid":       {"f947ba7b430d8a8f448550908bc74bef"},
						"sign_time":    {"1701913845"},
						"uuid":         {"49F5E96FDC20E8AE0C7BC6786EAAFF7D"},
						"version":      {"1.2.10.3"},
						"version_code": {"144"},
						"sign":         {"d56999551d9f6d99c4e0e1eb60fc6cde"},
						"user_id":      {"843974"},
					},
					Body: io.NopCloser(bytes.NewBuffer([]byte(""))),
				},
			},
			want:    true,
			wantErr: false,
		},
		{
			name: "title存在空字符参数",
			args: args{
				r: &http.Request{
					Form: map[string][]string{
						"title":        {""},
						"api":          {"i_1648893994"},
						"lan":          {"zh_TW"},
						"sessid":       {"f947ba7b430d8a8f448550908bc74bef"},
						"sign_time":    {"1701913845"},
						"uuid":         {"49F5E96FDC20E8AE0C7BC6786EAAFF7D"},
						"version":      {"1.2.10.3"},
						"version_code": {"144"},
						"sign":         {"cf8a02797a6a6ec6c99aa27df9a573e2"},
						"user_id":      {"843974"},
					},
					Body: io.NopCloser(bytes.NewBuffer([]byte(""))),
				},
			},
			want:    true,
			wantErr: false,
		},
		{
			name: "正常签名",
			args: args{
				r: &http.Request{
					Form: map[string][]string{
						"api":          {"a_1680591588"},
						"lan":          {"zh_CN"},
						"page":         {"1"},
						"sessid":       {"268dfaa7d9f7d5d87247697f238744c9"},
						"sign_time":    {"1703505811"},
						"uuid":         {"dd021fe0-cafa-4f39-b057-035a2d80c921"},
						"version":      {"2.0.0"},
						"version_code": {"1"},
						"sign":         {"a947a88e91c1c7986a160e1c0536aecf"},
						"user_id":      {"844446"},
					},
					Body: io.NopCloser(bytes.NewBuffer([]byte(""))),
				},
			},
			want:    true,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Auth(tt.args.r)
			if (err != nil) != tt.wantErr {
				t.Errorf("Auth() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("Auth() got = %v, want %v", got, tt.want)
			}
		})
	}
}
