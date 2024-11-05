package application

import (
	"context"
	"errors"
	"fmt"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"iditusi/internal/controllers"
	"iditusi/internal/repositories/postgres"
	"iditusi/internal/response"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/limiter"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/jackc/pgx/v5/pgxpool"
)

const version = "0.1.1"

func Run(ctx context.Context, args []string) error {

}