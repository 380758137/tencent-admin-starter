package main

import (
	"log"

	"github.com/liusheng/tencent-admin-starter/apps/api/internal/config"
	"github.com/liusheng/tencent-admin-starter/apps/api/internal/database"
	"github.com/liusheng/tencent-admin-starter/apps/api/internal/router"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("load config failed: %v", err)
	}

	db, err := database.New(cfg)
	if err != nil {
		log.Fatalf("init database failed: %v", err)
	}

	r := router.New(cfg, db)
	if err := r.Run(":" + cfg.Port); err != nil {
		log.Fatalf("start server failed: %v", err)
	}
}
