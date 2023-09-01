package main

import (
	"fmt"

	"github.com/AlexandrLitkevich/qwery/internal/config"
)

func main() {
	cfg := config.MustLoad()
	fmt.Println("run server", cfg)
}