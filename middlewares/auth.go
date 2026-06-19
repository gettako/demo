// auth.go
// Description: Router middleware to protect the transfer page with a PIN.

package middlewares

import (
	"gettako.dev/tako/contracts"
)

// Auth is a router middleware that protects the /transfer route.
func Auth(ctx contracts.Context, to contracts.RouteInfo) bool {
	if to.Path != "/pin" {
		pinEntered := ctx.State().Get("pin_entered")
		if !pinEntered.Bool() {
			// Save intended destination, redirect to pin
			ctx.State().Mutate("intended_url").Value(to.Name).Broadcast()
			ctx.Router().GoTo("pin")
			return false // block the current navigation
		}
	}

	// Proceed to the requested route
	return true
}
