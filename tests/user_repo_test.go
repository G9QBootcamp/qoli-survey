package test

import (
	"context"
	"testing"

	"github.com/G9QBootcamp/qoli-survey/internal/user/dto"
	"github.com/G9QBootcamp/qoli-survey/internal/user/models"
	"github.com/G9QBootcamp/qoli-survey/internal/util"
)

func TestCrudUserRepository(t *testing.T) {

	nationalId := util.GenerateNumericString(10)
	ctx := context.Background()

	user, err := testUserRepo.CreateUser(ctx, &models.User{FirstName: "test_name", LastName: "last_name", RoleID: 1, NationalID: nationalId, City: "tehran", Email: nationalId + "@gmail.com", PasswordHash: "kjfqwnion4"})
	if err != nil {
		t.Fatalf("failed to create user in user repository: %s", err.Error())

	}
	users, err := testUserRepo.GetUsers(ctx, dto.UserFilters{Name: "test_name", Email: nationalId + "@gmail.com", NationalID: nationalId, City: "tehran", Limit: 1})
	if err != nil {
		t.Fatalf("failed to get users in user repository: %s", err.Error())
	}

	if len(users) < 1 {
		t.Fatal("get users filter did not work properly because no user returned")
	}

	filtersUser := users[0]

	if filtersUser.NationalID != user.NationalID || filtersUser.City != user.City || filtersUser.Email != user.Email || filtersUser.LastName != user.LastName || filtersUser.FirstName != user.FirstName {
		t.Fatal("filtered user did not match with created user there are some problem with create or get Users function")
	}
	_, err = testUserRepo.GetUserByID(ctx, user.ID)

	if err != nil {
		t.Fatalf("failed to get user by id in user repository: %s", err.Error())
	}
	err = testUserRepo.DeleteUser(ctx, user.ID)

	if err != nil {
		t.Fatalf("failed to delete user by id in user repository: %s", err.Error())
	}

	t.Log("Crud Operations in User Repository executed successfully")

}
