package service

// UserStorage get users from storage
type connectionGetter interface {
	GetUserConnections(eventSubdomain string, audienceType string, connections *[]string) error
}

type messageSender interface {
	SendMessage(connections []string, msg interface{})
}

type service struct {
	dbUser connectionGetter
	sender messageSender
}

// New creates new service
func New(dbUser connectionGetter, sender messageSender) service {
	return service{
		dbUser: dbUser,
		sender: sender,
	}
}
