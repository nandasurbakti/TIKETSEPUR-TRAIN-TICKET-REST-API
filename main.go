package main

import (
	"log"
	"os"
	config "tiketsepur/configs"
	"tiketsepur/database/connection"
	"tiketsepur/database/migrations"
	_ "tiketsepur/docs"
	"tiketsepur/routes"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

// @title TiketSepur API
// @version 1.0
// @description API Server untuk aplikasi pemesanan tiket kereta api TiketSepur
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.url http://www.swagger.io/support
// @contact.email support@swagger.io

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @host localhost:8080
// @BasePath /api/v1
// @schemes http https

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description Type "Bearer" followed by a space and JWT token.
func main() {
	 if os.Getenv("ENVIRONMENT") != "production" {
        godotenv.Load()
    }
	
	cfg, err := config.LoadConfig()
    if err != nil {
        log.Fatal("gagal load config:", err)
    }

    connection.InitDB(cfg)
    defer connection.CloseDB()

    migrations.GetDBMigrate(connection.DB)

    connection.InitRedis(cfg)
    connection.InitRabbitMQ(cfg)

	r := routes.StartServer()


	log.Println("server berjalan di port 8080")
	log.Println("Swagger UI tersedia di http://localhost:8080/swagger/index.html")
	
	if err := r.Run(":8080"); err != nil {
		log.Fatal("gagal menjalankan server:", err)
	}
}