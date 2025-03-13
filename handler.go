package main

import (
	"fmt"
	"net/http"
)

type HelloHandler struct{}

func NewHelloHandler(router *http.ServeMux) {
	handler := &HelloHandler{}
	router.HandleFunc("/hello", handler.HelloToEveryone())
	router.HandleFunc("/hello/Ivan", handler.HelloToMe())
}

func (h *HelloHandler) HelloToEveryone() http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		fmt.Println("hello to everyone")
	}
}

func (h *HelloHandler) HelloToMe() http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		fmt.Println("Hello to Me")
	}
}
