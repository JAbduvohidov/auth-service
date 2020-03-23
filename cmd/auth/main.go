package main

import (
	"auth-service/cmd/auth/app"
	"auth-service/pkg/core/services"
	"auth-service/pkg/core/token"
	"auth-service/pkg/core/user"
	"context"
	"flag"
	"fmt"
	"github.com/JAbduvohidov/di/pkg/di"
	"github.com/JAbduvohidov/jwt"
	"github.com/JAbduvohidov/mux"
	"github.com/jackc/pgx/v4/pgxpool"
	"log"
	"net"
	"net/http"
)

var (
	hostF   = flag.String("host", "", "Server host")
	portF   = flag.String("port", "", "Server port")
	secretF = flag.String("secret", "", "Server secret")
	dsnF    = flag.String("dsn", "", "Postgres DSN")
)

var (
	EHOST   = "HOST"
	EPORT   = "PORT"
	ESECRET = "SECRET"
	EDSN    = "DATABASE_URL"
)

type DSN string

func main() {
	flag.Parse()

	host, ok := FlagOrEnv(*hostF, EHOST)
	if !ok {
		log.Panic("can't get host")
	}
	port, ok := FlagOrEnv(*portF, EPORT)
	if !ok {
		log.Panic("can't get port")
	}
	secret, ok := FlagOrEnv(*secretF, ESECRET)
	if !ok {
		log.Panic("can't get secret")
	}
	dsn, ok := FlagOrEnv(*dsnF, EDSN)
	if !ok {
		log.Panic("can't get dsn")
	}

	addr := net.JoinHostPort(host, port)

	start(addr, dsn, jwt.Secret(secret))
}

func start(addr string, dsn string, secret jwt.Secret) {

	err := services.InitDB(dsn)
	if err != nil {
		panic(err)
	}

	container := di.NewContainer()

	err = container.Provide(
		app.NewServer,
		mux.NewExactMux,
		func() jwt.Secret { return secret },
		func() DSN { return DSN(dsn) },
		func(dsn DSN) *pgxpool.Pool {
			pool, err := pgxpool.Connect(context.Background(), string(dsn))
			if err != nil {
				panic(fmt.Errorf("can't create pool: %w", err))
			}
			return pool
		},
		token.NewService,
		user.NewService,
	)

	if err != nil {
		panic(fmt.Errorf("unable to provide di: %w", err))
	}

	container.Start()
	var appServer *app.Server
	container.Component(&appServer)

	panic(http.ListenAndServe(addr, appServer))
}
