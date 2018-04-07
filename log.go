package main

import (
	"log"
)

func Debug(i interface{}) {
	log.Println("[Debug]", i)
}

func Error(i interface{}) {
	log.Println("[Error]", i)
}
