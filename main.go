package main

import (
	"blog-aggregator/internal/config"
	"fmt"
)

func main() {
	cfg, err := config.Read()
	if err != nil {
		fmt.Println(err)
	}

	cfg.SetUser("adam")

	cfg, err = config.Read()
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println(cfg.CurrentUser)
	fmt.Println(cfg.DbUrl)
}
