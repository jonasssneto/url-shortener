package config

import (
	"flag"
	"log"

	"github.com/joho/godotenv"
	"github.com/kelseyhightower/envconfig"
)

type Config struct {
	Postgres struct {
		URI string `envconfig:"POSTGRES_URI" default:"postgres://user:password@localhost:5432/dbname?sslmode=disable"`
	}

	Otel struct {
		URL string `envconfig:"OTEL_COLLECTOR_URL" default:"localhost:4317"`
	}

	Development bool
}

var Env Config

func Load() error {
	production := flag.Bool("production", false, "Run app on Production mode")
	flag.Parse()

	if !(*production) {
		log.Println("Development mode")

		if err := godotenv.Load(".env"); err != nil {
			return err
		}

		log.Println("Loaded .env file")
	}

	if err := envconfig.Process("", &Env); err != nil {
		log.Fatal("LoadConfig - Erro ao processar variáveis de ambiente", err)
		return err
	}

	Env.Development = !(*production)

	if Env.Development {
		log.Printf("Configurações: %+v \n", Env)
	}

	return nil
}
