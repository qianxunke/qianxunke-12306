package http

import (
	"net/http"
	"strings"

	"github.com/micro/cli"
	"github.com/micro/micro/plugin"
	"github.com/rs/cors"
)

type allowedCors struct {
	allowedHeaders []string
	allowedOrigins []string
	allowedMethods []string
}

func (ac *allowedCors) Flags() []cli.Flag {
	return []cli.Flag{
		cli.StringFlag{
			Name:   "cors-allowed-headers",
			Usage:  "Comma-seperated list of allowed headers",
			EnvVar: "CORS_ALLOWED_HEADERS",
		},
		cli.StringFlag{
			Name:   "cors-allowed-origins",
			Usage:  "Comma-seperated list of allowed origins",
			EnvVar: "CORS_ALLOWED_ORIGINS",
		},
		cli.StringFlag{
			Name:   "cors-allowed-methods",
			Usage:  "Comma-seperated list of allowed methods",
			EnvVar: "CORS_ALLOWED_METHODS",
		},
	}
}

func (ac *allowedCors) Commands() []cli.Command {
	return nil
}

func (ac *allowedCors) Handler() plugin.Handler {
	return func(ha http.Handler) http.Handler {
		hf := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ha.ServeHTTP(w, r)
		})

		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			cors.New(cors.Options{
				AllowedOrigins:   []string{"*", "http://127.0.0.1"},
				AllowedMethods:   []string{"POST", "GET", "OPTIONS", "PUT", "DELETE"},
				AllowedHeaders:   []string{"Content-Type", "AccessToken", "Origin", "Accept", "X-Requested-With", "Authorization", "Token", "Content-Length"},
				AllowCredentials: true,
			}).ServeHTTP(w, r, hf)
		})
	}
}

func (ac *allowedCors) Init(ctx *cli.Context) error {
	ac.allowedHeaders = ac.parseAllowed(ctx, "cors-allowed-headers")
	ac.allowedMethods = ac.parseAllowed(ctx, "cors-allowed-methods")
	ac.allowedOrigins = ac.parseAllowed(ctx, "cors-allowed-origins")

	return nil
}

func (ac *allowedCors) parseAllowed(ctx *cli.Context, flagName string) []string {
	fv := ctx.String(flagName)

	// no op
	if len(fv) == 0 {
		return nil
	}

	return strings.Split(fv, ",")
}

func (ac *allowedCors) String() string {
	return "cors-allowed-(headers|origins|methods)"
}

// NewPlugin Creates the CORS Plugin
func NewPlugin() plugin.Plugin {
	return &allowedCors{
		allowedHeaders: []string{},
		allowedOrigins: []string{},
		allowedMethods: []string{},
	}
}
