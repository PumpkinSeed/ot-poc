package main

import (
	"log/slog"
	"math/rand"
	"net/http"
	"time"
)

func main() {
	for {
		if _, err := http.Get("http://localhost:8080/multi"); err != nil {
			slog.Error(err.Error())
		}
		sleepTime := time.Duration(300+rand.Intn(500)) * time.Millisecond
		time.Sleep(sleepTime)
		slog.Info("Request sent")
	}
}
