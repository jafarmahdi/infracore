package identity

import "testing"

func TestUserFullName(t *testing.T) {
	tests := []struct {
		name string
		user User
		want string
	}{
		{name: "first and last name", user: User{FirstName: "Omar", LastName: "Hassan", Username: "omar"}, want: "Omar Hassan"},
		{name: "first name only", user: User{FirstName: "Omar", Username: "omar"}, want: "Omar"},
		{name: "last name only", user: User{LastName: "Hassan", Username: "omar"}, want: "Hassan"},
		{name: "username fallback", user: User{Username: "omar"}, want: "omar"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.user.FullName(); got != tt.want {
				t.Fatalf("FullName() = %q, want %q", got, tt.want)
			}
		})
	}
}
