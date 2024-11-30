package middleware

import (
	"context"

	"github.com/obanoff/gator/internal/commands"
	"github.com/obanoff/gator/internal/config"
)

func MiddlewareLoggedIn(handler commands.CommandHandler) commands.CommandHandler {
	return commands.CommandHandler(func(s *config.State, args []string) error {
		_, err := s.DB.Queries.GetUserByName(context.Background(), s.Config.Username)
		if err != nil {
			return err
		}

		return handler(s, args)
	})
}
