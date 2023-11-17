package drivers

type Driver interface {
	GetName() string
	GetConnectionString() string
	NoRowsError() string
	UniqueViolationError() string
}
