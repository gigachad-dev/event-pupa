package main

import (
	"context"
	"eventloop/internal/httpapi"
	"eventloop/internal/logger"
	"fmt"
	"os"
	"os/signal"
	"strings"
	"syscall"
)

const _LOGLEVEL = "debug"

func inputMonitor(sc chan<- os.Signal) {
	var input string
	for strings.ToLower(input) != "stop" {
		_, _ = fmt.Scanln(&input)
	}
	sc <- syscall.SIGTERM
}

func main() {
	srvLogger, err := initLogger()
	if err != nil {
		fmt.Println(err)
		return
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		errServer := httpapi.StartServer(8090, srvLogger)
		if errServer != nil {
			fmt.Println(err)
			return
		}
	}()

	fmt.Println("Server started...")

	go inputMonitor(sc)
	<-sc

	fmt.Println("Server stopping...")
	if errStop := httpapi.StopServer(ctx, srvLogger); err != nil {
		fmt.Println(errStop)
	}

	fmt.Println("Server stopped.")
}

func initLogger() (logger.Interface, error) {
	appLogger, err := logger.NewLogger(_LOGLEVEL, "logs", "")
	if err != nil {
		return nil, err
	}

	appLogger.Infof("Server starting...")
	return appLogger, nil
}
