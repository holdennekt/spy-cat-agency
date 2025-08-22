package cmd

import (
	"fmt"
	"os"
	"spy-cat-agency/config"
	"spy-cat-agency/internal/database"
	"spy-cat-agency/internal/handler"
	"spy-cat-agency/internal/middleware"
	"spy-cat-agency/internal/repository"
	"spy-cat-agency/internal/service"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/spf13/cobra"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

var rootCmd = &cobra.Command{
	Use:   "spy-cat-agency",
	Short: "Spy Cat Agency App",
	Run:   run,
}

func Execute() error {
	return rootCmd.Execute()
}

func run(cmd *cobra.Command, args []string) {
	cfg := config.Load()

	db, err := database.Connect(cfg.DatabaseURL)
	if err != nil {
		fmt.Printf("Failed to connect to database: %v\n", err)
		os.Exit(1)
	}

	err = database.Migrate(db)
	if err != nil {
		fmt.Printf("Failed to run migrations: %v\n", err)
		os.Exit(1)
	}

	repos := repository.New(db)
	services := service.New(repos, cfg)
	handlers := handler.New(services)

	// gin includes logger middleware by default
	r := gin.Default()

	corsConfig := cors.DefaultConfig()
	corsConfig.AllowOrigins = []string{"http://localhost"}
	corsConfig.AllowCredentials = true
	r.Use(cors.New(corsConfig))

	r.Use(middleware.ErrorHandler())

	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	api := r.Group("/api/v1")

	cats := api.Group("/cats")
	cats.POST("", handlers.CreateCat)
	cats.GET("", handlers.GetCats)
	cats.GET("/:id", handlers.GetCat)
	cats.PATCH("/:id", handlers.UpdateCat)
	cats.DELETE("/:id", handlers.DeleteCat)

	missions := api.Group("/missions")
	missions.POST("", handlers.CreateMission)
	missions.GET("", handlers.GetMissions)
	missions.GET("/:id", handlers.GetMission)
	missions.PATCH("/:id", handlers.UpdateMission)
	missions.DELETE("/:id", handlers.DeleteMission)
	missions.PATCH("/:id/assign/:cat_id", handlers.AssignCatToMission)
	missions.POST("/:id/targets", handlers.CreateTarget)
	missions.PATCH("/:id/targets/:target_id", handlers.UpdateTarget)
	missions.DELETE("/:id/targets/:target_id", handlers.DeleteTarget)

	fmt.Printf("Server starting on port %s\n", cfg.Port)
	r.Run(":" + cfg.Port)
}
