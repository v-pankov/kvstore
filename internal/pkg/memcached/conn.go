package memcached

type Connection interface {
	Client
	Close() error
}

type ConnectionFactory interface {
	CreateConnection() (Connection, error)
}
