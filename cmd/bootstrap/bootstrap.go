package bootstrap

import (
	"log"
	"os"
	"os/signal"

	"github.com/jrmanes/salary-calculator/pkg/server"
)

const (
	infoMsg = "[INFO]"
	errMsg  = "[ERROR]"
)

func Run() error {
	port := os.Getenv("SERVER_PORT")

	serv, err := server.New(port)
	if err != nil {
		log.Fatal(errMsg, err)
	}

	// start the server.
	go serv.Start()

	// Wait for an in interrupt .
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	<-c

	return nil
}
