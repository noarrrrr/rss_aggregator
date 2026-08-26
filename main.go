package main

import (
	"fmt"

	"github.com/noarrrrr/rss_aggregator/internal/config"
)

func main() {
	cfg, err := config.Read()
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(cfg)
	cfg.Current_username = "Noar"
	err = cfg.SetUser()
	if err != nil {
		fmt.Println(err)
		return
	}
	cfg, err = config.Read()
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(cfg)
}
