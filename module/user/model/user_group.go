package model

func NewUserGroup() *UserGroup {
	return &UserGroup{}
}

func (u *UserGroup) String() string {
	return u.Title
}
