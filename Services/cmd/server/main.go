package main

import (
	"flag"
	"log"
	"os"

	server "github.com/ChangerzaryX1602/SkillSync/internal/infrastructure"
	"github.com/ChangerzaryX1602/SkillSync/pkg/config"
)

var (
	version string
	build   string
	runEnv  string
)

func init() {
	// read running flag
	if len(os.Getenv("ENV")) != 0 {
		runEnv = os.Getenv("ENV")
	} else {
		flagEnv := flag.String("env", "dev", "A config file name without .env")
		flag.Parse()
		runEnv = *flagEnv
	}
	// load config by running flag
	if err := config.LoadConfig(runEnv); err != nil {
		log.Fatalf("error while loading the env:\n %+v", err)
	}
}

// @title SkillSync API
// @version 1.0
// @description This is a sample server for Test Services.
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.url http://www.swagger.io/support
// @contact.email support@swagger.io

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @host localhost:8080
// @BasePath /api/v1
// @schemes http https

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
func main() {

	// init server
	server, err := server.NewServer(version, build, runEnv)
	if err != nil {
		log.Fatalf("error while create server:\n %+v", err)
	}

	server.Run()
}
