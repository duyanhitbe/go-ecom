package repositories

import (
	"github.com/go-faker/faker/v4"
	"github.com/stretchr/testify/require"
	"testing"
)

func createTestUser(t *testing.T) *User {
	user, err := testQueries.CreateUser(ctx, &CreateUserParams{
		Username: faker.Username(),
		Password: faker.Password(),
	})

	require.NoError(t, err)
	require.NotEmpty(t, user.ID)
	require.NotEmpty(t, user.Username)
	require.NotEmpty(t, user.Password)
	require.NotEmpty(t, user.CreatedAt)
	require.NotEmpty(t, user.UpdatedAt)

	return user
}

func TestQueries_CreateUser(t *testing.T) {
	createTestUser(t)
}

func TestQueries_FindUser(t *testing.T) {
	createTestUser(t)

	users, err := testQueries.FindUser(ctx, &FindUserParams{
		Offset: 0,
		Limit:  10,
	})

	require.NoError(t, err)
	require.NotEmpty(t, users)
	for _, user := range users {
		require.NotEmpty(t, user.ID)
		require.NotEmpty(t, user.Username)
		require.NotEmpty(t, user.Password)
		require.NotEmpty(t, user.CreatedAt)
		require.NotEmpty(t, user.UpdatedAt)
	}
}

func TestQueries_CountUser(t *testing.T) {
	createTestUser(t)
	count, err := testQueries.CountUser(ctx)
	require.NoError(t, err)
	require.NotZero(t, count)
}

func TestQueries_FindOneUserByUsername(t *testing.T) {
	user := createTestUser(t)

	foundUser, err := testQueries.FindOneUserByUsername(ctx, user.Username)
	require.NoError(t, err)
	require.NotEmpty(t, foundUser.ID)
	require.NotEmpty(t, foundUser.Username)
	require.NotEmpty(t, foundUser.Password)
	require.NotEmpty(t, foundUser.CreatedAt)
	require.NotEmpty(t, foundUser.UpdatedAt)
	require.Equal(t, user.Username, foundUser.Username)
}
