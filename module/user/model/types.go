package model

func (u *User) String() string {
	if u.DeletedAt != nil {
		return "Deleted user"
	}
	return u.FirstName + " " + u.LastName
}

func (u *User) IsActive() bool {
	return u.DeletedAt == nil
}

func (u *UserGroup) String() string {
	return u.Title
}
