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
	mPassF  = flag.String("mpass", "", "Moderators password")
)

var (
	EHOST   = "HOST"
	EPORT   = "PORT"
	ESECRET = "SECRET"
	EDSN    = "DATABASE_URL"
	EMPASS  = "MODERATOR_PASS"
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
	mPass, ok := FlagOrEnv(*mPassF, EMPASS)
	if !ok {
		log.Panic("can't get moderators pass")
	}

	addr := net.JoinHostPort(host, port)

	start(addr, dsn, jwt.Secret(secret), mPass)
}

func start(addr string, dsn string, secret jwt.Secret, mPass string) {

	err := services.InitDB(dsn, mPass)
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
