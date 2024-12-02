package test

import (
	"context"
	"testing"

	"github.com/G9QBootcamp/qoli-survey/internal/user/dto"
)

func TestGetUsersInUserService(t *testing.T) {
	ctx := context.Background()
	users, err := testUserService.GetUsers(ctx, dto.UserGetRequest{Name: "first_name"})
	if err != nil {
		t.Fatalf("error in get Users in user service :%s", err.Error())
	}
	if len(users) < 1 {
		t.Log("get users in user service executed successfully but no user returned in response")
		return
	}
	t.Log("get users in user service executed successfully")
}
