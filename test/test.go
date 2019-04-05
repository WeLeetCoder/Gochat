package main

import (
	"Gochat/model"
	"fmt"
)

func main() {
	_, err := model.GetToken("e02486ffffcb679bb86a472b3f654bd19438")
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(model.AuthToken("e02486cb679bb86a472b3f654bd19438"))
}
