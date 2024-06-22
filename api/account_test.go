package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
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
	server := NewTestServer(store)
	recorder := httptest.NewRecorder()

	url := fmt.Sprintf("/api/accounts/%d", account.ID)
	req := httptest.NewRequest(http.MethodGet, url, nil)
	server.router.ServeHTTP(recorder, req)
	require.Equal(t, http.StatusOK, recorder.Code)
	requireResBodyMatchAccount(t, recorder.Body, *account)
}

func randomAccount() *db.Account {
	return &db.Account{
		ID:       int64(util.RandomMoney()),
		Balance:  int64(util.RandomMoney()),
		Owner:    util.RandomOwner(),
		Currency: util.RandomCurrencies(),
	}
}

func requireResBodyMatchAccount(t *testing.T, body *bytes.Buffer, account db.Account) {
	data, err := io.ReadAll(body)
	require.NoError(t, err)

	var gotAccount db.Account
	json.Unmarshal(data, &gotAccount)
	require.Equal(t, account, gotAccount)

}
