// main.go
// Description: Entry point for the TakoBank Terminal application.

package main

import (
	"fmt"
	"os"

	"banking/components"
	"banking/console"
	"banking/middlewares"
	"banking/plugins"
	"banking/services"

	"gettako.dev/tako"
	"gettako.dev/tako/contracts"
	btadapter "gettako.dev/tako/pkg/adapter/ui/bubbletea"
)

func main() {
	// Initialize core framework
	app := tako.NewApp()

	// Fluent Configuration
	app.WithMiddleware(middlewares.Auth).
		WithUIAdapter(btadapter.New(app)).
		WithProviders(
			services.NewExchangeProvider(),
			services.NewFraudProvider(),
			plugins.NewCRMProvider(),
		)

	// Register Standalone Auth Route
	app.Router().Route("/pin", components.NewPin(app)).Name("pin")

	// Register Demo Route Components inside the Main Layout Root
	app.Router().Root(components.NewLayout(app), func(r contracts.Router) {
		r.Route("/dashboard", components.NewDashboard(app)).Name("dashboard")
		r.Route("/transfer", components.NewTransfer(app)).Name("transfer")
		r.Route("/history", components.NewHistory(app)).Name("history")
		r.Route("/profile", components.NewProfile(app)).Name("profile")
	})

	// Set Default Route
	app.Router().GoTo("dashboard")

	// Register Custom CLI Commands
	app.CliRegistry().Register(console.NewExportCommand(app))

	// Boot up the framework!
	if err := tako.Run(app); err != nil {
		fmt.Printf("Application Error: %v\n", err)
		os.Exit(1)
	}
}
