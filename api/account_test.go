package api

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/golang/mock/gomock"
	mockdb "github.com/huzaifac137/simple_bank/db/mock"
	db "github.com/huzaifac137/simple_bank/db/sqlc"
	"github.com/huzaifac137/simple_bank/util"
	"github.com/stretchr/testify/require"
)

func TestGetAccountApi(t *testing.T) {
	account := randomAccount()
	ctrl := gomock.NewController(t)
	store := mockdb.NewMockStore(ctrl)
	store.EXPECT().GetAccount(gomock.Any(), account.ID).Times(1).Return(*account, nil)

	// satrt http mock server
	server := NewServer(store)
	recorder := httptest.NewRecorder()

	url := fmt.Sprintf("/api/accounts/%d", account.ID)
	req := httptest.NewRequest(http.MethodGet, url, nil)
	server.router.ServeHTTP(recorder, req)
	require.Equal(t, http.StatusOK, recorder.Code)

}

func randomAccount() *db.Account {
	return &db.Account{
		ID:       int64(util.RandomMoney()),
		Balance:  int64(util.RandomMoney()),
		Owner:    util.RandomOwner(),
		Currency: util.RandomCurrencies(),
	}
}
