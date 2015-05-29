package auth

// Google returns Google user from the given ID token
func Google(IDToken string) (*User, error) {

	return nil, nil
}

// User represents the authenticated user
type User struct {
	ID    string
	Name  string
	Email string
}
