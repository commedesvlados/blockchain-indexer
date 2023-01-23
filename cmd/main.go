package main

import (
	"context"
	"fmt"
	"github.com/commedesvlados/blockchain-indexer/pkg/handler"
	"github.com/commedesvlados/blockchain-indexer/pkg/repository"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"github.com/spf13/viper"
	"log"
	"os"
	"time"
)

const (
	//httpLocalURL = "http://localhost:8545"
	//wsLocalURL   = "ws://localhost:8545"
	sepoliaURL = "wss://sepolia.infura.io/ws/v3/3b7a103e898c41c1960043a2dc3cf6ca"
)

func main() {

	// cfg
	if err := initConfig(); err != nil {
		log.Fatalf("error initializing configs: %v", err)
	}

	fmt.Println("config initialized")

	// env
	if err := godotenv.Load(); err != nil {
		log.Fatalf("error loading env variables: %v", err)
	}

	fmt.Println("dotenv initialized")

	// db
	db, err := repository.NewPostgresDB(repository.Config{
		Host:     viper.GetString("db.host"),
		Port:     viper.GetString("db.port"),
		Username: viper.GetString("db.username"),
		DBName:   viper.GetString("db.dbname"),
		SSLMode:  viper.GetString("db.sslmode"),
		Password: os.Getenv("DB_PASSWORD"),
	})
	if err != nil {
		log.Fatalf("failed to initialize db: %s", err.Error())
	}

	fmt.Println("database initialized")

	// eth
	client, err := ethclient.DialContext(context.Background(), sepoliaURL)
	if err != nil {
		log.Fatalf("error initializing ethereum client: %v", err)
	}

	fmt.Println("connected to blockchain node")

	repos := repository.NewRepositiry(db)
	handlers := handler.NewHandler(repos, client)

	// subscribe to blockchain
	go func() {
		fmt.Println("go routine start")
		handlers.StartListen(context.Background())
	}()

	// TODO shutdown
	time.Sleep(time.Second * 40)

	fmt.Println("\nApp close")
}

func initConfig() error {
	viper.SetConfigType("yaml")
	viper.SetConfigName("config")
	viper.AddConfigPath("configs")
	return viper.ReadInConfig()
}
