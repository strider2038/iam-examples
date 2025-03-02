package domain

import "context"

type User struct {
	ID        string
	Email     string
	FirstName string
	LastName  string
}

type UserRepository interface {
	FindByIDs(ctx context.Context, ids []string) (map[string]*User, error)
}
