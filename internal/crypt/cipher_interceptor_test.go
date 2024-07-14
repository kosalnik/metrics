package crypt_test

import (
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/kosalnik/metrics/internal/crypt"
	"github.com/kosalnik/metrics/internal/crypt/mock"
	"github.com/stretchr/testify/require"
)

func TestCipherInterceptor(t *testing.T) {
	cases := []struct {
		name            string
		body            string
		encodeError     error
		wantEncodedBody string
		wantErr         bool
	}{
		{
			name:            "Success",
			body:            "very crypted body",
			encodeError:     nil,
			wantEncodedBody: "Decoded body",
			wantErr:         false,
		},
		{
			name:        "Fail decoding",
			body:        "very crypted body",
			encodeError: errors.New("random error"),
			wantErr:     true,
		},
	}
	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			encoderMock := mock.NewMockEncoder(ctrl)
			encoderMock.EXPECT().Encode([]byte(tt.body)).Times(1).Return([]byte(tt.wantEncodedBody), tt.encodeError)

			mux := http.NewServeMux()
			mux.HandleFunc(`/`, func(writer http.ResponseWriter, request *http.Request) {
				defer request.Body.Close()
				got, err := io.ReadAll(request.Body)
				require.NoError(t, err)
				require.Equal(t, tt.wantEncodedBody, string(got))
			})
			srv := httptest.NewServer(mux)

			client := srv.Client()
			oldTransport := srv.Client().Transport
			client.Transport = crypt.NewCipherInterceptor(encoderMock, oldTransport)

			resp, err := client.Post(srv.URL, "application/json", strings.NewReader(tt.body))
			if tt.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				require.NoError(t, resp.Body.Close())
				require.Equal(t, http.StatusOK, resp.StatusCode)
			}
		})
	}
}
