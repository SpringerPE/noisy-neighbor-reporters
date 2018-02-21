package main

import (
	"github.com/SpringerPE/noisy-neighbor-reporters/pkg/apps/graphite-reporter/app"
)

func main() {
	cfg := app.LoadConfig()
	app.NewReporter(cfg).Run()
}
