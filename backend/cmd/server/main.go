package main

import (
	"log"
	router "akichat/backend/internal/http"
    "akichat/backend/internal/config"
)



func main() {
    cfg := config.Load()
    r := router.SetupRouter()
    log.Println("http://localhost:" + cfg.Port + " で起動中")
    _ = r.Run(":" + cfg.Port)
}
