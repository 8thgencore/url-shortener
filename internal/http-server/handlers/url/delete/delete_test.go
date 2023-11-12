package delete_test

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
	resp2 "url-shortener/internal/lib/api/response"
	"url-shortener/internal/lib/logger/handlers/slogdiscard"
	"url-shortener/internal/storage"

	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"url-shortener/internal/http-server/handlers/url/delete"
	"url-shortener/internal/http-server/handlers/url/delete/mocks"
)

func TestDeleteHandler(t *testing.T) {
	cases := []struct {
		name      string
		uri       string
		respError string
		mockError error
		code      int
	}{
		{
			name: "Success",
			uri:  "/url/testalias",
			code: http.StatusOK,
		},
		{
			name:      "Invalid Alias",
			uri:       "/url/" + url.QueryEscape("!@#$%"),
			respError: "url alias not valid",
			code:      http.StatusOK,
		},
		{
			name: "Omitted Alias",
			uri:  "/url/",
			respError: "Don't delete it! We do not check this message because it is generated" +
				" by the router. But it must not be empty for the test to work correctly.",
			code: http.StatusNotFound,
		},
		{
			name:      "Alias Not Found",
			uri:       "/url/testalias",
			respError: "url alias not found",
			mockError: storage.ErrURLNotFound,
			code:      http.StatusOK,
		},
		{
			name:      "Delete Error",
			uri:       "/url/testalias",
			respError: "failed to delete url",
			mockError: errors.New("unexpected error"),
			code:      http.StatusOK,
		},
	}

	for _, tc := range cases {
		tc := tc

		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			urlDeleterMock := mocks.NewURLDeleter(t)

			if tc.respError == "" || tc.mockError != nil {
				urlDeleterMock.On("DeleteURL", mock.AnythingOfType("string")).
					Return(tc.mockError).
					Once()
			}

			handler := chi.NewRouter()
			handler.Delete("/url/{alias}", delete.New(slogdiscard.NewDiscardLogger(), urlDeleterMock))

			req, err := http.NewRequest(http.MethodDelete, tc.uri, nil)
			require.NoError(t, err)

			rr := httptest.NewRecorder()
			handler.ServeHTTP(rr, req)

			require.Equal(t, tc.code, rr.Code)

			if rr.Code == http.StatusOK {
				body := rr.Body.String()

				var resp resp2.Response

				require.NoError(t, json.Unmarshal([]byte(body), &resp))

				require.Equal(t, tc.respError, resp.Error)
			}
		})
	}
}
