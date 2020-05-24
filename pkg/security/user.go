package security

// Represents a user
type User struct {
	// The name of the user
	Name string
	// The roles the user has access too
	Roles []string
}

// Represents a user that's not logged in
var NotLoggedInUser = &User{"", []string{Read}}

// Check to see if user is logged in
func (u *User) IsLoggedIn() bool {
	return u.Name != ""
}

// Check to see if user has role
func (u *User) HasRole(role string) bool {
	for _, r := range u.Roles {
		if r == role {
			return true
		}
	}
	return false
}

// Check to see if this user has the supplied roles
func (u *User) HasRoles(roles []string) bool {
	for _, r := range roles {
		if !u.HasRole(r) {
			return false
		}
	}
	return true
}
