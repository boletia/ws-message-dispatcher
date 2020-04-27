package service

// UserStorage get users from storage
type userStorage interface {
	GetUserConnections(eventSubdomain string, audienceType string, connections *[]string) error
}

type service struct {
	dbUser userStorage
}

// New creates new service
func New(dbUser userStorage) service {
	return service{
		dbUser: dbUser,
	}
}
