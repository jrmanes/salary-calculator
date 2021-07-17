package main

import (
	"log"

	"github.com/jrmanes/my-salary/cmd/bootstrap"
)

const (
	infoMsg = "[INFO]"
	errMsg  = "[ERROR]"
)

func main() {
	log.Println(infoMsg, "Stating service...")
	err := bootstrap.Run()
	if err != nil {
		log.Println(errMsg, err)
	}
}
