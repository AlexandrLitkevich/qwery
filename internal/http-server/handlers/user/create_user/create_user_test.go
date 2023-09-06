package create_user

import (
	"bytes"
	"encoding/json"
	"errors"
	"github.com/AlexandrLitkevich/qwery/internal/http-server/handlers/user"
	"github.com/AlexandrLitkevich/qwery/internal/http-server/handlers/user/create_user/mocks"
	"github.com/AlexandrLitkevich/qwery/internal/lib/logger/handlers/slogdiscard"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestCreateUserHandler(t *testing.T) {

	tests := []struct {
		name       string
		wantStatus bool
		wantError  string
		userInfo   user.Request
	}{
		{
			name:       "first case",
			wantStatus: true,
			userInfo: user.Request{
				Age:      30,
				Name:     "Pops-test",
				Position: "Team lead",
			},
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			mockMethod := mocks.NewCreateUserProvider(t)
			mockLogger := slogdiscard.NewDiscardLogger()

			if tt.wantStatus == true {
				//Настраиваем в соответсвии с кейсами
				mockMethod.On("CreateUser", tt.userInfo, mock.AnythingOfType("string")).Return(true, errors.New(tt.wantError)).Once()
			}
			handler := New(mockLogger, mockMethod)

			//body request
			in, err := json.Marshal(tt.userInfo)
			require.NoError(t, err)

			req, err := http.NewRequest(http.MethodPost, "/user", bytes.NewReader([]byte(in)))
			require.NoError(t, err)

			rr := httptest.NewRecorder()
			handler.ServeHTTP(rr, req)

			require.Equal(t, rr.Code, http.StatusOK)

			body := rr.Body.String()

			var resp user.Response

			require.NoError(t, json.Unmarshal([]byte(body), &resp))
		})
	}
}

//func Test_responseOK(t *testing.T) {
//	type args struct {
//		w           http.ResponseWriter
//		r           *http.Request
//		createdUser user.User
//	}
//	tests := []struct {
//		name string
//		args args
//	}{
//		// TODO: Add test cases.
//	}
//	for _, tt := range tests {
//		t.Run(tt.name, func(t *testing.T) {
//			responseOK(tt.args.w, tt.args.r, tt.args.createdUser)
//		})
//	}
//}
