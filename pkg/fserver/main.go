package main

import (
	"flag"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

func main() {
	configFile := flag.String("config", "config.yaml", "path to config file")
	flag.Parse()

	config := GetFileServerConfig(*configFile)

	engine := gin.New()

	engine.StaticFS("/", http.Dir(config.Root))
	engine.Use(gin.Logger())

	engine.Run(fmt.Sprintf(":%d", config.Port))
}
