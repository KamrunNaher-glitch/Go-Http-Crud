package main

import "fmt"

var users = []string{
	"Mizan",
	"Mezba",
	"Mir",
	"firoz",
}

func main() {
	idx := 1

	users = append(users[:idx], users[idx+1:]...)
	fmt.Println(users)
}