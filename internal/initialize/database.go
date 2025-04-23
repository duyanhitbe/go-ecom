package initialize

import (
	"database/sql"
	"fmt"
	"github.com/duyanhitbe/go-ecom/internal/config"
	"github.com/duyanhitbe/go-ecom/internal/global"
	"github.com/duyanhitbe/go-ecom/internal/repositories"
	_ "github.com/lib/pq"
	"log"
)

func InitDatabase() {
	dbCfg := config.Cfg.Database
	connStr := fmt.Sprintf(
		"%s://%s:%s@%s:%d/%s?sslmode=%s",
		dbCfg.DataSource,
		dbCfg.Username,
		dbCfg.Password,
		dbCfg.Host,
		dbCfg.Port,
		dbCfg.DatabaseName,
		dbCfg.SslMode,
	)
	db, err := sql.Open(dbCfg.Driver, connStr)
	if err != nil {
		log.Fatalf("Fail to connect database: %v", err)
	}
	global.Repository = repositories.New(db)

	global.Logger.Info("Connect database successfully!")
}
