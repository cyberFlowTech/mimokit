package utils

import (
	"net/http"
	"net/url"
	"strings"
)

func GetPostFormData(r *http.Request) (map[string]string, error) {
	args := map[string]string{}
	err := r.ParseForm()
	if err != nil {
		return args, err
	}
	encode := r.PostForm.Encode()
	if encode != "" {
		params := strings.Split(encode, "&")
		for k := range params {
			kvs := strings.Split(params[k], "=")
			if len(kvs) != 2 {
				continue
			}
			var value string
			value, err = url.PathUnescape(kvs[1])
			if err != nil {
				return args, err
			}
			args[kvs[0]] = value
		}
	}
	return args, nil
}
