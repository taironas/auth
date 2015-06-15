package auth

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"strings"
)

// Google returns Google user from the given ID token.
//
func Google(IDToken string) (*User, error) {

	var err error

	// validate the provided ID token
	var ti *tokenInfo
	if ti, err = validateIDToken(IDToken); err != nil {
		return nil, err
	}

	// Get the user name
	var name string
	if name, err = extractNameFromToken(IDToken); err != nil {
		return nil, err
	}

	user := User{
		ID:    ti.UserID,
		Name:  name,
		Email: ti.Email,
	}

	return &user, err
}

// User represents the authenticated user.
//
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

	var err error
	// Verification of the integrity of the ID token
	var resp *http.Response
	if resp, err = http.Get("https://www.googleapis.com/oauth2/v2/tokeninfo?id_token=" + token); err != nil {
		return nil, err
	}

	if resp == nil || resp.StatusCode != http.StatusOK {
		// Bad request
		return nil, errors.New("Provided ID token is not valid")
	}

	// read the JSON token info
	var body []byte

	if body, err = ioutil.ReadAll(resp.Body); err != nil {
		return nil, err
	}

	resp.Body.Close()

	// decode the JSON token info
	var ti *tokenInfo
	json.Unmarshal(body, &ti)

	return ti, err
}

func extractNameFromToken(token string) (string, error) {

	parts := strings.Split(token, ".")
	// token is a JWT (http://jwt.io/)
	if len(parts) != 3 {
		return "", errors.New("Provided token is not valid")
	}

	// decode the second part (claims) containing the name
	var err error
	var claimBytes []byte
	if claimBytes, err = decodePart(parts[1]); err != nil {
		return "", err
	}

	var claims map[string]interface{}
	json.Unmarshal(claimBytes, &claims)

	var name string
	var ok bool
	if name, ok = claims["name"].(string); !ok {
		return "", errors.New("Unable to decode name from token")
	}

	return name, nil
}

func decodePart(p string) ([]byte, error) {

	// need padding to have a multiple of four characters and avoid a base64 error
	if l := len(p) % 4; l > 0 {
		p += strings.Repeat("=", 4-l)
	}

	return base64.URLEncoding.DecodeString(p)
}
