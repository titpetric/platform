package internal

type DatabaseOption struct {
	MaxOpenConns int
	MaxIdleConns int
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
