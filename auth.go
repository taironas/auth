package auth

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"strings"
)

// Google returns Google user from the given ID token
func Google(IDToken string) (*User, error) {
	// validate the provided ID token
	ti, err := validateIDToken(IDToken)

	if err != nil {
		return nil, err
	}

	// Get the user name
	name, err := extractNameFromToken(IDToken)

	if err != nil {
		return nil, err
	}

	return &User{ID: ti.UserID, Name: name, Email: ti.Email}, err
}

// User represents the authenticated user
type User struct {
	ID    string
	Name  string
	Email string
}

type tokenInfo struct {
	IssuedTo      string `json:"issued_to"`
	Audience      string `json:"audience"`
	UserID        string `json:"user_id"`
	ExpiresIn     int64  `json:"expires_in"`
	Email         string `json:"email"`
	EmailVerified bool   `json:"email_verified"`
}

func validateIDToken(token string) (*tokenInfo, error) {
	// Verification of the integrity of the ID token
	res, err := http.Get("https://www.googleapis.com/oauth2/v2/tokeninfo?id_token=" + token)

	if err != nil {
		return nil, err
	}

	if res == nil || res.StatusCode != http.StatusOK {
		// Bad request
		return nil, errors.New("Provided ID token is not valid.")
	}

	// read the JSON token info
	body, err := ioutil.ReadAll(res.Body)
	res.Body.Close()

	if err != nil {
		return nil, err
	}

	// decode the JSON token info
	var ti *tokenInfo
	json.Unmarshal(body, &ti)

	return ti, nil
}

func extractNameFromToken(token string) (string, error) {
	// split the token by '.'
	parts := strings.Split(token, ".")
	// token is a JWT (http://jwt.io/)
	if len(parts) != 3 {
		return "", errors.New("Provided token is not valid.")
	}

	// decode the second part (claims) containing the name
	claimBytes, err := decodePart(parts[1])

	if err != nil {
		return "", err
	}

	var claims map[string]interface{}
	json.Unmarshal(claimBytes, &claims)

	name, _ := claims["name"].(string)

	// read the JSON and get the name
	return name, nil
}

func decodePart(p string) ([]byte, error) {
	// need padding to have a multiple of four characters and avoid a base64 error
	if l := len(p) % 4; l > 0 {
		p += strings.Repeat("=", 4-l)
	}

	return base64.URLEncoding.DecodeString(p)
}
