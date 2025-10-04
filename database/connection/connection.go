package connection

import (
	"fmt"
	"log"
	"time"

	config "tiketsepur/configs"
	"tiketsepur/utils"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/spf13/viper"
)

var DB *sqlx.DB
var Redis *utils.RedisClient
var RabbitMQ *utils.RabbitMQ

func InitDB(cfg *config.Config) {
    dsn := viper.GetString("DATABASE_URL")
    
    if dsn == "" {
        dsn = fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
            cfg.Database.Host, cfg.Database.Port, cfg.Database.User,
            cfg.Database.Password, cfg.Database.DBName, cfg.Database.SSLMode)
    }

    db, err := sqlx.Connect("postgres", dsn)
	if err != nil {
		log.Fatal("gagal menghubungkan ke database:: ", err)
	}

	db.SetMaxOpenConns(cfg.Database.MaxOpenConns)
	db.SetMaxIdleConns(cfg.Database.MaxIdleConns)
	db.SetConnMaxLifetime(time.Duration(cfg.Database.ConnMaxLifetime))

	DB = db
	log.Println("koneksi ke database sukses")
}

func InitRedis(cfg *config.Config) {
    if cfg.Redis.URL == "" {
        log.Fatal("Redis URL tidak terdefinisi di config.json")
    }
    
    Redis = utils.NewRedisClient(cfg.Redis.URL)
    
    log.Println("Koneksi ke Redis sukses")
}

func InitRabbitMQ(cfg *config.Config) {
	rabbitmq, err := utils.NewRabbitMQ(cfg.RabbitMQ.URL, cfg.RabbitMQ.QueueName)
	if err != nil {
		log.Fatal("Failed to initialize RabbitMQ: ", err)
	}
	
	RabbitMQ = rabbitmq
	rabbitmq.ConsumeNotifications()
	log.Println("Koneksi ke RabbitMQ sukses")
}

func CloseDB() {
	if DB != nil {
		err := DB.Close()
		if err != nil {
			fmt.Println("error saat menutup koneksi database: ", err)
		} else {
			fmt.Println("koneksi database sukses ditutup")
		}
	}
}