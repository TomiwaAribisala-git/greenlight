package main

import (
	"context"
	"database/sql"
	"expvar"
	"flag"
	"os"
	"runtime"
	"sync"
	"time"

	"github.com/TomiwaAribisala-git/greenlight.git/internal/data"
	"github.com/TomiwaAribisala-git/greenlight.git/internal/jsonlog"
	"github.com/TomiwaAribisala-git/greenlight.git/internal/mailer"
	_ "github.com/lib/pq"
)

const version = "1.0.0"

type config struct {
	port int
	env  string
	db   struct {
		dsn          string
		maxOpenConns int
		maxIdleConns int
		maxIdleTime  string
	}
	limiter struct {
		rps     float64
		burst   int
		enabled bool
	}
	smtp struct {
		host     string
		port     int
		username string
		password string
		sender   string
	}
}

type application struct {
	config config
	//logger *log.Logger
	logger *jsonlog.Logger
	models data.Models
	mailer mailer.Mailer
	wg     sync.WaitGroup
}

func main() {

	var cfg config

	flag.IntVar(&cfg.port, "port", 4000, "API server port")
	flag.StringVar(&cfg.env, "env", "development", "Environment (development|staging|production)")
	flag.StringVar(&cfg.db.dsn, "db-dsn", os.Getenv("GREENLIGHT_DB_DSN"), "PostgreSQL DSN")
	// "postgres://greenlight:pa55word@localhost/greenlight

	flag.IntVar(&cfg.db.maxOpenConns, "db-max-open-conns", 25, "PostgreSQL max open connections")
	flag.IntVar(&cfg.db.maxIdleConns, "db-max-idle-conns", 25, "PostgreSQL max idle connections")
	flag.StringVar(&cfg.db.maxIdleTime, "db-max-idle-time", "15m", "PostgreSQL max connection idle time")

	flag.Float64Var(&cfg.limiter.rps, "limiter-rps", 2, "Rate limiter maximum requests per second")
	flag.IntVar(&cfg.limiter.burst, "limiter-burst", 4, "Rate limiter maximum burst")
	flag.BoolVar(&cfg.limiter.enabled, "limiter-enabled", true, "Enable rate limiter")

	flag.StringVar(&cfg.smtp.host, "smtp-host", "sandbox.smtp.mailtrap.io", "SMTP host")
	flag.IntVar(&cfg.smtp.port, "smtp-port", 2525, "SMTP port")
	flag.StringVar(&cfg.smtp.username, "smtp-username", "37b331b8503719", "SMTP username")
	flag.StringVar(&cfg.smtp.password, "smtp-password", "9ed24e0f21b21e", "SMTP password")
	flag.StringVar(&cfg.smtp.sender, "smtp-sender", "Greenlight <no-reply@greenlight.tomiwa.net>", "SMTP sender")

	flag.Parse()

	// logger := log.New(os.Stdout, "", log.Ldate|log.Ltime)
	logger := jsonlog.New(os.Stdout, jsonlog.LevelInfo)

	db, err := openDB(cfg)
	if err != nil {
		//logger.Fatal(err)
		logger.PrintFatal(err, nil)
	}

	defer db.Close()
	//logger.Printf("database connection pool established")
	logger.PrintInfo("database connection pool established", nil)

	// it’s OK to manipulate this value at runtime from your application handlers
	//  you’ll get a runtime panic when the duplicate variable is registered

	// leverage an authentication process for security:
	// create a metrics:view permission so that only certain trusted users can access the endpoint
	// use HTTP Basic Authentication to restrict access to the endpoint
	// Caddy set up, we’ll restrict access to the GET /debug/vars endpoint

	// Custom Metrics

	expvar.NewString("version").Set(version)

	// Publish the number of active goroutines.
	expvar.Publish("goroutines", expvar.Func(func() any {
		return runtime.NumGoroutine()
	}))
	// Publish the database connection pool statistics.
	expvar.Publish("database", expvar.Func(func() any {
		return db.Stats()
	}))
	// Publish the current Unix timestamp.
	expvar.Publish("timestamp", expvar.Func(func() any {
		return time.Now().Unix()
	}))

	app := &application{
		config: cfg,
		logger: logger,
		models: data.NewModels(db),
		mailer: mailer.New(cfg.smtp.host, cfg.smtp.port, cfg.smtp.username, cfg.smtp.password, cfg.smtp.sender),
	}
	// app.models.Movies.Insert(...)

	/*
		srv := &http.Server{
			Addr:    fmt.Sprintf(":%d", cfg.port),
			Handler: app.routes(),
			// log.Logger instance should not use a prefix or any flags
			// log messages of http.Server writes will be passed to Logger.Write() method, which in turn will output a log entry in JSON format at the ERROR
			ErrorLog:     log.New(logger, "", 0),
			IdleTimeout:  time.Minute,
			ReadTimeout:  10 * time.Second,
			WriteTimeout: 30 * time.Second,
		}

		//logger.Printf("starting %s server on %s", cfg.env, srv.Addr)
		logger.PrintInfo("starting server", map[string]string{
			"addr": srv.Addr,
			"env":  cfg.env,
		})

		err = srv.ListenAndServe() // = instead of := because err variable is already declared above
		// logger.Fatal(err)
		logger.PrintFatal(err, nil)
	*/
	err = app.serve()
	if err != nil {
		logger.PrintFatal(err, nil)
	}
}

func openDB(cfg config) (*sql.DB, error) {

	db, err := sql.Open("postgres", cfg.db.dsn)
	if err != nil {
		return nil, err
	}

	db.SetMaxOpenConns(cfg.db.maxOpenConns)
	db.SetMaxIdleConns(cfg.db.maxIdleConns)

	// Use the time.ParseDuration() function to convert the idle timeout duration string
	// to a time.Duration type.
	duration, err := time.ParseDuration(cfg.db.maxIdleTime)
	if err != nil {
		return nil, err
	}
	db.SetConnMaxIdleTime(duration)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err = db.PingContext(ctx)
	if err != nil {
		return nil, err
	}

	return db, nil
}

// Creating an end-to-end test for the GET /v1/healthcheck endpoint to verify that the
// headers and response body are what you expect.
// Creating a unit-test for the rateLimit() middleware to confirm that it sends a
// 429 Too Many Requests response after a certain number of requests.
// Creating an end-to-end integration test, using a test database instance, which confirms
// that the authenticate() and requirePermission() middleware work together correctly
// to allow or disallow access to specific endpoints.
