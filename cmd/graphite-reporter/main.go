package main

import "github.com/SpringerPE/noisy-neighbor-reporters/cmd/graphite-reporter/app"

func main() {
	cfg := app.LoadConfig()
	app.NewReporter(cfg).Run()
}
