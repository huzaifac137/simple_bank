package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	db "github.com/huzaifac137/simple_bank/db/sqlc"
	testdb "github.com/huzaifac137/simple_bank/db/test_db"
	"github.com/stretchr/testify/require"
)

func TestCreateUserSuccess(t *testing.T) {

	conn, err := testdb.ConnectTestDb()
	require.NoError(t, err)
	store := db.NewStore(conn)

	server := NewTestServer(store)

	testReq := &createUserRequest{
		Username: "test99",
		Password: "123456",
		Email:    "test99@gamil.com",
		FullName: "test test 99",
	}

	jsonReq, err := json.Marshal(testReq)
	require.NoError(t, err)
	if err != nil {
		t.Fatal(err)
	}

	reader := bytes.NewReader(jsonReq)

	req := httptest.NewRequest("POST", "/api/users/", reader)
	rec := httptest.NewRecorder()
	res := rec.Result()
	resJson, err := io.ReadAll(res.Body)
	require.NoError(t, err)

	server.router.ServeHTTP(rec, req)

	var fromResToStruct *CreateUserResponse
	json.Unmarshal(resJson, fromResToStruct)

	require.Equal(t, http.StatusOK, res.StatusCode)
	fmt.Printf("fromResToStruct: %v", fromResToStruct)
	// require.Equal(t, testReq.Username, fromResToStruct.Username)
	// require.Equal(t, testReq.Email, fromResToStruct.Email)
	// require.Equal(t, testReq.FullName, fromResToStruct.FullName)
	// require.NotZero(t, fromResToStruct.CreatedAt)
}

func TestLoginUserSuccess(t *testing.T) {

	conn, err := testdb.ConnectTestDb()
	require.NoError(t, err)
	store := db.NewStore(conn)

	server := NewTestServer(store)

	testReq := &LoginUserRequest{
		Username: "test99",
		Password: "123456",
	}

	jsonReq, err := json.Marshal(testReq)
	require.NoError(t, err)

	reader := bytes.NewReader(jsonReq)

	req := httptest.NewRequest("POST", "/api/users/login", reader)
	rec := httptest.NewRecorder()

	server.router.ServeHTTP(rec, req)

	res := rec.Result()
	resBody, err := io.ReadAll(res.Body)
	require.NoError(t, err)

	var fromResToStruct LoginUserResponse
	json.Unmarshal(resBody, &fromResToStruct)

	require.Equal(t, http.StatusOK, res.StatusCode)
	require.Equal(t, testReq.Username, fromResToStruct.Username)
	require.NotZero(t, fromResToStruct.CreatedAt)
}
