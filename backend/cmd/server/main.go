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
    log.Println("http://0.0.0.0:" + cfg.Port + " で起動中です")
    _ = r.Run("0.0.0.0:" + cfg.Port)
}
