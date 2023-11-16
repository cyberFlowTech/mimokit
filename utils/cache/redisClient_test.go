package cache

import (
	"testing"
	"time"
)

func TestRedisClient_Setex(t *testing.T) {

	type args struct {
		key    string
		value  string
		expire time.Duration
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
		{
			name: "test_Setex",
			args: args{
				key:    "test_setex_key",
				value:  "test_setex_value",
				expire: time.Second * 3600 * 24,
			},
			wantErr: false,
		},
	}

	c := NewRedisClient([]string{"127.0.0.1:6379"})

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			if err := c.Setex(tt.args.key, tt.args.value, tt.args.expire); (err != nil) != tt.wantErr {
				t.Errorf("Setex() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestRedisClient_Get(t *testing.T) {

	type args struct {
		key string
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		// TODO: Add test cases.
		{
			name: "test_Setex",
			args: args{
				key: "test_setex_key",
			},
			want:    "test_setex_value",
			wantErr: false,
		},
	}

	c := NewRedisClient([]string{"127.0.0.1:6379"})

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := c.Get(tt.args.key)
			if (err != nil) != tt.wantErr {
				t.Errorf("Get() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("Get() got = %v, want %v", got, tt.want)
			}
		})
	}
}
