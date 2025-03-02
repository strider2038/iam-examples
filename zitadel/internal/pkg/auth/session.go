package auth

import (
	"errors"
	"net/http"
)

var ErrUnauthorized = errors.New("unauthorized")

type Session struct {
	UserID    string
	CompanyID string
}

func (s *Session) IsValid() bool {
	return s.UserID != "" && s.CompanyID != ""
}

func ParseSessionFromRequest(request *http.Request) (*Session, error) {
	session := &Session{
		UserID:    request.Header.Get("X-User-Id"),
		CompanyID: request.Header.Get("X-Company-Id"),
	}
	if !session.IsValid() {
		return nil, ErrUnauthorized
	}

	return session, nil
}
