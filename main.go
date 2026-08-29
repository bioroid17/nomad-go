package main

import "fmt"

func main() {
	names := []string{"Nico", "Lynn", "Dal"}
	names = append(names, "flynn")
	fmt.Println(names)
}