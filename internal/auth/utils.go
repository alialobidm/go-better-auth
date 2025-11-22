package auth

import (
	"github.com/GoBetterAuth/go-better-auth/pkg/domain"
)

// callHook safely calls a hook function if it's not nil
func (s *Service) callHook(hook func(domain.User) error, user *domain.User) {
	if hook != nil && user != nil {
		go hook(*user)
	}
}
