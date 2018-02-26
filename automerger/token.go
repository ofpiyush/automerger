package automerger

import (
	"net/http"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
)

type TokenResponse struct {
	Token     string    `json:"token,omitempty"`
	ExpiresAt time.Time `json:"expires_at,omitempty"`
}

type Token struct {
	Config *Config
}

func (t *Token) GetorDie(installationID int) (string, []error) {
	key, err := jwt.ParseRSAPrivateKeyFromPEM(t.Config.Key)
	if err != nil {
		return "", []error{err}
	}
	token, err := jwt.NewWithClaims(jwt.SigningMethodRS256, jwt.StandardClaims{
		IssuedAt:  time.Now().Unix(),
		ExpiresAt: time.Now().Add(10 * time.Minute).Unix(),
		Issuer:    t.Config.IntegrationID,
	}).SignedString(key)
	if err != nil {
		return "", []error{err}
	}

	var result = &TokenResponse{}

	return result.Token, GithubRequest("POST", URL(t.Config.ApiURL, "installations", installationID, "access_tokens"), token, http.StatusCreated, nil, &result)
}
