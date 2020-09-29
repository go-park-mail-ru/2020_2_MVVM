package main

import (
	"github.com/go-park-mail-ru/2020_2_MVVM.git/pkg/api"
	yconfig "github.com/rowdyroad/go-yaml-config"
)

func main() {
	var config api.Config
	yconfig.LoadConfig(&config, "configs/config.yaml", nil)
	app := api.NewApp(config)
	defer app.Close()
	app.Run()
}
