package main

import (
	"fmt"

	"github.com/manuelbanchero/blog-aggregator/internal/config"
)

func main() {
	cfg, err := config.Read()
	if err != nil {
		fmt.Println(err)
		return
	}

	cfg.CurrentUserName = "mbanchero"
	cfg.SetUser()

	updatedConfig, err := config.Read()
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println(updatedConfig.CurrentUserName)
	fmt.Println(updatedConfig.DbUrl)

}
