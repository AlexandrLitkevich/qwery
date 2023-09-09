package request_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/AlexandrLitkevich/qwery/internal/http-server/handlers/helpers"
	mwRSize "github.com/AlexandrLitkevich/qwery/internal/http-server/middleware/request"
	resp "github.com/AlexandrLitkevich/qwery/internal/lib/api/response"
	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/require"
	"net/http"
	"net/http/httptest"
	"testing"
	"unsafe"
)

func TestRequestSize(t *testing.T) {

	res := mwRSize.RequestSize(12)

	handler := http.HandlerFunc(helpers.SizeHandlerHelpers)

	r := res(handler)

	input := fmt.Sprint("lkdfgnjklhjfkldsahdfkl;haskljdfhljkasdhfk ")

	t.Log("This size input", unsafe.Sizeof(input))

	req, err := http.NewRequest(http.MethodPost, "/size", bytes.NewReader([]byte(input)))
	require.NoError(t, err)

	rr := httptest.NewRecorder()

	r.ServeHTTP(rr, req)

	t.Log("This second req", req.Body)
	var js resp.Response
	body := rr.Body.String()
	err = json.Unmarshal([]byte(body), &js)
	require.NoError(t, err, "Fail json.Unmarshal([]byte(body), &js)")
	t.Log("this rr", js.Error)

}

func TestRequestSizeWithNewServer(t *testing.T) {
	type args struct {
		bytes int64
	}
	tests := []struct {
		name string
		desc string
		args args
		want string
	}{
		{
			name: "first request",
			args: args{
				bytes: 1,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			router := chi.NewRouter()
			router.Use(mwRSize.RequestSize(tt.args.bytes))
			router.Post("/size", helpers.SizeHandlerHelpers)
			ts := httptest.NewServer(router)
			defer ts.Close()

			handler := http.HandlerFunc(helpers.SizeHandlerHelpers)

			input := fmt.Sprintf(`{"name": "%s", "desc": "%s"}`, tt.name, tt.desc)

			t.Log("This size input", unsafe.Sizeof(input))

			req, err := http.NewRequest(http.MethodPost, ts.URL+"/size", bytes.NewReader([]byte(input)))
			require.NoError(t, err)

			rr := httptest.NewRecorder()
			handler.ServeHTTP(rr, req)

			//body := rr.Body.String()
			t.Log("This code", rr.Code)
			//t.Log("This code", body)

			//var resp2 struct{}
			//
			//err = json.Unmarshal([]byte(body), &resp2)
			//t.Log(err)
			//
			//require.NoError(t, err)

		})
	}
}
