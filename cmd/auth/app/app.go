package app

import (
	"auth-service/pkg/core/token"
	"auth-service/pkg/core/user"
	"errors"
	"github.com/JAbduvohidov/jwt"
	"github.com/JAbduvohidov/mux"
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

func (s *Server) handleProfile() http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		response, err := s.userSvc.Profile(request.Context())
		if err != nil {
			writer.WriteHeader(http.StatusBadRequest)
			_ = rest.WriteJSONBody(writer, &ErrorDTO{
				[]string{"err.bad_request"},
			})
			log.Print(err)
			return
		}

		if request.Context().Value("id") != strconv.Itoa(int(response.Id)) {
			writer.WriteHeader(http.StatusBadRequest)
			err = rest.WriteJSONBody(writer, &ErrorDTO{
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
				[]string{"err.bad_request"},
			})
			log.Print(err)
			return
		}
	}
}
