package main

import (
	"fmt"

	"github.com/manuelbanchero/blog-aggregator/internal/config"
)

func main() {
	config := config.Read()

	fmt.Println(config.DbUrl)
}
