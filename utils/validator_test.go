package utils

import (
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestValidate(t *testing.T) {
	type request struct {
		URL      string    `json:"url" validate:"required,url"`
		Custom   string    `json:"custom,omitempty" validate:"omitempty,alphanum,max=16"`
		ExpireAt time.Time `json:"expire_at,omitempty" validate:"omitempty,gte"`
	}

	type args struct {
		r *request
	}
	tests := []struct {
		name          string
		args          args
		wantErrorsLen int
	}{
		{
			name: "OK",
			args: args{&request{
				URL:      "https://github.com/txfs19260817/url-shortener",
				Custom:   "shortener",
				ExpireAt: time.Now().Add(1 * time.Hour),
			}},
			wantErrorsLen: 0,
		},
		{
			name: "url",
			args: args{&request{
				URL:      "txfs19260817/url-shortener",
				Custom:   "shortener",
				ExpireAt: time.Now().Add(1 * time.Hour),
			}},
			wantErrorsLen: 1,
		},
		{
			name: "alphanum",
			args: args{&request{
				URL:      "https://github.com/txfs19260817/url-shortener",
				Custom:   "url-shortener",
				ExpireAt: time.Now().Add(1 * time.Hour),
			}},
			wantErrorsLen: 1,
		},
		{
			name: "max",
			args: args{&request{
				URL:      "https://github.com/txfs19260817/url-shortener",
				Custom:   "shortenershortenershortener",
				ExpireAt: time.Now().Add(1 * time.Hour),
			}},
			wantErrorsLen: 1,
		},
		{
			name: "gte",
			args: args{&request{
				URL:      "https://github.com/txfs19260817/url-shortener",
				Custom:   "shortener",
				ExpireAt: time.Now().Add(-1 * time.Hour),
			}},
			wantErrorsLen: 1,
		},
		{
			name: "omitempty",
			args: args{&request{
				URL: "https://github.com/txfs19260817/url-shortener",
			}},
			wantErrorsLen: 0,
		},
		{
			name: "required",
			args: args{&request{
				Custom: "shortener",
			}},
			wantErrorsLen: 1,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			errors := Validate(tt.args.r)
			assert.Len(t, errors, tt.wantErrorsLen)
			for _, errorResponse := range errors {
				t.Log(errorResponse)
			}
		})
	}
}

func TestCheckNoLoopRisk(t *testing.T) {
	type args struct {
		url  string
		host string
		port int
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "http://localhost:3000",
			args: args{
				url:  "http://localhost:3000",
				host: "localhost",
				port: 3000,
			},
			want: false,
		},
		{
			name: "localhost:3000",
			args: args{
				url:  "http://localhost:3000",
				host: "localhost",
				port: 3000,
			},
			want: false,
		},
		{
			name: "localhost",
			args: args{
				url:  "localhost",
				host: "localhost",
			},
			want: false,
		},
		{
			name: "localhost:80",
			args: args{
				url:  "http://localhost:3000",
				host: "localhost",
				port: 80,
			},
			want: false,
		},
		{
			name: "abc:3000",
			args: args{
				url:  "https://abc:3000",
				host: "localhost",
				port: 80,
			},
			want: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := CheckNoLoopRisk(tt.args.url, tt.args.host, tt.args.port); got != tt.want {
				t.Errorf("CheckNoLoopRisk() = %v, want %v", got, tt.want)
			}
		})
	}
}
