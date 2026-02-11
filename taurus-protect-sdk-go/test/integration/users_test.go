package integration

import (
	"context"
	"testing"

	"github.com/taurushq-io/taurus-protect-sdk/taurus-protect-sdk-go/pkg/protect/model"
)

func TestIntegration_GetMe(t *testing.T) {
	skipIfNotIntegration(t)
	client := getTestClient(t)
	defer client.Close()

	ctx := context.Background()
	user, err := client.Users().GetMe(ctx)
	if err != nil {
		t.Fatalf("GetMe() error = %v", err)
	}

	t.Logf("Current user:")
	t.Logf("  ID: %s", user.ID)
	t.Logf("  Username: %s", user.Username)
	t.Logf("  Email: %s", user.Email)
	t.Logf("  Status: %s", user.Status)
}

func TestIntegration_ListUsers(t *testing.T) {
	skipIfNotIntegration(t)
	client := getTestClient(t)
	defer client.Close()

	ctx := context.Background()
	result, err := client.Users().ListUsers(ctx, &model.ListUsersOptions{
		Limit: 10,
	})
	if err != nil {
		t.Fatalf("ListUsers() error = %v", err)
	}

	t.Logf("Found %d users", len(result.Users))
	for _, u := range result.Users {
		t.Logf("User: ID=%s, Username=%s, Email=%s", u.ID, u.Username, u.Email)
	}
}
