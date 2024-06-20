package utils

import (
	"net/http"
)

func GetPostFormData(r *http.Request) (map[string]string, error) {
	args := map[string]string{}
	err := r.ParseForm()
	if err != nil {
		return args, err
	}
	for key, value := range r.Form {
		if len(value) > 0 {
			args[key] = value[0]
		} else {
			args[key] = ""
		}
	}
	return args, nil
}
