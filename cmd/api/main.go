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
)

const version = "1.0"

type dbConfig struct {
	host     string
	dbname   string
	port     int64
	user     string
	password string
	sslmode  string
}
type config struct {
	servPort int
	envMode  string
	db       dbConfig
}
type application struct {
	config config
	logger *log.Logger
}

func main() {
	if err := initEnv(); err != nil {
		log.Fatalf("fatal error to read env file. error: %s", err)
	}

	app, err := initConfigApp()
	if err != nil {
		panic(err)
	}

	connect, err := initDB(&app.config)

	if err != nil {
		app.logger.Fatalf("fatal error to open connection with db. error: %s", err)
	}

	defer connect.Close()

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

func initConfigApp() (*application, error) {
	logger := log.New(os.Stdout, "", log.Ldate|log.Ltime)
	port, err := strconv.ParseInt(os.Getenv("DB_PORT"), 10, 32)
	if err != nil {
		return nil, err
	}
	dbConfig := dbConfig{
		user:     os.Getenv("DB_USER"),
		password: os.Getenv("DB_PASSWORD"),
		host:     os.Getenv("DB_HOST"),
		port:     port,
		dbname:   os.Getenv("DB_DATABASE"),
		sslmode:  os.Getenv("DB_SSL_MODE"),
	}
	app := &application{
		config: config{
			db: dbConfig,
		},
		logger: logger,
	}
	flag.IntVar(&app.config.servPort, "port", 8000, "API server port")
	flag.StringVar(&app.config.envMode, "env", "dev", "Environment (dev|staging|production)")
	flag.Parse()
	return app, nil
}

func initDB(cfg *config) (*sql.DB, error) {
	dsn := fmt.Sprintf("user=%s password=%s host=%s port=%d dbname=%s sslmode=%s",
		cfg.db.user, cfg.db.password, cfg.db.host, cfg.db.port, cfg.db.dbname, cfg.db.sslmode)
	conn, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, err
	}
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err = conn.PingContext(ctx); err != nil {
		return nil, err
	}

	fmt.Println("db connection success! 🐭")
	initEnv()
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
