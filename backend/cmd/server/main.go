package main

import (
	"log"
	router "akichat/backend/internal/http"
    websocket "akichat/backend/internal/handler/webSocket"
    "akichat/backend/internal/config"
)



func main() {
    go websocket.GlobalHub.Run()
    cfg := config.Load()
    r := router.SetupRouter()
    log.Println("http://localhost:" + cfg.Port + " で起動中")
    _ = r.Run(":" + cfg.Port)
}
