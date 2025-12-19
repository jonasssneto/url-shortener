package main

import "main/internal/config"

func main() {
	if err := config.Load(); err != nil {
		panic(err)
	}
}
