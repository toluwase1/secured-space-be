package server

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/decadevs/rentals-api/db"
	"github.com/decadevs/rentals-api/router"
	"github.com/decadevs/rentals-api/server/middleware"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

// Server serves requests to DB with router
type Server struct {
	DB     db.DB
	Mail   db.Mailer
	Router *router.Router
}

func (s *Server) defineRoutes(router *gin.Engine) {
	apirouter := router.Group("/api/v1")
	apirouter.POST("/auth/signup_tenant", s.handleSignupTenant())
	apirouter.POST("/auth/signup_agent", s.handleSignupAgent())
	apirouter.POST("/auth/login", s.handleLogin())
	apirouter.GET("/features/interior", s.handleGetInteriorFeatures())
	apirouter.GET("/features/exterior", s.handleGetExteriorFeatures())
	apirouter.GET("/categories", s.handleGetCategories())
	apirouter.POST("/reset-password/:userID", s.ResetPassword())
	apirouter.GET("/search-apartment", s.SearchApartment())
	apirouter.GET("/apartment-details/:apartmentID", s.GetApartmentDetails())
	apirouter.POST("/new/user", s.registerNewUser())
	apirouter.POST("/pusher/auth", s.pusherAuth())
	apirouter.POST("/chat/create", s.CreateChat())
	apirouter.POST("/pusher/message", s.SendNewMessage())
	apirouter.GET("/apartment", s.GetAllApartments())
	apirouter.POST("/verify-email/:userID/:userToken", s.VerifyEmail())
	apirouter.POST("/forgot-password", s.ForgotPassword())

	authorized := apirouter.Group("/")
	authorized.Use(middleware.Authorize(s.DB.FindUserByEmail, s.DB.TokenInBlacklist))
	authorized.POST("/logout", s.handleLogout())
	authorized.GET("/users", s.handleGetUsers())
	authorized.GET("/bookmark/apartments", s.GetBookmarkedApartments())
	authorized.GET("/user-apartment", s.handleGetUserApartments())
	authorized.PUT("/me/update", s.handleUpdateUserDetails())
	authorized.GET("/me", s.handleShowProfile())
	authorized.PUT("/me/uploadpic", s.handleUploadProfilePic())
	authorized.POST("/user/change-password", s.ChangePassword())
	// apartment routes
	authorized.POST("/user/apartments", s.handleCreateApartment())
	authorized.DELETE("/user/apartment/:apartmentID", s.DeleteApartment())
	authorized.PUT("/user/:apartmentID/update", s.handleUpdateApartmentDetails())
	authorized.GET("/user/:apartmentID/bookmark", s.SaveBookmarkApartment())
	authorized.DELETE("/user/apartment/:apartmentID/removebookmark", s.RemoveBookmarkedApartment())
	authorized.GET("/user/:apartmentID", s.handleGetApartmentByID())

}

func (s *Server) setupRouter() *gin.Engine {
	ginMode := os.Getenv("GIN_MODE")
	if ginMode == "test" {
		r := gin.New()
		s.defineRoutes(r)
		return r
	}

	r := gin.New()
	// LoggerWithFormatter middleware will write the logs to gin.DefaultWriter
	// By default gin.DefaultWriter = os.Stdout
	r.Use(gin.LoggerWithFormatter(func(param gin.LogFormatterParams) string {
		// your custom format
		return fmt.Sprintf("%s - [%s] \"%s %s %s %d %s \"%s\" %s\"\n",
			param.ClientIP,
			param.TimeStamp.Format(time.RFC1123),
			param.Method,
			param.Path,
			param.Request.Proto,
			param.StatusCode,
			param.Latency,
			param.Request.UserAgent(),
			param.ErrorMessage,
		)
	}))
	r.Use(gin.Recovery())
	// setup cors
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"POST", "GET", "PUT", "PATCH", "DELETE"},
		AllowHeaders:     []string{"*"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		// AllowOriginFunc: func(origin string) bool {
		// 	return origin == "https://github.com"
		// },
		MaxAge: 12 * time.Hour,
	}))
	s.defineRoutes(r)
	return r
}

// Start starts the whole server by preparing everything it needs
// like router
func (s *Server) Start() {
	r := s.setupRouter()
	PORT := fmt.Sprintf(":%s", os.Getenv("PORT"))
	if PORT == ":" {
		PORT = ":8080"
	}
	srv := &http.Server{
		Addr:    PORT,
		Handler: r,
	}

	// Initializing the server in a goroutine so that
	// it won't block the graceful shutdown handling below
	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
	}()

	log.Printf("Server started on %s\n", PORT)

	s.DB.PopulateTables()
	// Wait for interrupt signal to gracefully shutdown the server with
	// a timeout of 5 seconds.
	quit := make(chan os.Signal)
	// kill (no param) default send syscall.SIGTERM
	// kill -2 is syscall.SIGINT
	// kill -9 is syscall.SIGKILL but can't be catch, so don't need add it
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutting down server...")

	// The context is used to inform the server it has 5 seconds to finish
	// the request it is currently handling
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("Server forced to shutdown:", err)
	}

	log.Println("Server exiting")
}
