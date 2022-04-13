package main

import (
	//j "encoding/json"
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
	//P []Person
	M map[string]int
}

func main()  {
	m := make(map[string]int, 4)
	m["ccc"] = 2
	m["1"] = 1
	var u User = User{
		Name:    "ddd",
		Address: "wdwd",
		Age:     0,
		People:  []int{1, 2, 3},
		//P: []Person{{H:1}, {H:2}},
		M: m,
	}

	//by, _ :=j.Marshal(u)
	//fmt.Println(string(by))

	by, _ := json.Marshal(u)
	fmt.Println(string(by))

	//var u1 User
	//
	//_ = json.Unmarshal(by, &u1)
	//
	//fmt.Println(u1)
}
