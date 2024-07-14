package crypt_test

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/kosalnik/metrics/internal/crypt"
)

var headerName = "HashSHA256"

func TestExtractSign(t *testing.T) {
	tests := map[string]struct {
		r    func() *http.Request
		want string
	}{
		"Success": {
			r: func() *http.Request {
				r := httptest.NewRequest(`POST`, `/`, strings.NewReader(`Hello`))
				r.Header.Add(headerName, "aa00")
				return r
			},
			want: "aa00",
		},
		"Negative": {
			r: func() *http.Request {
				r := httptest.NewRequest(`POST`, `/`, strings.NewReader(`Hello`))
				r.Header.Add("Tra", "aa00")
				return r
			},
			want: "",
		},
	}
	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			if got := crypt.ExtractSign(tt.r()); got != tt.want {
				t.Errorf("ExtractSign() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetSign(t *testing.T) {
	data := []byte(`Hello`)
	key := []byte(`Secret`)
	want := `6cd4180752f6880f553f2c73c89efe222166924f7bb9707c6240e4b88de77122`
	if got := crypt.GetSign(data, key); got != want {
		t.Errorf("GetSign() = %v, want %v", got, want)
	}
}

func TestToSignRequest(t *testing.T) {
	r := httptest.NewRequest(http.MethodPut, `/asd`, strings.NewReader("a"))
	want := "41234a5b"
	crypt.ToSignRequest(r, want)
	got := r.Header.Get(headerName)
	if got != want {
		t.Errorf("ToSignRequest(), got %v, want %v", got, want)
	}
}

func TestVerifySign(t *testing.T) {
	tests := map[string]struct {
		data []byte
		sign string
		key  []byte
		want bool
	}{
		"Success": {
			data: []byte(`Hello`),
			sign: `6cd4180752f6880f553f2c73c89efe222166924f7bb9707c6240e4b88de77122`,
			key:  []byte(`Secret`),
			want: true,
		},
		"Success From Lecture": {
			data: []byte("Видишь гофера? Нет. И я нет. А он есть."),
			sign: `2fd903d51a40a74cc3f79caa861cdc4e2a7edafaa6383b4975d840ebd2200c63`,
			key:  []byte("secretkeytratata"),
			want: true,
		},
		"WrongKey": {
			data: []byte(`Hello`),
			sign: `6cd4180752f6880f553f2c73c89efe222166924f7bb9707c6240e4b88de77122`,
			key:  []byte(`WrongSecret`),
			want: false,
		},
		"WrongData": {
			data: []byte(`Hello, World`),
			sign: `6cd4180752f6880f553f2c73c89efe222166924f7bb9707c6240e4b88de77122`,
			key:  []byte(`Secret`),
			want: false,
		},
		"WrongSign": {
			data: []byte(`Hello`),
			sign: `981298319823`,
			key:  []byte(`Secret`),
			want: false,
		},
	}
	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			if got := crypt.VerifySign(tt.data, tt.sign, tt.key); got != tt.want {
				t.Errorf("VerifySign() = %v, want %v", got, tt.want)
			}
		})
	}
}
