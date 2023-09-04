package create_user

import (
	"github.com/AlexandrLitkevich/qwery/internal/http-server/handlers/user/create_user/mocks"
	"github.com/AlexandrLitkevich/qwery/internal/lib/logger/handlers/slogdiscard"
	"github.com/stretchr/testify/mock"
	"testing"
)

func TestCreateUserHandler(t *testing.T) {
	type userInfo struct {
		name     string
		age      int
		position string
	}
	tests := []struct {
		name       string
		wantStatus bool
		wantError  string
		userInfo
	}{
		{
			name:       "first case",
			wantStatus: true,
			userInfo: userInfo{
				age:      11,
				name:     "Pops",
				position: "PM",
			},
		},
	}
	for _, tt := range tests {
		//tt := tt
		t.Run(tt.name, func(t *testing.T) {
			mockMethod := mocks.NewCreateUserProvider(t)
			mockLogger := slogdiscard.NewDiscardLogger()

			if tt.wantStatus == true {
				//Настраиваем в соответсвии с кейсами
				mockMethod.On("CreateUser", tt.userInfo, mock.AnythingOfType("bool")).Return(true, tt.wantError).Once()
			}
			got := New(mockLogger, mockMethod)
			t.Log(got)
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
