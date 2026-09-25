package service

import (
	"context"
	"encoding/json"
	"fmt"
	"net/url"
	"os"
	"path/filepath"
	"testing"

	"github.com/taurushq-io/taurus-protect-sdk/taurus-protect-sdk-go/pkg/protect/model"
)

// enforcedInRulesVectorsRelPath points at the vectors every SDK test suite consumes, so the four
// SDKs read the enforced-in-rules flags the same way: false only where the endpoint computes them,
// nil where it does not.
const enforcedInRulesVectorsRelPath = "../../../../scripts/resources/enforced-in-rules-vectors.json"

// enforcedInRulesVectorCount moves with the vector file, so a dropped vector fails instead of passing.
const enforcedInRulesVectorCount = 7

type enforcedInRulesVector struct {
	Name      string `json:"name"`
	Operation string `json:"operation"`
	Request   struct {
		ID string `json:"id"`
	} `json:"request"`
	Reply    json.RawMessage `json:"reply"`
	Expected struct {
		Query map[string]string `json:"query"`
		Users []struct {
			EnforcedInRules          *bool   `json:"enforcedInRules"`
			PublicKeyEnforcedInRules *bool   `json:"publicKeyEnforcedInRules"`
			GroupsEnforcedInRules    []*bool `json:"groupsEnforcedInRules"`
		} `json:"users"`
		Groups []struct {
			EnforcedInRules      *bool   `json:"enforcedInRules"`
			UsersEnforcedInRules []*bool `json:"usersEnforcedInRules"`
		} `json:"groups"`
	} `json:"expected"`
}

func TestEnforcedInRulesVectors(t *testing.T) {
	raw, err := os.ReadFile(filepath.Clean(enforcedInRulesVectorsRelPath))
	if err != nil {
		t.Fatalf("cannot read shared vectors file: %v", err)
	}
	var vectors []enforcedInRulesVector
	if err := json.Unmarshal(raw, &vectors); err != nil {
		t.Fatalf("cannot parse shared vectors file: %v", err)
	}
	if len(vectors) != enforcedInRulesVectorCount {
		t.Fatalf("shared vectors file has %d vectors, want %d", len(vectors), enforcedInRulesVectorCount)
	}

	for _, v := range vectors {
		t.Run(v.Name, func(t *testing.T) {
			ctx := context.Background()
			switch v.Operation {
			case "getMe":
				client, srv := newToleranceClient(t, "/api/rest/v1/users/me", v.Reply)
				user, err := NewUserService(client).GetMe(ctx)
				if err != nil {
					t.Fatalf("GetMe failed: %v", err)
				}
				assertEnforcedInRulesQuery(t, srv, v.Expected.Query)
				assertEnforcedInRulesUsers(t, v, []*model.User{user})
			case "getUser":
				client, _ := newToleranceClient(t, "/api/rest/v1/users/"+v.Request.ID, v.Reply)
				user, err := NewUserService(client).GetUser(ctx, v.Request.ID)
				if err != nil {
					t.Fatalf("GetUser failed: %v", err)
				}
				assertEnforcedInRulesUsers(t, v, []*model.User{user})
			case "listUsers":
				client, _ := newToleranceClient(t, "/api/rest/v1/users", v.Reply)
				result, err := NewUserService(client).ListUsers(ctx, nil)
				if err != nil {
					t.Fatalf("ListUsers failed: %v", err)
				}
				assertEnforcedInRulesUsers(t, v, result.Users)
			case "listGroups":
				client, _ := newToleranceClient(t, "/api/rest/v1/groups", v.Reply)
				result, err := NewGroupService(client).ListGroups(ctx, nil)
				if err != nil {
					t.Fatalf("ListGroups failed: %v", err)
				}
				assertEnforcedInRulesGroups(t, v, result.Groups)
			default:
				t.Fatalf("no runner registered for operation %s", v.Operation)
			}
		})
	}
}

func assertEnforcedInRulesQuery(t *testing.T, srv *toleranceServer, want map[string]string) {
	t.Helper()
	srv.mu.Lock()
	defer srv.mu.Unlock()
	if len(srv.requests) != 1 {
		t.Fatalf("got %d requests, want 1", len(srv.requests))
	}
	query, err := url.ParseQuery(srv.requests[0].query)
	if err != nil {
		t.Fatalf("cannot parse request query %q: %v", srv.requests[0].query, err)
	}
	for key, value := range want {
		if got := query.Get(key); got != value {
			t.Fatalf("request query %s = %q, want %q (query %q)", key, got, value, srv.requests[0].query)
		}
	}
}

func assertEnforcedInRulesUsers(t *testing.T, v enforcedInRulesVector, users []*model.User) {
	t.Helper()
	if len(users) != len(v.Expected.Users) {
		t.Fatalf("got %d users, want %d", len(users), len(v.Expected.Users))
	}
	for i, want := range v.Expected.Users {
		got := users[i]
		assertFlag(t, fmt.Sprintf("users[%d].EnforcedInRules", i), got.EnforcedInRules, want.EnforcedInRules)
		assertFlag(t, fmt.Sprintf("users[%d].PublicKeyEnforcedInRules", i), got.PublicKeyEnforcedInRules, want.PublicKeyEnforcedInRules)
		if len(got.Groups) != len(want.GroupsEnforcedInRules) {
			t.Fatalf("users[%d] has %d groups, want %d", i, len(got.Groups), len(want.GroupsEnforcedInRules))
		}
		for j, wantGroup := range want.GroupsEnforcedInRules {
			assertFlag(t, fmt.Sprintf("users[%d].Groups[%d].EnforcedInRules", i, j), got.Groups[j].EnforcedInRules, wantGroup)
		}
	}
}

func assertEnforcedInRulesGroups(t *testing.T, v enforcedInRulesVector, groups []*model.Group) {
	t.Helper()
	if len(groups) != len(v.Expected.Groups) {
		t.Fatalf("got %d groups, want %d", len(groups), len(v.Expected.Groups))
	}
	for i, want := range v.Expected.Groups {
		got := groups[i]
		assertFlag(t, fmt.Sprintf("groups[%d].EnforcedInRules", i), got.EnforcedInRules, want.EnforcedInRules)
		if len(got.Users) != len(want.UsersEnforcedInRules) {
			t.Fatalf("groups[%d] has %d users, want %d", i, len(got.Users), len(want.UsersEnforcedInRules))
		}
		for j, wantUser := range want.UsersEnforcedInRules {
			assertFlag(t, fmt.Sprintf("groups[%d].Users[%d].EnforcedInRules", i, j), got.Users[j].EnforcedInRules, wantUser)
		}
	}
}

func assertFlag(t *testing.T, name string, got, want *bool) {
	t.Helper()
	if (got == nil) != (want == nil) || (got != nil && *got != *want) {
		t.Errorf("%s = %s, want %s", name, flagString(got), flagString(want))
	}
}

func flagString(b *bool) string {
	if b == nil {
		return "nil"
	}
	return fmt.Sprint(*b)
}
