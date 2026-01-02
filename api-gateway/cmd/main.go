package main

import (
	"log"

	"wealth-vault/api-gateway/configs"
	"wealth-vault/api-gateway/internal/bootstrap"
)

func main() {
	cfg := configs.LoadConfigs()

	app := bootstrap.InitApp(cfg)

	log.Println("API Gateway running on :", cfg.APP.Port)
	log.Fatal(app.Listen(":" + cfg.APP.Port))
}
