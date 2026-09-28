package services

import (
	"encoding/json"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/rojgarsetu/backend/internal/db"
)

func TestPublicUserOmitsPasswordHash(t *testing.T) {
	user := &db.User{
		ID:           uuid.New(),
		Name:         "Candidate",
		Email:        "candidate@example.invalid",
		PasswordHash: "must-never-be-serialized",
		Role:         "candidate",
	}

	encoded, err := json.Marshal(toPublicUser(user))
	if err != nil {
		t.Fatalf("marshal public user: %v", err)
	}
	response := string(encoded)
	if strings.Contains(response, "password_hash") || strings.Contains(response, user.PasswordHash) {
		t.Fatalf("public user response exposed credential data: %s", response)
	}
}

func TestAuthUserResponseContainsNoCredentials(t *testing.T) {
	user := &db.User{
		ID:           uuid.New(),
		Name:         "Candidate",
		Email:        "candidate@example.invalid",
		PasswordHash: "must-never-be-serialized",
		Role:         "candidate",
	}

	encoded, err := json.Marshal(authUserResponse(user))
	if err != nil {
		t.Fatalf("marshal auth response: %v", err)
	}
	response := string(encoded)
	for _, forbidden := range []string{"password_hash", user.PasswordHash, `"token"`} {
		if strings.Contains(response, forbidden) {
			t.Fatalf("auth response exposed %q: %s", forbidden, response)
		}
	}
}

func TestDatabaseUserOmitsPasswordHash(t *testing.T) {
	user := db.User{PasswordHash: "must-never-be-serialized"}
	encoded, err := json.Marshal(user)
	if err != nil {
		t.Fatalf("marshal database user: %v", err)
	}
	if strings.Contains(string(encoded), "password_hash") || strings.Contains(string(encoded), user.PasswordHash) {
		t.Fatalf("database user serialization exposed credential data: %s", encoded)
	}
}

func TestAuthCookiesAreSecureInProduction(t *testing.T) {
	t.Setenv("ENVIRONMENT", "production")
	t.Setenv("COOKIE_SECURE", "false")
	gin.SetMode(gin.TestMode)
	ctx, _ := gin.CreateTestContext(httptest.NewRecorder())
	ctx.Request = httptest.NewRequest("POST", "/login", nil)

	setAuthCookies(ctx, "access", "refresh")
	cookies := ctx.Writer.Header().Values("Set-Cookie")
	if len(cookies) == 0 {
		t.Fatal("expected auth cookies to be set")
	}
	for _, cookie := range cookies {
		if !strings.Contains(cookie, "; Secure") || !strings.Contains(cookie, "; HttpOnly") {
			t.Errorf("production auth cookie lacks Secure or HttpOnly: %s", cookie)
		}
	}
}

func TestAuthCookiesRespectSecureFlagOutsideProduction(t *testing.T) {
	t.Setenv("ENVIRONMENT", "development")
	t.Setenv("COOKIE_SECURE", "true")
	gin.SetMode(gin.TestMode)
	ctx, _ := gin.CreateTestContext(httptest.NewRecorder())
	ctx.Request = httptest.NewRequest("POST", "/login", nil)

	setAuthCookies(ctx, "access", "refresh")
	if got := os.Getenv("COOKIE_SECURE"); got != "true" {
		t.Fatalf("test environment changed unexpectedly: %q", got)
	}
	for _, cookie := range ctx.Writer.Header().Values("Set-Cookie") {
		if !strings.Contains(cookie, "; Secure") {
			t.Errorf("COOKIE_SECURE=true was not applied: %s", cookie)
		}
	}
}
