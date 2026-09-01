package main

import (
	"fmt"
	"nomad-go/accounts"
)

func main() {
	account := accounts.NewAccount("nico")
	fmt.Println(account)
}