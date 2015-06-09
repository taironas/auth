package auth

import "testing"
import "fmt"

func TestGoogle(t *testing.T) {
	// test with an empty ID token
	fmt.Println("Test with an empty ID token")
	user, err := Google("")
	if err == nil {
		t.Fatal("expected an error and got nothing")
	}

	if user != nil {
		t.Fatal(user)
	}

	// test with a wrong ID token
	fmt.Println("Test with a wrong ID token")
	user, err = Google("XYZ123")
	if err == nil {
		t.Fatal("expected an error and got nothing")
	}

	if user != nil {
		t.Fatal(user)
	}
}
