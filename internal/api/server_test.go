package api

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	"github.com/golang/mock/gomock"
	modelUser "github.com/igortoigildin/go-rewards-app/internal/entities/user"
	modelWithdrawal "github.com/igortoigildin/go-rewards-app/internal/entities/withdrawal"
	"github.com/igortoigildin/go-rewards-app/mocks"
)

func Test_registerUserHandler(t *testing.T) {
	input := struct {
		Login    string `json:"login"`
		Password string `json:"password"`
	}{
		Login:    "dummyLogin",
		Password: "abc",
	}
	jsonUserData, _ := json.Marshal(input)
	req, err := http.NewRequest("POST", "/api/user/register", bytes.NewBuffer(jsonUserData))
	if err != nil {
		t.Fatal(err)
	}
	user := modelUser.User{}
	_ = user.Password.Set(input.Password)

	ctrl := gomock.NewController(t)
	u := mocks.NewMockUserService(ctrl)
	tok := mocks.NewMockTokenService(ctrl)

	u.EXPECT().Create(gomock.Any(), gomock.Any()).Times(1).Return(nil)
	tok.EXPECT().NewToken(gomock.Any(), gomock.Any(), gomock.Any()).Times(1).Return(&user.Token, nil)

	handler := http.HandlerFunc(registerUserHandler(u, tok))
	srv := httptest.NewServer(handler)
	defer srv.Close()

	rr := httptest.NewRecorder()

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}
}

func Test_balanceHandler(t *testing.T) {
	req, err := http.NewRequest("GET", "/api/user/balance", nil)
	if err != nil {
		t.Fatal(err)
	}
	user := modelUser.User{
		UserID: int64(1),
	}
	req = contextSetUser(req, &user)
	req.AddCookie(&http.Cookie{Name: "token", Value: "dummy_cookie"})
	ctrl := gomock.NewController(t)
	u := mocks.NewMockUserService(ctrl)
	w := mocks.NewMockWithdrawalService(ctrl)
	currentBalance := float64(500)
	withdrawn := []modelWithdrawal.Withdrawal{}

	u.EXPECT().Balance(gomock.Any(), user.UserID).Times(1).Return(currentBalance, nil)
	w.EXPECT().WithdrawalsForUser(gomock.Any(), user.UserID).Times(1).Return(withdrawn, nil)

	handler := http.HandlerFunc(balanceHandler(u, w))
	srv := httptest.NewServer(handler)
	defer srv.Close()

	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	data := struct {
		Current   float64 `json:"current"`
		Withdrawn float64 `json:"withdrawn"`
	}{}
	want := struct {
		Current   float64 `json:"current"`
		Withdrawn float64 `json:"withdrawn"`
	}{
		Current:   float64(500),
		Withdrawn: float64(0),
	}
	err = json.NewDecoder(rr.Body).Decode(&data)
	if err != nil {
		t.Fatal(err)
	}

	if got := data; !reflect.DeepEqual(got, want) {
		t.Errorf("balanceHandler() = %v, want %v", got, want)
	}
}
