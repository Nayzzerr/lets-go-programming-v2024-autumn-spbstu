package main

import (
	"flag"
	"log"
	"os"

	"github.com/EmptyInsid/task-9/cmd/app"
	"github.com/EmptyInsid/task-9/internal/config"
)

func main() {

	CfigPathFlag := flag.String("config", "../../configs/config.yml", "Path to YAML config")
	flag.Parse()

	// Open config file
	configFile, err := os.Open(*CfigPathFlag)
	if err != nil {
		panic(err)
	}

	// Loade date from config file
	cfg, err := config.LoadConfig(configFile)
	if err != nil {
		panic(err)
	}

	log.Printf("Succsess read config file %s", cfg.DBCfg.DBName)

	//make new app with config
	app := app.NewApp(cfg)

	app.Run()
	defer app.Close()

	//start app
	//defer a.Close()
}
