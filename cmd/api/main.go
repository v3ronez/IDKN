package main

import (
	"context"
	"database/sql"
	"expvar"
	"flag"
	"fmt"
	"log"
	"os"
	"runtime"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"github.com/v3ronez/IDKN/internal/data"
	"github.com/v3ronez/IDKN/internal/jsonlog"
	"github.com/v3ronez/IDKN/internal/mailer"
)

const version = "1.0"

type dbConfig struct {
	host         string
	dbname       string
	port         int64
	user         string
	password     string
	sslmode      string
	maxOpenConns int
	maxIdleConns int
	maxIndleTime string
}

type config struct {
	servPort int
	envMode  string
	db       dbConfig
	smtp     struct {
		host     string
		port     int
		username string
		password string
	}
}
type application struct {
	config  config
	logger  *jsonlog.Logger
	models  data.Models
	limiter struct {
		rps     float64
		burst   int
		enabled bool
	}
	mailer mailer.Mailer
	wg     sync.WaitGroup
}

func main() {
	if err := initEnv(); err != nil {
		log.Fatalf("fatal error to read env file. error: %s", err)
	}

	config := &config{}
	flag.IntVar(&config.servPort, "port", 8000, "API server port")
	flag.StringVar(&config.envMode, "env", "dev", "Environment (dev|staging|production)")
	flag.IntVar(&config.db.maxOpenConns, "db-max-open-conns", 25, "set default value to db max open conns")
	flag.IntVar(&config.db.maxIdleConns, "db-max-idle-conns", 25, "set default value to db max idle conns")
	flag.StringVar(&config.db.maxIndleTime, "db-max-idle-time", "15m", "set default value db to idle time conn")
	flag.Parse()

	app, err := initConfigApp(config)
	if err != nil {
		panic(err)
	}

	connect, err := initDB(&app.config)

	if err != nil {
		app.logger.PrintFatal(err, nil)
	}

	defer connect.Close()
	app.models = data.NewModels(connect)

	//metrics
	expvar.NewString("version").Set(version)
	expvar.Publish("goroutines", expvar.Func(func() any {
		return runtime.NumGoroutine()
	}))
	expvar.Publish("CPU", expvar.Func(func() any {
		return runtime.NumCPU()
	}))
	expvar.Publish("database", expvar.Func(func() any {
		return app.models.Users.DB.Stats()
	}))
	expvar.Publish("timestamp", expvar.Func(func() any {
		return time.Now().Unix()
	}))

	err = app.server()
	if err != nil {
		app.logger.PrintFatal(err, nil)
	}
}

func initConfigApp(cfg *config) (*application, error) {
	port, err := strconv.ParseInt(os.Getenv("DB_PORT"), 10, 32)
	if err != nil {
		return nil, err
	}
	portSmtp, err := strconv.Atoi(os.Getenv("EMAIL_PORT"))
	if err != nil {
		return nil, err
	}
	cfg.db.user = os.Getenv("DB_USER")
	cfg.db.password = os.Getenv("DB_PASSWORD")
	cfg.db.host = os.Getenv("DB_HOST")
	cfg.db.port = port
	cfg.db.dbname = os.Getenv("DB_DATABASE")
	cfg.db.sslmode = os.Getenv("DB_SSL_MODE")
	cfg.smtp.host = os.Getenv("EMAIL_HOST")
	cfg.smtp.port = portSmtp
	cfg.smtp.username = os.Getenv("EMAIL_USERNAME")
	cfg.smtp.password = os.Getenv("EMAIL_PASSWORD")

	app := &application{
		config: *cfg,
		mailer: mailer.New(
			cfg.smtp.host,
			cfg.smtp.port,
			cfg.smtp.username,
			cfg.smtp.password,
			"Greenlight <no-reply@greenlight.alexedwards.net>"),
		logger: jsonlog.New(os.Stdout, jsonlog.LevelInfo),
		limiter: struct {
			rps     float64
			burst   int
			enabled bool
		}{2, 4, true},
	}

	return app, nil
}

func initDB(cfg *config) (*sql.DB, error) {
	dsn := fmt.Sprintf("user=%s password=%s host=%s port=%d dbname=%s sslmode=%s",
		cfg.db.user,
		cfg.db.password,
		cfg.db.host,
		cfg.db.port,
		cfg.db.dbname,
		cfg.db.sslmode)

	conn, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, err
	}
	conn.SetConnMaxIdleTime(time.Duration(cfg.db.maxIdleConns))
	conn.SetMaxOpenConns(cfg.db.maxOpenConns)
	durantion, err := time.ParseDuration(cfg.db.maxIndleTime)
	if err != nil {
		return nil, err
	}
	conn.SetConnMaxLifetime(durantion)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err = conn.PingContext(ctx); err != nil {
		return nil, err
	}

	fmt.Println("db connection success! 🐭")
	return conn, nil
}

func initEnv() error {
	path, _ := os.Getwd()
	fullPath := strings.Join([]string{path, "/.env"}, "")
	if err := godotenv.Load(fullPath); err != nil {
		return err
	}
	return nil
}
