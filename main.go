package main

import (
	"fmt"

	"github.com/SaschaRunge/gator/internal/config"
)

func main() {
	cfg, _ := config.Read()
	cfg.SetUser("Sascha")
	cfg, _ = config.Read()

	fmt.Printf("%+v", cfg)
}
