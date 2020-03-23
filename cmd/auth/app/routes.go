package app

import (
	"auth-service/pkg/core/token"
	"github.com/JAbduvohidov/mux/middleware/authenticated"
	"github.com/JAbduvohidov/mux/middleware/jwt"
	"github.com/JAbduvohidov/mux/middleware/logger"
	"reflect"
)

func (s *Server) InitRoutes() {
	s.router.POST("/api/tokens",
		s.handleCreateToken(),
		logger.Logger("TOKEN"),
	)
	s.router.GET("/api/users/{id}",
		s.handleProfile(),
		authenticated.Authenticated(jwt.IsContextNonEmpty, false, ""),
		jwt.JWT(jwt.SourceAuthorization, reflect.TypeOf((*token.Payload)(nil)).Elem(), s.secret),
		logger.Logger("USERS"),
	)

	s.router.POST("/api/users",
		s.handleNewUser(),
		jwt.JWT(jwt.SourceAuthorization, reflect.TypeOf((*token.Payload)(nil)).Elem(), s.secret),
		logger.Logger("USERS"),
	)
}
