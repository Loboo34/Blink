package main

import (
	"fmt"
	"log"

	"github.com/joho/godotenv"
)


func init(){
	err := godotenv.Load()
	if err != nil{
		log.Fatal("Failed to load env")
	}
}
func main(){
	fmt.Println("Hello")
}