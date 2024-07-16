package otpoc

import (
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"testing"
	"time"
)

func TestRun(t *testing.T) {
	go Run()

	time.Sleep(1 * time.Second)

	resp, err := http.Get("http://localhost:8080/multi")
	if err != nil {
		t.Fatal(err)
	}
	data, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(string(data))

	time.Sleep(2 * time.Second)
}

func TestReq(t *testing.T) {
	for {
		if _, err := http.Get("http://localhost:8080/multi"); err != nil {
			t.Error(err)
		}
		sleepTime := time.Duration(300+rand.Intn(500)) * time.Millisecond
		time.Sleep(sleepTime)
		t.Log("Request sent")
	}

}
