package internal

import "github.com/jmoiron/sqlx"

type DatabaseOption struct {
	MaxOpenConns int
	MaxIdleConns int
}

func (o *DatabaseOption) Apply(client *sqlx.DB) {
	if o == nil {
		return
	}
	client.SetMaxOpenConns(o.MaxOpenConns)
	client.SetMaxIdleConns(o.MaxIdleConns)
}

var databaseOptions = map[string]DatabaseOption{
	"sqlite": {
		MaxOpenConns: 1,
		MaxIdleConns: 1,
	},
	"mysql": {
		MaxOpenConns: 10,
		MaxIdleConns: 10,
	},
}
