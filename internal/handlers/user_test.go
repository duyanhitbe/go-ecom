package handlers

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/duyanhitbe/go-ecom/internal/dto"
	"github.com/duyanhitbe/go-ecom/internal/repositories"
	"github.com/duyanhitbe/go-ecom/pkg/utils"
	mock "github.com/duyanhitbe/go-ecom/test/mocks"
	"github.com/go-faker/faker/v4"
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

	mockQuerier := mock.NewMockQuerier(ctrl)
	mockHash := mock.NewMockHash(ctrl)
	handler := NewUserHandler(mockQuerier, mockHash)

	fakeUsername := faker.Username()
	fakePassword := faker.Password()
	hashedPassword := faker.Password()

	tests := []testcase[*dto.CreateUserRequest]{
		{
			name: "Successful user creation",
			input: &dto.CreateUserRequest{
				Username: fakeUsername,
				Password: fakePassword,
			},
			mockBehavior: func() {
				mockQuerier.EXPECT().
					FindOneUserByUsername(gomock.Any(), fakeUsername).
					Times(1).
					Return(&repositories.User{
						ID:        uuid.New(),
						CreatedAt: time.Now(),
						UpdatedAt: time.Now(),
					}, nil)
				mockHash.EXPECT().
					Hash(fakePassword).
					Times(1).
					Return(hashedPassword, nil)
				mockQuerier.EXPECT().
					CreateUser(gomock.Any(), &repositories.CreateUserParams{
						Username: fakeUsername,
						Password: hashedPassword,
					}).
					Times(1).
					Return(&repositories.User{
						ID:        uuid.New(),
						Username:  fakeUsername,
						Password:  hashedPassword,
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

			},
			expectedCode: http.StatusBadRequest,
		},
		{
			name: "Database error from CreateUser",
			input: &dto.CreateUserRequest{
				Username: fakeUsername,
				Password: fakePassword,
			},
			mockBehavior: func() {
				mockHash.EXPECT().
					Hash(fakePassword).
					Times(1).
					Return(hashedPassword, nil)
				mockQuerier.EXPECT().
					FindOneUserByUsername(gomock.Any(), fakeUsername).
					Times(1).
					Return(&repositories.User{}, nil)
				mockQuerier.EXPECT().
					CreateUser(gomock.Any(), &repositories.CreateUserParams{
						Username: fakeUsername,
						Password: hashedPassword,
					}).
					Times(1).
					Return(nil, errors.New("database error"))
			},
			expectedCode: http.StatusInternalServerError,
		},
		{
			name: "Database error from FindOneUserByUsername",
			input: &dto.CreateUserRequest{
				Username: fakeUsername,
				Password: fakePassword,
			},
			mockBehavior: func() {
				mockQuerier.EXPECT().
					FindOneUserByUsername(gomock.Any(), fakeUsername).
					Times(1).
					Return(nil, errors.New("database error"))
			},
			expectedCode: http.StatusInternalServerError,
		},
		{
			name: "Username already exists",
			input: &dto.CreateUserRequest{
				Username: fakeUsername,
				Password: fakePassword,
			},
			mockBehavior: func() {
				mockQuerier.EXPECT().
					FindOneUserByUsername(gomock.Any(), fakeUsername).
					Times(1).
					Return(&repositories.User{
						ID:        uuid.New(),
						Username:  fakeUsername,
						Password:  hashedPassword,
						CreatedAt: time.Now(),
						UpdatedAt: time.Now(),
					}, nil)
			},
			expectedCode: http.StatusBadRequest,
		},
		{
			name: "Hash password error",
			input: &dto.CreateUserRequest{
				Username: fakeUsername,
				Password: fakePassword,
			},
			mockBehavior: func() {
				mockQuerier.EXPECT().
					FindOneUserByUsername(gomock.Any(), fakeUsername).
					Times(1).
					Return(&repositories.User{
						ID:        uuid.New(),
						CreatedAt: time.Now(),
						UpdatedAt: time.Now(),
					}, nil)
				mockHash.EXPECT().
					Hash(fakePassword).
					Times(1).
					Return("", errors.New("hash error"))
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
				var response dto.Response[dto.CreateUserResponse]
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
			}
		})
	}
}

