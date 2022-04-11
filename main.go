package main

import (
	"fmt"
	"go-json/json"
)

type Person struct {
	H int
}

type User struct {
	Name,Address string
	Age int `json:"gender"`
	People []int `json:"haha"`
	P []Person
}

func main()  {
	var u User = User{
		Name:    "ddd",
		Address: "wdwd",
		Age:     0,
		People:  []int{1, 2, 3},
		P: []Person{{H:1}, {H:2}},
	}

	fmt.Println(u)
	by, _ := json.Marshal(u)
	fmt.Println(string(by))
	ii := []int{1,2,3,4}
	by, _ = json.Marshal(ii)
	fmt.Println(string(by))
}
