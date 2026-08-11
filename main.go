package main

import (
	"fmt"

	"github.com/manuelbanchero/blog-aggregator/internal/config"
)

func main() {
	config, err := config.Read()
	if err != nil {
		fmt.Printf("%w", err)
		return
	}

	fmt.Println(config.DbUrl)
}