func TestUserHandler_FindUser(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockQuerier := mock.NewMockQuerier(ctrl)
	mockHash := mock.NewMockHash(ctrl)
	handler := NewUserHandler(mockQuerier, mockHash)

	tests := []testcase[*dto.FindUserRequest]{
		{
			name: "Successful user finding",
			path: fmt.Sprintf("/users?page=%d&per_page=%d", 1, 10),
			mockBehavior: func() {
				pg := int32(1)
				perPg := int32(10)
				_, limit, offset := utils.GetPaginationMeta(&pg, &perPg)
				mockQuerier.EXPECT().
					FindUser(gomock.Any(), &repositories.FindUserParams{
						Offset: offset,
						Limit:  limit,
					}).
					Times(1).
					Return([]*repositories.User{
						{
							ID:        uuid.New(),
							Username:  "test",
							Password:  "test",
							CreatedAt: time.Now(),
							UpdatedAt: time.Now(),
						},
					}, nil)
				mockQuerier.EXPECT().
					CountUser(gomock.Any()).
					Times(1).
					Return(int32(1), nil)
			},
			expectedCode: http.StatusOK,
		},
		{
			name: "Validation error - negative page",
			path: "/users?page=-1&per_page=10",
			mockBehavior: func() {
				// No mock calls expected for validation error
			},
			expectedCode: http.StatusBadRequest,
		},
		{
			name: "Validation error - negative per_page",
			path: "/users?page=1&per_page=-10",
			mockBehavior: func() {
				// No mock calls expected for validation error
			},
			expectedCode: http.StatusBadRequest,
		},
		{
			name: "Database error - FindUser fails",
			path: fmt.Sprintf("/users?page=%d&per_page=%d", 1, 10),
			mockBehavior: func() {
				mockQuerier.EXPECT().
					FindUser(gomock.Any(), &repositories.FindUserParams{
						Offset: 0,
						Limit:  10,
					}).
					Times(1).
					Return(nil, errors.New("database error"))
			},
			expectedCode: http.StatusInternalServerError,
		},
		{
			name: "Database error - CountUser fails",
			path: fmt.Sprintf("/users?page=%d&per_page=%d", 1, 10),
			mockBehavior: func() {
				mockQuerier.EXPECT().
					FindUser(gomock.Any(), &repositories.FindUserParams{
						Offset: 0,
						Limit:  10,
					}).
					Times(1).
					Return([]*repositories.User{}, nil)
				mockQuerier.EXPECT().
					CountUser(gomock.Any()).
					Times(1).
					Return(int32(0), errors.New("database error"))
			},
			expectedCode: http.StatusInternalServerError,
		},
		{
			name: "Invalid page parameter format",
			path: "/users?page=invalid&per_page=10",
			mockBehavior: func() {
				// No mock calls expected for validation error
			},
			expectedCode: http.StatusBadRequest,
		},
		{
			name: "Invalid per_page parameter format",
			path: "/users?page=1&per_page=invalid",
			mockBehavior: func() {
				// No mock calls expected for validation error
			},
			expectedCode: http.StatusBadRequest,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			tc.mockBehavior()
			w := httptest.NewRecorder()
			r := httptest.NewRequest(http.MethodGet, tc.path, bytes.NewBuffer([]byte("")))
			handler.FindUser(w, r)

			require.Equal(t, tc.expectedCode, w.Code)

			if tc.expectedCode == http.StatusOK {
				var response dto.Response[[]*dto.FindUserResponse]
				err := json.NewDecoder(w.Body).Decode(&response)
				require.NoError(t, err)
				require.Equal(t, http.StatusOK, *response.StatusCode)
				require.Equal(t, http.StatusText(http.StatusOK), *response.Message)
				require.Nil(t, response.Error)
				require.Nil(t, response.Errors)
				require.True(t, *response.Success)
				require.NotNil(t, response.Data)
				require.NotNil(t, response.Meta)

				for _, user := range *response.Data {
					require.NotEmpty(t, user.ID)
					require.NotEmpty(t, user.UpdatedAt)
					require.NotEmpty(t, user.CreatedAt)
					require.Equal(t, "test", user.Username)
				}
			}
		})
	}
}
