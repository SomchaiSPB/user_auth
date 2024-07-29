package main

import (
	"context"
	"github.com/SomchaiSPB/user-auth/internal/app"
	"github.com/SomchaiSPB/user-auth/internal/config"
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"
)

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
func main() {
	conf, err := config.Load()

	if err != nil {
		log.Fatal(err)
	}

	a := app.New(conf)

	if err := a.Init(); err != nil {
		log.Fatalf("init application error: %v", err)
	}

	sigCh := make(chan os.Signal)
	signal.Notify(sigCh, syscall.SIGTERM, syscall.SIGINT, syscall.SIGQUIT, syscall.SIGSTOP, syscall.SIGKILL)

	ctx, cancel := context.WithCancel(context.Background())

	var wg = &sync.WaitGroup{}

	log.Println("🚀 application started...")
	go a.Run(ctx, wg)

	select {
	case sig := <-sigCh:
		log.Printf("🚦 received signal: %s \n", sig.String())
		if err := a.ShutDown(); err != nil {
			log.Printf("shutdown error: %v", err)
		}
		cancel()
	}

	wg.Wait()

	log.Println("🛑 application gracefully stopped...")
}
