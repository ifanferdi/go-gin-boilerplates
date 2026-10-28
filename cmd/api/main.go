package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"go-gin-boilerplate/internal/bootstrap"
	"go-gin-boilerplate/internal/config"
	"go-gin-boilerplate/internal/infrastructure/database"
	"go-gin-boilerplate/internal/infrastructure/database/migrations"
	"go-gin-boilerplate/internal/infrastructure/database/seeders"
)

func main() {
	migrateFlag := flag.Bool("migrate", false, "Run pending database migrations")
	rollbackFlag := flag.Bool("rollback", false, "Rollback the last migration batch")
	dropFlag := flag.Bool("drop", false, "Drop all database tables")
	freshFlag := flag.Bool("fresh", false, "Drop all tables and re-run migrations")
	seedFlag := flag.Bool("seed", false, "Run database seeders")
	runFlag := flag.Bool("run", false, "Keep HTTP & gRPC servers running after DB commands")
	flag.Parse()

	cfg := config.Load()

	db, err := database.NewPostgresConnection(cfg)
	if err != nil {
		log.Fatalf("❌ Database error: %v", err)
	}

	dbCommandExecuted := false

	if *dropFlag {
		dbCommandExecuted = true
		if err := migrations.DropAll(db); err != nil {
			log.Fatalf("❌ Drop failed: %v", err)
		}
	}

	if *freshFlag {
		dbCommandExecuted = true
		if err := migrations.Fresh(db); err != nil {
			log.Fatalf("❌ Fresh migration failed: %v", err)
		}
	} else if *migrateFlag {
		dbCommandExecuted = true
		if err := migrations.Migrate(db); err != nil {
			log.Fatalf("❌ Migration failed: %v", err)
		}
	}

	if *rollbackFlag {
		dbCommandExecuted = true
		if err := migrations.Rollback(db); err != nil {
			log.Fatalf("❌ Rollback failed: %v", err)
		}
	}

	if *seedFlag {
		dbCommandExecuted = true
		if err := seeders.RunAllSeeders(db); err != nil {
			log.Fatalf("❌ Seeding failed: %v", err)
		}
	}

	if dbCommandExecuted && !*runFlag {
		return
	}

	app, err := bootstrap.NewApp(cfg, db)
	if err != nil {
		log.Fatalf("❌ Bootstrap failed: %v", err)
	}
	defer app.Close()

	errChan := make(chan error, 2)

	go func() {
		if err := app.GRPCServer.Start(cfg.GRPC.Port); err != nil {
			errChan <- err
		}
	}()

	go func() {
		addr := fmt.Sprintf(":%s", cfg.App.Port)
		log.Printf("🚀 HTTP Server listening on %s\n", addr)
		if err := app.HTTPServer.Run(addr); err != nil {
			errChan <- err
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	select {
	case sig := <-quit:
		log.Printf("🛑 Shutting down servers (signal: %v)...\n", sig)
	case err := <-errChan:
		log.Fatalf("❌ Server error: %v", err)
	}
}
