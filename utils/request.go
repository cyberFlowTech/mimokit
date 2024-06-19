package utils

import (
	"net/http"
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
			args[kvs[0]] = kvs[1]
		}
	}
	return args, nil
}
