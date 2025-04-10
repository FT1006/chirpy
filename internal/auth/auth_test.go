package auth

import (
	"net/http"
	"strings"
	"testing"
	"time"

	"github.com/google/uuid"
)

func TestHashPassword(t *testing.T) {
	hash, err := HashPassword("password")
	if err != nil {
		t.Errorf("Error hashing password: %s", err)
	}
	if len(hash) == 0 {
		t.Errorf("Hash is empty")
	}
}

func TestCheckPasswordHash(t *testing.T) {
	hash, err := HashPassword("password")
	if err != nil {
		t.Errorf("Error hashing password: %s", err)
	}
	err = CheckPasswordHash(hash, "password")
	if err != nil {
		t.Errorf("Error checking password hash: %s", err)
	}
}

func TestMakeJWT(t *testing.T) {
	var testUserID uuid.UUID
	testUserID, err := uuid.Parse("b466f0d7-0059-487a-bc82-07067494b577")
	if err != nil {
		t.Errorf("Error parsing user ID: %s", err)
	}

	testTokenSecret := "test_secret"
	testExpiresIn := time.Hour

	if _, err := MakeJWT(testUserID, testTokenSecret, testExpiresIn); err != nil {
		t.Errorf("Error making token: %s", err)
	}
}

func TestValidateJWT(t *testing.T) {
	var testUserID uuid.UUID
	testUserID, err := uuid.Parse("b466f0d7-0059-487a-bc82-07067494b577")
	if err != nil {
		t.Errorf("Error parsing user ID: %s", err)
	}

	testTokenSecret := "test_secret"
	testExpiresIn := time.Hour

	testTokenString, err := MakeJWT(testUserID, testTokenSecret, testExpiresIn)
	if err != nil {
		t.Errorf("Error making token: %s", err)
	}

	if _, err := ValidateJWT(testTokenString, testTokenSecret); err != nil {
		t.Errorf("Error validating token: %s", err)
	}
}

func TestValidateJWTTokenExpired(t *testing.T) {
	var testUserID uuid.UUID
	testUserID, err := uuid.Parse("b466f0d7-0059-487a-bc82-07067494b577")
	if err != nil {
		t.Errorf("Error parsing user ID: %s", err)
	}

	testTokenSecret := "test_secret"
	testExpiresIn := time.Second

	testTokenString, err := MakeJWT(testUserID, testTokenSecret, testExpiresIn)
	if err != nil {
		t.Errorf("Error making token: %s", err)
	}

	time.Sleep(testExpiresIn * 3)

	_, err = ValidateJWT(testTokenString, testTokenSecret)
	if err == nil {
		t.Errorf("Expected token validation to fail due to expiration, but it succeeded.")
	} else if !strings.Contains(err.Error(), "expired") {
		// Check if the error message describes expiration.
		t.Errorf("Expected expiration error, but got: %s", err)
	}
}

func TestValidateJWTWrongSecret(t *testing.T) {
	var testUserID uuid.UUID
	testUserID, err := uuid.Parse("b466f0d7-0059-487a-bc82-07067494b577")
	if err != nil {
		t.Errorf("Error parsing user ID: %s", err)
	}

	testTokenSecret := "test_secret"
	wrongTokenSecret := "wrong_secret"
	testExpiresIn := time.Second

	testTokenString, err := MakeJWT(testUserID, wrongTokenSecret, testExpiresIn)
	if err != nil {
		t.Errorf("Error making token: %s", err)
	}

	_, err = ValidateJWT(testTokenString, testTokenSecret)
	if err == nil {
		t.Errorf("Expected token validation to fail due to wrong secret, but it succeeded.")
	}
}

func TestGetBearerToken(t *testing.T) {
	testHeader := http.Header{}
	testHeader.Set("Authorization", "Bearer test_token")

	token, err := GetBearerToken(testHeader)
	if err != nil {
		t.Errorf("Error getting bearer token: %s", err)
	}
	if token != "test_token" {
		t.Errorf("Expected token to be test_token, but got: %s", token)
	}
}

func TestGetBearerNoHeader(t *testing.T) {
	testHeader := http.Header{}
	testHeader.Set("", "Bearer test_token")

	_, err := GetBearerToken(testHeader)
	if err == nil {
		t.Errorf("Expected error due to missing header, but got nil")
	}
}

func TestGetBearerNoPrefix(t *testing.T) {
	testHeader := http.Header{}
	testHeader.Set("Authorization", "test_token")

	_, err := GetBearerToken(testHeader)
	if err == nil {
		t.Errorf("Expected error due to missing prefix, but got nil")
	}
}
