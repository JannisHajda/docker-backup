package drivers

type Driver interface {
	GetName() string
	GetConnectionString() string
}
