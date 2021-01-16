package app

import (
	"auth-service/pkg/core/token"
	"context"
	"github.com/JAbduvohidov/mux/middleware/authenticated"
	"github.com/JAbduvohidov/mux/middleware/jwt"
	"github.com/JAbduvohidov/mux/middleware/logger"
	"reflect"
)

func (s *Server) InitRoutes() {
	s.router.POST("/api/tokens",
		s.handleCreateToken(),
		logger.Logger("CREATE_TOKEN"),
	)
	s.router.GET("/api/users/{id}",
		s.handleGetProfile(),
		authenticated.Authenticated(func(ctx context.Context) bool { return !jwt.IsContextNonEmpty(ctx) }, false, ""),
		jwt.JWT(jwt.SourceAuthorization, reflect.TypeOf((*token.Payload)(nil)).Elem(), s.secret),
		logger.Logger("PROFILE"),
	)

	s.router.POST("/api/users/{id}",
		s.handlePostProfile(),
		authenticated.Authenticated(func(ctx context.Context) bool { return !jwt.IsContextNonEmpty(ctx) }, false, ""),
		jwt.JWT(jwt.SourceAuthorization, reflect.TypeOf((*token.Payload)(nil)).Elem(), s.secret),
		logger.Logger("PROFILE"),
	)

	s.router.POST("/api/users",
		s.handleNewUser(),
		logger.Logger("NEW_USER"),
	)

	s.router.GET("/api/users",
		s.handleGetUsers(),
		authenticated.Authenticated(func(ctx context.Context) bool { return !jwt.IsContextNonEmpty(ctx) }, false, ""),
		jwt.JWT(jwt.SourceAuthorization, reflect.TypeOf((*token.Payload)(nil)).Elem(), s.secret),
		logger.Logger("GET_USERS"),
	)

	s.router.DELETE("/api/users/{id}",
		s.handleDeleteUser(),
		authenticated.Authenticated(func(ctx context.Context) bool { return !jwt.IsContextNonEmpty(ctx) }, false, ""),
		jwt.JWT(jwt.SourceAuthorization, reflect.TypeOf((*token.Payload)(nil)).Elem(), s.secret),
		logger.Logger("DELETE_USER"),
	)

	s.router.PUT("/api/users/{id}/{role}",
		s.handleUpgradeUser(),
		authenticated.Authenticated(func(ctx context.Context) bool { return !jwt.IsContextNonEmpty(ctx) }, false, ""),
		jwt.JWT(jwt.SourceAuthorization, reflect.TypeOf((*token.Payload)(nil)).Elem(), s.secret),
		logger.Logger("UPGRADE_USER"),
	)
}
