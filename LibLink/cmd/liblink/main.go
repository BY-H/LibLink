package main

import (
	"liblink/internal/global"
	_ "liblink/internal/global"
	"liblink/internal/router"
	"log"
)

func main() {
	r := router.Router()

	if err := r.Run(global.Conf.Port); err != nil {
		log.Fatal("server start error with msg: ", err.Error())
		return
	}
}
