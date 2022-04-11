package main

import (
	"fmt"
	"go-json/json"
)

type User struct {
	Name,Address string
	Age int `json:"gender"`
	People []int `json:"haha"`
}

func main()  {
	var u User = User{
		Name:    "ddd",
		Address: "wdwd",
		Age:     0,
		People:  []int{1, 2, 3},
	}
	by, _ := json.Marshal(u)
	fmt.Println(string(by))
	ii := []int{1,2,3,4}
	by, _ = json.Marshal(ii)
	fmt.Println(string(by))
}
