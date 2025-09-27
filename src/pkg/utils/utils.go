package utils

import (
	"fmt"
	"log"
)

func CheckErr(err error, args ...interface{}) {
	if err != nil {
		panic(fmt.Sprint(append([]interface{}{err}, args...)...))
	}
}

func LogErr(err error, args ...interface{}) {
	if err != nil {
		log.Println(fmt.Sprint(append([]interface{}{err}, args...)...))
	}
}
