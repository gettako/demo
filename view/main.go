package main

import (
	"fmt"
	"os"

	"viewdemo/components"
	"viewdemo/providers"

	"gettako.dev/tako"
	"gettako.dev/tako/contracts"
	btadapter "gettako.dev/tako/pkg/adapter/ui/bubbletea"
)

func main() {
	app := tako.NewApp()

	app.WithUIAdapter(btadapter.New(app)).
		WithProviders(
			&providers.ViewServiceProvider{},
		)

	// Mount the root Vue/Blade-like component
	app.Router().Root(components.NewApp(app), func(r contracts.Router) {
		// Create a dummy route just to satisfy the router requirement for now
		// In the future, App could be a standalone route instead of Root
		r.Route("/", components.NewApp(app)).Name("index")
	})

	app.Router().GoTo("index")

	if err := tako.Run(app); err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}
}
