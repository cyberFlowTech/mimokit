package response

import (
	"google.golang.org/grpc/status"
	"net/http"
	"testing"
)

func TestHTTPResponse_JSON(t *testing.T) {
	type fields struct {
		ServerCommonErrorCode int
		TokenExpireErrorCode  int
		message               map[int]string
	}
	type args struct {
		r    *http.Request
		w    http.ResponseWriter
		resp interface{}
		err  error
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		{
			name: "grpc error",
			fields: fields{
				ServerCommonErrorCode: -1,
				TokenExpireErrorCode:  -100,
				message:               map[int]string{},
			},
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
						//"user_id":      {"844446"},
					},
				},
				//w: http.ResponseWriter(),
				resp: "",
				err:  status.Error(2, "ErrCode:5000016，ErrMsg:xxxxx"),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := &HTTPResponse{
				Config: Config{
					Trans:                 false,
					ServerCommonErrorCode: tt.fields.ServerCommonErrorCode,
					TokenExpireErrorCode:  tt.fields.TokenExpireErrorCode,
				},
				message: tt.fields.message,
			}
			h.JSON(tt.args.r, tt.args.w, tt.args.resp, tt.args.err)
		})
	}
}
