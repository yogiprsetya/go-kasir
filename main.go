package main

import (
	"fmt"
	"go-kasir-api/database"
	"go-kasir-api/handlers"
	"go-kasir-api/repositories"
	"go-kasir-api/services"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/spf13/viper"
)

type Config struct {
	Port   string `mapstructure:"PORT"`
	DBConn string `mapstructure:"DB_CONN"`
}

func main() {
	// Init ENV using viper
	viper.AutomaticEnv()
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "-"))

	if _, err := os.Stat(".env"); err == nil {
		viper.SetConfigFile(".env")
		_ = viper.ReadInConfig()
	}

	config := Config{
		Port:   viper.GetString("PORT"),
		DBConn: viper.GetString("DB_CONN"),
	}

	// Setup database
	db, err := database.InitDB(config.DBConn)
	if err != nil {
		log.Fatal("Failed to initialize database:", err)
	}
	defer db.Close()

	productRepo := repositories.NewProductRepository(db)
	productService := services.NewProductService(productRepo)
	productHandler := handlers.NewProductHandler(productService)

	categoryRepo := repositories.NewCategoryRepository(db)
	categoryService := services.NewCategoryService(categoryRepo)
	categoryHandler := handlers.NewCategoryHandler(categoryService)

	// Product routes
	http.HandleFunc("/api/product", productHandler.HandleProducts)
	http.HandleFunc("/api/product/", productHandler.HandleProductById)

	// Category routes
	http.HandleFunc("/api/category", categoryHandler.HandleCategory)
	http.HandleFunc("/api/category/", categoryHandler.HandleCategoryById)

	fmt.Println("Server running on http://localhost:" + config.Port)
	httpErr := http.ListenAndServe(":"+config.Port, nil)

	if httpErr != nil {
		fmt.Println("Server failed:", httpErr)
	}
}
