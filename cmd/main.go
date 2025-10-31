package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"path/filepath"
	"runtime"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	campbackend "github.com/lavatee/camp_backend"
	"github.com/lavatee/camp_backend/internal/endpoint"
	"github.com/lavatee/camp_backend/internal/repository"
	"github.com/lavatee/camp_backend/internal/service"
	_ "github.com/lib/pq"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

func main() {
	logrus.SetFormatter(&logrus.JSONFormatter{PrettyPrint: true})
	if err := InitConfig(); err != nil {
		logrus.Fatalf("error initializing configs: %s", err.Error())
	}

	logrus.Infof("Database configuration: host=%s, port=%s, user=%s, dbname=%s, sslmode=%s",
		viper.GetString("db.host"),
		viper.GetString("db.port"),
		viper.GetString("db.user"),
		viper.GetString("db.dbname"),
		viper.GetString("db.sslmode"))

	db, err := repository.NewPostgresDB(repository.PostgresConfig{
		Host:     viper.GetString("db.host"),
		Port:     viper.GetString("db.port"),
		User:     viper.GetString("db.user"),
		Password: viper.GetString("db.password"),
		DBName:   viper.GetString("db.dbname"),
		SSLMode:  viper.GetString("db.sslmode"),
	})
	if err != nil {
		logrus.Fatalf("error initializing db: %s", err.Error())
	}
	defer db.Close()
	driver, err := postgres.WithInstance(db.DB, &postgres.Config{})
	if err != nil {
		logrus.Fatalf("Failed to create migrate driver: %s", err.Error())
	}

	_, filename, _, _ := runtime.Caller(0)
	dir := filepath.Dir(filename)
	migrationsPath := "file://" + filepath.Join(dir, "../schema")
	migrations, err := migrate.NewWithDatabaseInstance(migrationsPath, "postgres", driver)
	if err != nil {
		logrus.Fatalf("Failed to create migrate instance: %s", err.Error())
	}
	if err = migrations.Up(); err != nil && err != migrate.ErrNoChange {
		logrus.Fatalf("Migrations error: %s", err.Error())
	}
	repo := repository.NewRepository(db)
	s3, err := service.ConnectS3(viper.GetString("s3.url"), viper.GetString("s3.accessKey"), viper.GetString("s3.secretKey"), viper.GetString("s3.region"))
	if err != nil {
		logrus.Fatalf("error connecting to s3: %s", err.Error())
	}
	services := service.NewService(repo, s3, viper.GetString("s3.bucket"))
	endp := endpoint.NewEndpoint(services)
	server := &campbackend.Server{}
	go func() {
		if err := server.Run(viper.GetString("port"), endp.InitRoutes()); err != nil {
			logrus.Fatalf("error running http server: %s", err.Error())
		}
	}()
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT)
	<-quit
	logrus.Print("Shutting down server...")
	if err := server.Shutdown(context.Background()); err != nil {
		logrus.Errorf("error shutting down server: %s", err.Error())
	}

}

func InitConfig() error {
	viper.AddConfigPath("configs")
	viper.SetConfigName("config")
	return viper.ReadInConfig()
}
