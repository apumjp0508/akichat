package router

import (
    "fmt"
    wsHandler "akichat/backend/internal/handler/webSocket"
    UserHandler "akichat/backend/internal/handler/userHandler"
    middleware "akichat/backend/internal/http/middleware"
    friendsHandler "akichat/backend/internal/handler/friendsHandler"
    jwtTokenHandler "akichat/backend/internal/handler/auth/token/JWTToken"
    sessionHandler "akichat/backend/internal/handler/auth/session"
    "akichat/backend/internal/repository"
    "akichat/backend/internal/db"
    "akichat/backend/internal/config"
    "github.com/gin-contrib/cors"
    "github.com/gin-gonic/gin"
    "github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
)	

func SetupRouter() *gin.Engine {
	r := gin.Default()

	database, err := db.InitDB()
    if err != nil {
        fmt.Printf("db接続失敗")
    }

    cfg := config.Load()

    store := cookie.NewStore([]byte(cfg.SessionSecretKey))
    store.Options(sessions.Options{
        Path:     cfg.SessionPath,
        HttpOnly: cfg.SessionHttpOnly,
        Secure:   cfg.SessionSecure,
        SameSite: cfg.SessionSameSite,
        MaxAge:   cfg.SessionMaxAge,
    })

    r.Use(sessions.Sessions("mysession", store))
	
    r.Use(cors.New(cors.Config{
        AllowOrigins:     cfg.CorsOrigins,
        AllowMethods:     cfg.CorsMethods,
        AllowHeaders:     cfg.CorsHeaders,
        ExposeHeaders:    cfg.CorsExpose,
        AllowCredentials: cfg.CorsCredentials,
    }))

    userRepo := repository.NewUserRepository(database)
    friendsRepo := repository.NewFriendShipRepository(database)
    friendRequestRepo := repository.NewFriendRequestRepository(database)

    registerHandler := UserHandler.NewRegisterHandler(userRepo)
    loginHandler := UserHandler.NewLoginHandler(userRepo)
    getMeHandler := UserHandler.NewGetMeHandler(userRepo)

    FetchFriendsHandler := friendsHandler.NewFriendShipHandler(friendsRepo)
    searchNotFriendHandler := friendsHandler.NewSearchNotFriendHandler(userRepo)
    friendRequestHandler := friendsHandler.NewFriendRequestHandler(friendRequestRepo)

    r.POST("/api/register", registerHandler.RegisterHandler)
    r.POST("/api/login", loginHandler.LoginHandler)
    r.POST("/api/auth/refresh", jwtTokenHandler.RefreshHandler)

    auth:= r.Group("/api")
    auth.Use(middleware.JWTMiddleware())
    {
        auth.POST("/websocket/init", sessionHandler.SetupSessionHandler)
        auth.GET("/connected-users", wsHandler.GetConnectedUsersHandler)
        auth.GET("/getMe", getMeHandler.GetMeHandler)
        auth.GET("/fetch/friends" ,FetchFriendsHandler.GetFriendsHandler) 
        auth.POST("/search/notfriends", searchNotFriendHandler.SearchNotFriendHandler)
        auth.POST("/friend/request", friendRequestHandler.FriendRequestHandler)
        auth.POST("/friend/request/approve", FetchFriendsHandler.ApproveRequestHandler)
    }

    wsGroup := r.Group("/api/session")
    wsGroup.Use(middleware.SessionMiddleware())
    wsGroup.GET("/websocket", wsHandler.WebSocketHandler)

	fmt.Println("http://localhost:8080 で起動中")

	return r
}	

