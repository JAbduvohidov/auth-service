package app

import (
	"auth-service/pkg/core/token"
	"auth-service/pkg/core/user"
	"errors"
	"github.com/JAbduvohidov/jwt"
	"github.com/JAbduvohidov/mux"
	jwt2 "github.com/JAbduvohidov/mux/middleware/jwt"
	"github.com/JAbduvohidov/rest"
	"github.com/jackc/pgx/v4/pgxpool"
	"log"
	"net/http"
	"strconv"
)

type Server struct {
	router   *mux.ExactMux
	pool     *pgxpool.Pool
	secret   jwt.Secret
	tokenSvc *token.Service
	userSvc  *user.Service
}

func NewServer(router *mux.ExactMux, pool *pgxpool.Pool, secret jwt.Secret, tokenSvc *token.Service, userSvc *user.Service) *Server {
	return &Server{router: router, pool: pool, secret: secret, tokenSvc: tokenSvc, userSvc: userSvc}
}

func (s *Server) Start() {
	s.InitRoutes()
}

func (s *Server) Stop() {
	// TODO: make server stop
}

type ErrorDTO struct {
	Errors []string `json:"errors"`
}

func (s *Server) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	s.router.ServeHTTP(writer, request)
}

func (s *Server) handleCreateToken() http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		var body token.RequestDTO
		err := rest.ReadJSONBody(request, &body)
		if err != nil {
			writer.WriteHeader(http.StatusBadRequest)
			err := rest.WriteJSONBody(writer, &ErrorDTO{
				[]string{"err.json_invalid"},
			})
			log.Print(err)
			return
		}

		response, err := s.tokenSvc.Generate(request.Context(), &body)
		if err != nil {
			switch {
			case errors.Is(err, token.ErrInvalidLogin):
				writer.WriteHeader(http.StatusBadRequest)
				_ = rest.WriteJSONBody(writer, &ErrorDTO{
					[]string{"err.login_mismatch"},
				})
				log.Print(err)
			case errors.Is(err, token.ErrInvalidPassword):
				writer.WriteHeader(http.StatusBadRequest)
				_ = rest.WriteJSONBody(writer, &ErrorDTO{
					[]string{"err.password_mismatch"},
				})
				log.Print(err)
			default:
				writer.WriteHeader(http.StatusBadRequest)
				_ = rest.WriteJSONBody(writer, &ErrorDTO{
					[]string{"err.unknown"},
				})
				log.Print(err)
			}
			return
		}

		err = rest.WriteJSONBody(writer, &response)
		if err != nil {
			log.Print(err)
		}
	}
}

func (s *Server) handleGetProfile() http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		idFromContext := request.Context().Value("id")
		payload := request.Context().Value(jwt2.ContextKey("jwt")).(*token.Payload)
		if idFromContext == payload.Id {
			writer.WriteHeader(http.StatusBadRequest)
			_ = rest.WriteJSONBody(writer, &ErrorDTO{
				[]string{"err.bad_request"},
			})
			return
		}

		response, err := s.userSvc.GetProfile(request.Context())
		if err != nil {
			writer.WriteHeader(http.StatusBadRequest)
			_ = rest.WriteJSONBody(writer, &ErrorDTO{
				[]string{"err.bad_request"},
			})
			log.Print(err)
			return
		}

		err = rest.WriteJSONBody(writer, &response)
		if err != nil {
			log.Print(err)
		}

	}
}

func (s *Server) handlePostProfile() http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		idFromContext := request.Context().Value("id")
		payload := request.Context().Value(jwt2.ContextKey("jwt")).(*token.Payload)
		if idFromContext == payload.Id {
			writer.WriteHeader(http.StatusBadRequest)
			_ = rest.WriteJSONBody(writer, &ErrorDTO{
				[]string{"err.bad_request"},
			})
			return
		}

		var dto user.RequestDTO
		err := rest.ReadJSONBody(request, &dto)
		if err != nil {
			writer.WriteHeader(http.StatusBadRequest)
			_ = rest.WriteJSONBody(writer, &ErrorDTO{
				[]string{"err.bad_request"},
			})
			log.Print(err)
			return
		}

		err = s.userSvc.EditProfile(request.Context(), &dto)
		if err != nil {
			writer.WriteHeader(http.StatusBadRequest)
			_ = rest.WriteJSONBody(writer, &ErrorDTO{
				[]string{"err.bad_request"},
			})
			log.Print(err)
			return
		}

		writer.WriteHeader(http.StatusOK)
	}
}

