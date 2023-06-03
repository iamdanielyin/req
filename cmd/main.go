package main

import (
	"github.com/iamdanielyin/req"
	"log"
)

func main() {
	var posts []struct {
		UserId int    `json:"userId"`
		Id     int    `json:"id"`
		Title  string `json:"title"`
		Body   string `json:"body"`
	}
	if err := req.GET("https://jsonplaceholder.typicode.com/posts", &posts); err != nil {
		log.Fatal(err)
	}
}
