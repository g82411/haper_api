package utils

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"golang.org/x/oauth2"
	"hyper_api/internal/config"
	"strings"
)

type Claims struct {
	Email        string `json:"email"`
	Sub          string `json:"sub"`
	CogUsername  string `json:"cognito:username"`
	Name         string `json:"name"`
	IsDoneSurvey string `json:"custom:isDoneSurvey"`
}

func ExtractTokenFromCode(code string) (*oauth2.Token, error) {
	env := config.GetConfig()

	oauthConfig := &oauth2.Config{
		ClientID:    env.CognitoClientId,
		RedirectURL: env.RedirectURL,
		Endpoint: oauth2.Endpoint{
			AuthURL:  env.CognitoDomain + "/oauth2/authorize",
			TokenURL: env.CognitoDomain + "/oauth2/token",
		},
		Scopes: []string{"email", "profile", "openid"},
	}
	token, err := oauthConfig.Exchange(oauth2.NoContext, code)
	return token, err
}

func ExtractUserInfoFromToken(token string) (Claims, error) {
	parts := strings.Split(token, ".")
	var claims Claims

	if len(parts) != 3 {
		return claims, fmt.Errorf("invalid token")
	}

	// Base64URL 解码 Payload
	payload, err := base64.RawURLEncoding.DecodeString(parts[1])
	if err != nil {
		return claims, err
	}

	// 解析 JSON
	if err := json.Unmarshal(payload, &claims); err != nil {
		return claims, err
	}
	return claims, nil
}
