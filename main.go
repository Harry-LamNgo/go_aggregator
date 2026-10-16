package main

import (
	"fmt"

	"github.com/Harry-LamNgo/go_aggregator/internal/config"
)

func main() {
	cfg, err := config.Read()
	if err != nil {
		fmt.Println("Error reading config: ", err)
		return
	}

	err = cfg.SetUser("harryngo")
	if err != nil {
		fmt.Println("Error setting current username: ", err)
		return
	}

	cfg, err = config.Read()
	if err != nil {
		fmt.Println("Error reading config: ", err)
		return
	}

	fmt.Println(cfg)
}
