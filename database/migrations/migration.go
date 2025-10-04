package migrations

import (
	"embed"
	"fmt"
	"log"

	"github.com/jmoiron/sqlx"
	migrate "github.com/rubenv/sql-migrate"
)

//go:embed sql_migrations/*.sql
var dbMigrations embed.FS

func GetDBMigrate(dbParam *sqlx.DB) {
	migrations := &migrate.EmbedFileSystemMigrationSource{
		FileSystem: dbMigrations,
		Root: "sql_migrations",
	}

	n, err := migrate.Exec(dbParam.DB, "postgres", migrations, migrate.Up)
	if err != nil {
		log.Fatal("migrate gagal: ", err)
	}
	
	fmt.Println("Migration sucess, applied", n, "migrations!")
}