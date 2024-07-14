package crypt_test

import (
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/kosalnik/metrics/internal/crypt/mock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/kosalnik/metrics/internal/crypt"
)

func TestCipherMiddleware(t *testing.T) {
	cases := []struct {
		name            string
		encodedBody     string
		decodeError     error
		wantDecodedBody string
		wantErr         bool
	}{
		{
			name:            "Success",
			encodedBody:     "very crypted body",
			decodeError:     nil,
			wantDecodedBody: "Decoded body",
			wantErr:         false,
		},
		{
			name:        "Fail decoding",
			encodedBody: "very crypted body",
			decodeError: errors.New("random error"),
			wantErr:     true,
		},
	}
	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			decoderMock := mock.NewMockDecoder(ctrl)
			decoderMock.EXPECT().Decode([]byte(tt.encodedBody)).Times(1).Return([]byte(tt.wantDecodedBody), tt.decodeError)
			mw := crypt.CipherMiddleware(decoderMock)
			mux := http.NewServeMux()
			mux.HandleFunc(`/`, func(writer http.ResponseWriter, request *http.Request) {
				defer require.NoError(t, request.Body.Close())
				got, err := io.ReadAll(request.Body)
				require.NoError(t, err)
				require.Equal(t, tt.wantDecodedBody, string(got))
			})
			h := mw(mux)
			r := httptest.NewRequest(http.MethodPost, `/`, strings.NewReader(tt.encodedBody))
			w := httptest.NewRecorder()
			h.ServeHTTP(w, r)
			res := w.Result()
			defer assert.NoError(t, res.Body.Close())
			if tt.wantErr {
				require.Equal(t, http.StatusBadRequest, res.StatusCode)

			} else {
				require.Equal(t, http.StatusOK, res.StatusCode)
			}
		})
	}
}
