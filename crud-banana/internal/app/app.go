package app

import (
	"os"
	"github.com/berezovskiydeval/crud-banana/internal/delivery/rest"
	rmq "github.com/berezovskiydeval/crud-banana/internal/delivery/rabbitmq"
	"github.com/berezovskiydeval/crud-banana/internal/repository/pdb"
	"github.com/berezovskiydeval/crud-banana/internal/server"
	"github.com/berezovskiydeval/crud-banana/internal/service"
	"github.com/berezovskiydeval/crud-banana/pkg/database"
	"github.com/berezovskiydeval/crud-banana/pkg/hasher"

	"os/signal"
	"syscall"

	"context"
	"time"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

func Run() {
	//Инициализация логгера, конфига и env
	initLogger()

	if err := initConfig(); err != nil {
		log.Fatalf("error initializing configs: %s", err.Error())
	}
	
	if err := godotenv.Load(); err != nil {
		log.Fatalf("error loading env file: %s", err.Error())
	}
		
	//Инициализация структуры БД
	db, err := database.NewPostgresConnection(database.ConnectionInfo{
		Host:     os.Getenv("DB_HOST"),
		Port:     os.Getenv("DB_PORT"),
		Username: os.Getenv("DB_USERNAME"),
		DBName:   os.Getenv("DB_NAME"),
		SSLMode:  os.Getenv("DB_SSLMODE"),
		Password: os.Getenv("DB_PASSWORD"),
	})
	if err != nil {
		log.Fatalf("failed to initialize db: %s", err.Error())
	}
	defer db.Close()

	//Инициализация хеша
	hasher := hasher.NewSHA1Hasher("salt")

	//Инициализация репозиториев
	bananasRepo := pdb.NewBananas(db)
	tokensRepo := pdb.NewTokens(db)
	usersRepo := pdb.NewUsers(db)
	
	//Инициализация RabbitMQ клиента
	audit, err := rmq.NewRabbitMQClient("amqp://guest:guest@localhost:5672/", "logs")
	if err != nil {
		log.Fatal("failed to initialize RabbitMQ client: ", err)
	}
	defer audit.Close()

	//Инициализация сервисов
	bananasService := service.NewBananas(bananasRepo, audit)
	usersService := service.NewUsers(usersRepo, tokensRepo, audit, hasher, []byte("sample secret"))

	//Инициализация хендлера
	handlers := rest.NewHandler(bananasService, usersService)

	//Запуск сервера
	srv := new(server.Server)
	go func() {
		if err := srv.Run(viper.GetString("port"), handlers.InitRoutes()); err != nil {
			log.Fatalf("error occurred while running http server: %s", err.Error())
		}
	}()
	log.Info("Server started")

	//Graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	<-quit
	log.Info("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Errorf("Server forced to shutdown: %s", err.Error())
	}

	if err := db.Close(); err != nil {
		log.Errorf("Error closing database connection: %s", err.Error())
	}

	log.Info("Server exited gracefully")
}

func initLogger() {
	log.SetFormatter(new(log.JSONFormatter))
	log.SetOutput(os.Stdout)
	log.SetLevel(log.InfoLevel)
}

func initConfig() error {
	viper.AddConfigPath("configs")
	viper.SetConfigName("config")
	return viper.ReadInConfig()
}
