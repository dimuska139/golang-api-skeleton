package main

import (
	"context"
	"flag"
	"fmt"
	"github.com/dimuska139/golang-api-skeleton/middlewares"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	configPathPtr := flag.String("config", "config.yml", "Path to configuration file")
	flag.Parse()

	cfg, err := InitializeConfig(*configPathPtr)
	if err != nil {
		fmt.Println(err)
	}

	db, err := InitializeDatabase(cfg)
	if err != nil {
		fmt.Println(err)
	}

	usersApi, err := InitializeUsersAPI(db)
	if err != nil {
		fmt.Println(err)
	}

	authApi, err := InitializeAuthAPI(cfg, db)
	if err != nil {
		fmt.Println(err)
	}

	router := gin.Default()
	api := router.Group("/v1")
	{
		users := api.Group("/users")
		{
			users.GET("/total", usersApi.GetTotal)
			users.GET("", usersApi.GetList)
		}
		auth := api.Group("/auth")
		{
			auth.POST("/refresh-tokens", authApi.RefreshTokens)
			auth.POST("/login", authApi.Login)
			auth.POST("/registration", authApi.Registration)
		}
		private := api.Group("/private") // Authentication required
		{
			private.Use(middlewares.JwtMiddleware(cfg))
			private.GET("/profile", usersApi.GetProfile)
		}
	}

	srv := &http.Server{
		Addr:    ":8080",
		Handler: router,
	}
	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
	}()

	quit := make(chan os.Signal)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("Server forced to shutdown:", err)
	}

	log.Println("Server exiting")
}
