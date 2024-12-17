package main

import (
	"backend/service"
	"database/sql"
	"log/slog"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
)

type App struct {
	Addr      string
  Port      int
	Router    *gin.Engine
	Service   service.Service
	Database  *sql.DB
  Redis     *redis.Client
	Logger    *slog.Logger
}