package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"time"

	"bitbucket.org/atlassian/pushan"
	"github.com/caarlos0/env"
	"github.com/justinas/alice"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/hlog"
)

var log = zerolog.New(os.Stdout).With().
	Timestamp().
	Logger()

type appCfg struct {
	// server config
	HTTPPort string `env:"HTTP_PORT" envDefault:"8080"`
}

func main() {
	app := appCfg{}
	err := env.Parse(&app)
	if err != nil {
		log.Fatal().Err(err).Msg("service is misconfigured")
	}

	// Create an http mux and a server for that mux.
	handler := http.NewServeMux()
	addr := fmt.Sprintf(":%s", app.HTTPPort)
	srv := &http.Server{
		// It's important to set http server timeouts for the publicly available service api.
		// 5 seconds between when connection is accepted to when the body is fully reaad.
		ReadTimeout: 5 * time.Second,
		// 10 seconds from end of request headers read to end of response write.
		WriteTimeout: 10 * time.Second,
		// 120 seconds for an idle KeeP-Alive connection.
		IdleTimeout: 120 * time.Second,
		Addr:        addr,
		Handler:     handler,
	}

	// Create a middleware chain setup to log http access and inject
	// a logger into the request context.
	chain := alice.New(
		hlog.NewHandler(log),
		hlog.AccessHandler(func(r *http.Request, status, size int, duration time.Duration) {
			hlog.FromRequest(r).Info().
				Str("method", r.Method).
				Str("url", r.URL.String()).
				Int("status", status).
				Int("size", size).
				Dur("duration", duration).
				Msg("")
		}),
		hlog.RemoteAddrHandler("ip"),
		hlog.UserAgentHandler("user_agent"),
		hlog.RefererHandler("referer"),
		hlog.RequestIDHandler("req_id", "Request-Id"),
	)

	slashJuneBug := chain.ThenFunc(junebug.Slash)

	handler.Handle("/slash/junebug", slashJuneBug)
	handler.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, "buzz buzz")
	})

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt)
	go func() {
		log.Info().Msgf("listening on %s...", addr)
		err := srv.ListenAndServe()
		log.Fatal().Err(err).Msg("shutting server down")
	}()
	<-stop
	log.Info().Msg("shutting server down")
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()
	err = srv.Shutdown(ctx)
	if err != nil {
		log.Fatal().Err(err).Msg("unable to shutdown cleanly")
	}

}