func (s *Server) handleNewUser() http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		var dto user.RequestDTO
		err := rest.ReadJSONBody(request, &dto)
		if err != nil {
			writer.WriteHeader(http.StatusBadRequest)
			_ = rest.WriteJSONBody(writer, &ErrorDTO{
				[]string{"err.bad_request"},
			})
			log.Print(err)
			return
		}

		userProfile := user.RequestDTO{
			Id:       dto.Id,
			Name:     dto.Name,
			Surname:  dto.Surname,
			Login:    dto.Login,
			Password: dto.Password,
			Avatar:   dto.Avatar,
		}

		err = s.userSvc.AddUser(request.Context(), userProfile)
		if err != nil {
			writer.WriteHeader(http.StatusBadRequest)
			err := rest.WriteJSONBody(writer, &ErrorDTO{
				[]string{"err.invalid_values"},
			})
			log.Print(err)
			return
		}
	}
}

func (s *Server) handleGetUsers() http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		users, err := s.userSvc.GetUsers(request.Context())
		if err != nil {
			writer.WriteHeader(http.StatusBadRequest)
			err := rest.WriteJSONBody(writer, &ErrorDTO{
				[]string{"err.internal_server_error"},
			})
			log.Print(err)
			return
		}

		err = rest.WriteJSONBody(writer, &users)
		if err != nil {
			log.Print(err)
		}
	}
}

func (s *Server) handleDeleteUser() http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		idFromContext := request.Context().Value("id")
		payload := request.Context().Value(jwt2.ContextKey("jwt")).(*token.Payload)
		if idFromContext == payload.Id {
			writer.WriteHeader(http.StatusBadRequest)
			_ = rest.WriteJSONBody(writer, &ErrorDTO{
				[]string{"err.bad_request"},
			})
			return
		}

		if payload.Role	!= "MODERATOR" {
			writer.WriteHeader(http.StatusBadRequest)
			_ = rest.WriteJSONBody(writer, &ErrorDTO{
				[]string{"err.invalid_privilege"},
			})
			return
		}

		idStr := idFromContext.(string)
		id, _ := strconv.Atoi(idStr)

		err := s.userSvc.DeleteUser(request.Context(), id)
		if err != nil {
			writer.WriteHeader(http.StatusBadRequest)
			err := rest.WriteJSONBody(writer, &ErrorDTO{
				[]string{"err.internal_server_error"},
			})
			log.Print(err)
			return
		}
	}
}

func (s *Server) handleUpgradeUser() http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		idFromContext := request.Context().Value("id")
		roleFromContext := request.Context().Value("role")
		payload := request.Context().Value(jwt2.ContextKey("jwt")).(*token.Payload)
		if idFromContext == payload.Id {
			writer.WriteHeader(http.StatusBadRequest)
			_ = rest.WriteJSONBody(writer, &ErrorDTO{
				[]string{"err.bad_request"},
			})
			return
		}

		if payload.Role	!= "MODERATOR" {
			writer.WriteHeader(http.StatusBadRequest)
			_ = rest.WriteJSONBody(writer, &ErrorDTO{
				[]string{"err.invalid_privilege"},
			})
			return
		}

		idStr := idFromContext.(string)
		id, _ := strconv.Atoi(idStr)

		role := roleFromContext.(string)

		switch role {
		case "ADMIN":
			role = "USER"
		case "USER":
			role = "ADMIN"
		default:
			role = "USER"
		}

		err := s.userSvc.UpgradeUser(request.Context(), role, id)
		if err != nil {
			writer.WriteHeader(http.StatusBadRequest)
			err := rest.WriteJSONBody(writer, &ErrorDTO{
				[]string{"err.internal_server_error"},
			})
			log.Print(err)
			return
		}
	}
}
