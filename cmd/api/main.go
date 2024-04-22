package main

import (
	"context"
	"database/sql"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"github.com/v3ronez/IDKN/internal/data"
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
}
type application struct {
	config config
	logger *log.Logger
	models data.Models
}

func main() {
	if err := initEnv(); err != nil {
		log.Fatalf("fatal error to read env file. error: %s", err)
	}
	app := &application{}
	config := &config{}
	flag.IntVar(&config.servPort, "port", 8000, "API server port")
	flag.StringVar(&config.envMode, "env", "dev", "Environment (dev|staging|production)")
	flag.IntVar(&config.db.maxOpenConns, "db-max-open-conns", 25, "set default value to db max open conns")
	flag.IntVar(&config.db.maxIdleConns, "db-max-idle-conns", 25, "set default value to db max idle conns")
	flag.StringVar(&config.db.maxIndleTime, "db-max-idle-time", "15m", "set default value db to idle time conn")
	flag.Parse()

	if err := initConfigApp(app, config); err != nil {
		panic(err)
	}

	connect, err := initDB(&app.config)

	if err != nil {
		app.logger.Fatalf("fatal error to open connection with db. error: %s", err)
	}
	defer connect.Close()
	app.models = data.NewModels(connect)

	serv := &http.Server{
		Addr:         fmt.Sprintf("localhost:%d", app.config.servPort),
		Handler:      app.routes(),
		IdleTimeout:  time.Minute,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
	}

	app.logger.Printf("running %s server on %s", app.config.envMode, serv.Addr)
	if err := serv.ListenAndServe(); err != nil {
		log.Fatal(err)
	}
}

func initConfigApp(app *application, cfg *config) error {
	logger := log.New(os.Stdout, "", log.Ldate|log.Ltime)
	port, err := strconv.ParseInt(os.Getenv("DB_PORT"), 10, 32)
	if err != nil {
		return err
	}
	cfg.db.user = os.Getenv("DB_USER")
	cfg.db.password = os.Getenv("DB_PASSWORD")
	cfg.db.host = os.Getenv("DB_HOST")
	cfg.db.port = port
	cfg.db.dbname = os.Getenv("DB_DATABASE")
	cfg.db.sslmode = os.Getenv("DB_SSL_MODE")

	app.config = *cfg
	app.logger = logger
	return nil
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
