package routes

import (
	"log"
	config "tiketsepur/configs"
	"tiketsepur/controllers"
	"tiketsepur/database/connection"
	"tiketsepur/middleware"
	"tiketsepur/repository"
	"tiketsepur/service"

	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	"github.com/gin-gonic/gin"
)

func StartServer() *gin.Engine{
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatal("gagal melakukan konfigurasi: ", err)
	}
	
	userRepo := repository.NewUserRepository(connection.DB)
	trainRepo := repository.NewTrainRepository(connection.DB)
	scheduleRepo := repository.NewScheduleRepository(connection.DB)
	ticketRepo := repository.NewTicketRepository(connection.DB)
	paymentRepo := repository.NewPaymentRepository(connection.DB)

	authService := service.NewAuthService(userRepo, connection.Redis, cfg)
	userService := service.NewUserService(userRepo)
	trainService := service.NewTrainService(trainRepo)
	scheduleService := service.NewScheduleService(scheduleRepo, trainRepo)
	ticketService := service.NewTicketService(connection.DB, ticketRepo, scheduleRepo, userRepo, paymentRepo, connection.Redis, connection.RabbitMQ)
	paymentService := service.NewPaymentService(connection.DB, paymentRepo, ticketRepo, scheduleRepo, userRepo, connection.RabbitMQ)

	authControllers := controllers.NewAuthControllers(authService, userService)
	userControllers := controllers.NewUserControllers(userService)
	trainControllers := controllers.NewTrainControllers(trainService)
	scheduleControllers := controllers.NewScheduleControllers(scheduleService)
	ticketControllers := controllers.NewTicketControllers(ticketService)
	paymentControllers := controllers.NewPaymentHandler(paymentService)

	if cfg.Server.Mode == "release" {
		gin.SetMode(gin.ReleaseMode)
	}

	r := gin.Default()

	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	r.Use(func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}
		
		c.Next()
	})
	
	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status": "healthy",
			"service": "train-ticket-api",
		})
	})
	
	api := r.Group("/api")
	{
		public := api.Group("/public")
		{
			public.GET("schedules", scheduleControllers.GetAll)
			public.GET("/search", scheduleControllers.Search)
			public.GET("/:id", scheduleControllers.GetByID)
		}
		auth := api.Group("/auth")
		{
			auth.POST("/register", authControllers.Register)
			auth.POST("/register-admin", authControllers.RegisterAdmin)
			auth.POST("/login", authControllers.Login)
		}

		authenticated := api.Group("")
		authenticated.Use(middleware.JWTAuthMiddleware(authService))
		{

			authenticated.POST("/auth/logout", authControllers.Logout)
			authenticated.GET("/auth/me", authControllers.Me)

			tickets := authenticated.Group("/tickets")
			{
				tickets.POST("", ticketControllers.Create)
				tickets.GET("/my-tickets", ticketControllers.GetMyTickets)
				tickets.GET("/:id", ticketControllers.GetByID)
				tickets.PUT("/:id/cancel", ticketControllers.Cancel)
			}

			payments := authenticated.Group("/payments")
			{
				payments.POST("/confirm/:paymentCode", paymentControllers.ConfirmPayment)
				payments.GET("/status/:paymentCode", paymentControllers.GetPaymentStatus)
			}

			admin := authenticated.Group("")
			admin.Use(middleware.AdminOnly())
			{
				users := admin.Group("/users")
				{
					users.POST("", userControllers.Create)
					users.GET("", userControllers.GetAll)
					users.GET("/:id", userControllers.GetByID)
					users.PUT("/:id", userControllers.Update)
					users.DELETE("/:id", userControllers.Delete)
				}

				trains := admin.Group("/trains")
				{
					trains.POST("", trainControllers.Create)
					trains.GET("", trainControllers.GetAll)
					trains.GET("/:id", trainControllers.GetByID)
					trains.PUT("/:id", trainControllers.Update)
					trains.DELETE("/:id", trainControllers.Delete)
				}

				adminSchedules := admin.Group("/schedules")
				{
					adminSchedules.POST("", scheduleControllers.Create)
					adminSchedules.PUT("/:id", scheduleControllers.Update)
					adminSchedules.DELETE("/:id", scheduleControllers.Delete)
				}

				adminTickets := admin.Group("/tickets")
				{
					adminTickets.GET("/all", ticketControllers.GetAll)
				}
			}
		}
	}
	
	return r

}