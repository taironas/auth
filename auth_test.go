package auth

import "testing"

func TestGoogleWithEmptyIDToken(t *testing.T) {

	var user *User
	var err error
	if user, err = Google(""); err == nil {
		t.Fatal("expected an error and got nothing")
	}

	if user != nil {
		t.Fatal(user)
	}
}

func TestGoogleWithWrongIDToken(t *testing.T) {

	var user *User
	var err error
	if user, err = Google("XYZ123"); err == nil {
		t.Fatal("expected an error and got nothing")
	}

	if user != nil {
		t.Fatal(user)
	}
}
