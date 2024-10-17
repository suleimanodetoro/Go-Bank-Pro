package api

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/require"
	mockdb "github.com/suleimanodetoro/Go-Bank-Pro/db/mock"
	db "github.com/suleimanodetoro/Go-Bank-Pro/db/sqlc"
	"github.com/suleimanodetoro/Go-Bank-Pro/db/util"
	"go.uber.org/mock/gomock"
)

// TestGetAccountAPI tests the GetAccount API endpoint.
func TestGetAccountAPI(t *testing.T) {
	// Create a random account to be used in test cases
	account := randomAccount()

	// Define test cases using a table-driven approach
	testCases := []struct {
		name          string                                                  // Name of the test case
		accountID     int64                                                   // ID of the account to retrieve
		buildStubs    func(store *mockdb.MockStore)                           // Function to set up mock behavior
		checkResponse func(t *testing.T, recorder *httptest.ResponseRecorder) // Function to verify the response
	}{
		{
			name:      "OK",       // Test case where everything works as expected
			accountID: account.ID, // Use the ID of the random account
			buildStubs: func(store *mockdb.MockStore) {
				// Set up expected behavior on the mock store
				// Expect GetAccount to be called with any context and the specific account ID
				// Return the account and no error
				store.EXPECT().
					GetAccount(gomock.Any(), gomock.Eq(account.ID)).
					Times(1).
					Return(account, nil)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				// Check that the status code is 200 OK
				require.Equal(t, http.StatusOK, recorder.Code)
				// Check that the response body matches the account
				requireBodyMatchAccount(t, recorder.Body, account)
			},
		},
		{
			name:      "NotFound", // Test case where account is not found
			accountID: account.ID, // Use the ID of the random account
			buildStubs: func(store *mockdb.MockStore) {
				// Set up expected behavior on the mock store
				// Expect GetAccount to be called with any context and the specific account ID
				// Return an empty account and sql.ErrNoRows to simulate a "not found" scenario
				store.EXPECT().
					GetAccount(gomock.Any(), gomock.Eq(account.ID)).
					Times(1).
					Return(db.Account{}, sql.ErrNoRows)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				// Check that the status code is 404 Not Found
				require.Equal(t, http.StatusNotFound, recorder.Code)
			},
		},
		// TODO: add more test cases for different scenarios (e.g., NotFound, BadRequest)
	}

	// Loop over each test case and run it as a subtest
	for i := range testCases {
		tc := testCases[i]

		// Run each test case as a subtest using t.Run
		t.Run(tc.name, func(t *testing.T) {
			// Create a new mock controller and defer its finish
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			// Create a new mock store
			store := mockdb.NewMockStore(ctrl)
			// Build stubs (set up expected behavior) for the mock store
			tc.buildStubs(store)

			// Start a new server with the mock store
			server := NewServer(store)
			// Create a response recorder to capture the response
			recorder := httptest.NewRecorder()

			// Construct the request URL
			url := fmt.Sprintf("/accounts/%d", tc.accountID)
			// Create a new HTTP GET request
			request, err := http.NewRequest(http.MethodGet, url, nil)
			require.NoError(t, err)

			// Serve the HTTP request
			server.router.ServeHTTP(recorder, request)
			// Check the response using the provided function
			tc.checkResponse(t, recorder)
		})
	}
}

// randomAccount creates a random account for testing purposes
func randomAccount() db.Account {
	return db.Account{
		ID:       util.RandomInt(1, 1000),    // Random account ID
		Owner:    util.RandomOwner(),         // Random account owner
		Balance:  util.RandomInt(0, 1000000), // Random balance
		Currency: util.RandomCurrency(),      // Random currency
	}
}

// requireBodyMatchAccount checks that the response body matches the expected account
func requireBodyMatchAccount(t *testing.T, body *bytes.Buffer, account db.Account) {
	// Read all the data from the response body
	data, err := io.ReadAll(body)
	require.NoError(t, err)

	// Unmarshal the JSON data into a db.Account object
	var gotAccount db.Account
	err = json.Unmarshal(data, &gotAccount)
	require.NoError(t, err)

	// Check that the account matches what we expect
	require.Equal(t, account, gotAccount)
}
