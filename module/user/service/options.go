package service

import (
	"io/fs"

	"github.com/titpetric/platform/module/user/storage"
)

// Options are the dependencies to pass to the constructor.
type Options struct {
	ThemeFS  fs.FS
	ModuleFS fs.FS

	UserStorage    *storage.UserStorage
	SessionStorage *storage.SessionStorage
}
