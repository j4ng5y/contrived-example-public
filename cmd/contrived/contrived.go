package main

import (
	"errors"
	"fmt"
	"log/slog"
	"os"
	"time"

	"github.com/j4ng5y/contrived-example-public/pkg/app"
	"github.com/peterbourgon/ff/v4"
	"github.com/peterbourgon/ff/v4/ffhelp"
)

var (
	logger = slog.New(slog.NewTextHandler(os.Stdout, nil))
	fs     = ff.NewFlagSet("contrived")
	addr   = fs.StringLong("server.addr", "0.0.0.0", "The address to bind the server to")
	port   = fs.Uint64Long("server.port", 8080, "the TCP port to bind the server to")
	mock   = fs.BoolLong("mock", "Whether to use a mock database rather than a real database")
	_      = fs.StringLong("config", "", "The config file to use (optional)")
)

func init() {
	if err := ff.Parse(fs, os.Args[1:],
		ff.WithEnvVarPrefix("CONTRIVED"),
		ff.WithConfigFileFlag("config"),
		ff.WithConfigFileParser(ff.PlainParser),
	); err != nil {
		if errors.Is(err, ff.ErrHelp) {
			fmt.Printf("%s\n", ffhelp.Flags(fs))
			os.Exit(0)
		}
	}
}

func main() {
	opts := []app.Option{
		app.WithHTTPAddr(*addr),
		app.WithHTTPPort(*port),
		app.WithHTTPReadTimeout(600 * time.Second),
		app.WithHTTPReadHeaderTimeout(600 * time.Second),
		app.WithHTTPWriteTimeout(600 * time.Second),
		app.WithHTTPIdleTimeout(600 * time.Second),
		app.WithLogger(logger),
	}

	if *mock {
		opts = append(opts, app.WithMemoryRepo())
	} else {
		opts = append(opts, app.WithSqliteRepo())
	}

	srv := app.New(opts...)

	if err := srv.Run(); err != nil {
		logger.Error("error running server", "error", err)
	}
}
