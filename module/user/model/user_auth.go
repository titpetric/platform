package model

func (u *UserAuth) Valid() bool {
	if u == nil || u.Email == "" || u.Password == "" {
		return false
	}
	return true
}
