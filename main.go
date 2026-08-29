package main

import "fmt"

func main() {
	nico := map[string]string {
		"name": "Nico",
		"age": "12",
	}
	for key, value := range nico {
		fmt.Println(key, value)
	}
}