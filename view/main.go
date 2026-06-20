package main

import (
	"fmt"
	"os"

	"viewdemo/components"
	"viewdemo/providers"

	"gettako.dev/tako"
	btadapter "gettako.dev/tako/pkg/adapter/ui/bubbletea"
)

func main() {
	app := tako.NewApp()

	app.WithUIAdapter(btadapter.New(app)).
		WithProviders(
			&providers.ViewServiceProvider{},
		)

	app.Router().Route("/", components.NewApp(app)).Name("index")
	app.Router().GoTo("index")

	if err := tako.Run(app); err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}
}
