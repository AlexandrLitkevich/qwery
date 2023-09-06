package request

import (
	"bytes"
	"fmt"
	resp "github.com/AlexandrLitkevich/qwery/internal/lib/api/response"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	"github.com/stretchr/testify/require"
	"net/http"
	"net/http/httptest"
	"testing"
	"unsafe"
)

type Request struct {
	Name string
	Age  int
}

func TestRequestSize(t *testing.T) {
	SizeByteHandle := func(w http.ResponseWriter, r *http.Request) {
		var req Request

		err := render.DecodeJSON(r.Body, &req)
		t.Log("this req str 25", req)
		t.Log("this error", err)
		if err != nil {
			t.Log("this block error")
			render.JSON(w, r, resp.Error("failed to decode request"))
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		w.WriteHeader(http.StatusOK)
		render.JSON(w, r, "byte success")
	}

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
			router.Use(RequestSize(1)) //This this function
			router.Post("/size", SizeByteHandle)
			ts := httptest.NewServer(router)
			defer ts.Close()

			handler := http.HandlerFunc(SizeByteHandle)

			input := fmt.Sprintf(`{"name": "%s", "desc": "%s"}`, tt.name, tt.desc)

			size := unsafe.Sizeof(input)
			t.Log("This size input", size)

			req, err := http.NewRequest(http.MethodPost, ts.URL+"/size", bytes.NewReader([]byte(input)))
			require.NoError(t, err)

			rr := httptest.NewRecorder()
			handler.ServeHTTP(rr, req)

			body := rr.Body.String()
			t.Log("This code", rr.Code)
			t.Log("This code", body)

			//var resp2 struct{}
			//
			//err = json.Unmarshal([]byte(body), &resp2)
			//t.Log(err)
			//
			//require.NoError(t, err)

		})
	}
}
