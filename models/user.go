package models

import (
	"auth"
	"errors"
)

type User struct {
	Id            string `form:"-"`
	Username      string `form:"username"`
	Password      string `form:"password"`
	authenticated bool   `form:"-"`
}

// Return whether this user is logged in or not
func (m *User) IsAuthenticated() bool {
	return m.authenticated
}

// Set any flags or extra data that should be available
func (m *User) Login() {
	m.authenticated = true
}

// Clear any sensitive data out of the user
func (m *User) Logout() {
	m.authenticated = false
}

// Return the unique identifier of this user object
func (m *User) UniqueId() interface{} {
	return m.Id
}

// Populate this user object with values
func (m *User) GetById(id interface{}) error {

	if id.(string) == auth.MasterUserId {

		m.Id = auth.MasterUserId
		m.Username = auth.MasterUserUsername
		m.Password = auth.MasterUserPassword

		return nil
	} else {
		return errors.New("Bad Credentials")
	}

}
