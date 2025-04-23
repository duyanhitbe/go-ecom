package handlers

import (
	"bytes"
	"encoding/json"
	"errors"
	"github.com/duyanhitbe/go-ecom/internal/dto"
	"github.com/duyanhitbe/go-ecom/internal/repositories"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestUserHandler_CreateUser(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mock := repositories.NewMockQuerier(ctrl)
	handler := NewUserHandler(mock)

	tests := []testcase[*dto.CreateUserRequest]{
		{
			name: "Successful user creation",
			input: &dto.CreateUserRequest{
				Username: "test",
				Password: "test",
			},
			mockBehavior: func() {
				mock.EXPECT().
					CreateUser(gomock.Any(), &repositories.CreateUserParams{
						Username: "test",
						Password: "test",
					}).
					Return(&repositories.User{
						ID:        uuid.New(),
						Username:  "test",
						Password:  "test",
						CreatedAt: time.Now(),
						UpdatedAt: time.Now(),
					}, nil)
			},
			expectedCode: http.StatusCreated,
		},
		{
			name: "Validation error",
			input: &dto.CreateUserRequest{
				Username: "",
				Password: "",
			},
			mockBehavior: func() {
				// No interaction with mock for validation errors.
			},
			expectedCode: http.StatusBadRequest,
		},
		{
			name: "Database error",
			input: &dto.CreateUserRequest{
				Username: "test",
				Password: "test",
			},
			mockBehavior: func() {
				mock.EXPECT().
					CreateUser(gomock.Any(), &repositories.CreateUserParams{
						Username: "test",
						Password: "test",
					}).
					Return(nil, errors.New("database error"))
			},
			expectedCode: http.StatusInternalServerError,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			tc.mockBehavior()

			w := httptest.NewRecorder()
			body, _ := json.Marshal(tc.input)
			r := httptest.NewRequest(http.MethodPost, "/users", bytes.NewBuffer(body))
			r.Header.Set("Content-Type", "application/json")

			handler.CreateUser(w, r)

			require.Equal(t, tc.expectedCode, w.Code)

			if tc.expectedCode == http.StatusCreated {
				var response dto.Response[repositories.User]
				err := json.NewDecoder(w.Body).Decode(&response)
				require.NoError(t, err)
				require.Equal(t, http.StatusCreated, *response.StatusCode)
				require.Equal(t, http.StatusText(http.StatusCreated), *response.Message)
				require.Nil(t, response.Error)
				require.Nil(t, response.Errors)
				require.True(t, *response.Success)
				require.NotEmpty(t, response.Data.ID)
				require.NotEmpty(t, response.Data.CreatedAt)
				require.NotEmpty(t, response.Data.UpdatedAt)
				require.Equal(t, tc.input.Username, response.Data.Username)
				require.Equal(t, tc.input.Password, response.Data.Password)
			}
		})
	}
}
