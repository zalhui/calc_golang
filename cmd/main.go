package main

import (
	"github.com/zalhui/calc_golang/internal/application"
)

func main() {
	app := application.New()
	app.RunServer()
}
