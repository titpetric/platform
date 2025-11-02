package model

func NewUser() *User {
	return &User{}
}

func (u *User) String() string {
	if u.DeletedAt != nil {
		return "Deleted user"
	}
	return u.FirstName + " " + u.LastName
}

func (u *User) IsActive() bool {
	return u.DeletedAt == nil
}
