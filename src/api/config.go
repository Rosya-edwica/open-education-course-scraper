package api

import (
	"io"
	"net/http"
	"fmt"
	"time"
)


func getJson(url string) string {
	client := http.Client{
		Timeout: 60 * time.Second,
	}
	response, err := client.Get(url)
	checkErr(err)
	body, err := io.ReadAll(response.Body) 
	checkErr(err)
	return string(body)
}

func checkErr(err error) {
	if err != nil {
		panic(fmt.Sprintf("[%s]", err.Error()))
	}
}