// Package main provides the entry point for the Tako mini-application.
package main

import (
	"fmt"
	"log"

	"github.com/gettako/tako"
	"github.com/gettako/tako/contracts"
	"github.com/gettako/tako/demo/commands"
	"github.com/gettako/tako/demo/components"
	"github.com/gettako/tako/demo/layouts"
	"github.com/gettako/tako/demo/services"
	"github.com/gettako/tako/pkg/adapter/bubbletea"

	_ "github.com/gettako/tako/demo/plugins/fzf"
	_ "github.com/gettako/tako/demo/plugins/help"
	_ "github.com/gettako/tako/demo/plugins/status"
)

func main() {
	app := tako.NewApp()

	app.CliRegistry().Register(&commands.FooCommand{})

	// 0. Register services to the IoC Container (Singleton, Lazy, Transient)
	app.Container().Singleton(
		new(contracts.Config),
		services.NewMapConfig(map[string]any{"app.debug": true}),
	)

	app.Container().Lazy(
		new(contracts.KVStore),
		func() (any, error) {
			return services.NewInMemoryStorage(), nil
		},
	)

	app.Container().Transient(
		new(*services.APIClient),
		func() (any, error) {
			return services.NewAPIClient(), nil
		},
	)

	// 1. Push the base layer — auto-focus now fires automatically.
	//    The Router will focus zone "base" at level 0 without an explicit Focus() call.
	app.Stack().Push("base")

	// ── High-Level Overlay API ────────────────────────────────────────────────

	// 2. Dashboard toggle using the new OverlayManager.
	//    One call replaces the old three-step: Stack.Push + Router.Focus + Hooks.Set.
	app.Keys().Bind([]string{"ctrl+d", "f1"}, func() {
		if app.Overlay().Top() == "dashboard" {
			app.Overlay().Close()
		} else {
			app.Overlay().Show("dashboard", func() any {
				return " 📊 DASHBOARD OVERLAY \n\n" +
					" This layer is managed by app.Overlay() \n" +
					" NOT by manual Stack/Hook/Focus calls. \n\n" +
					" Press [ctrl+d] to toggle. "
			})
		}
	})

	// 3. Global ESC to close any active overlay.
	//    OverlayManager.Close() is a no-op when at the base layer — safe to call unconditionally.
	app.Keys().Bind("esc", func() {
		if app.Overlay().IsActive() {
			top := app.Overlay().Top()
			app.Overlay().Close()
			app.Emit(top+":closed", "Closed via global esc")
		}
	})

	// 4. Base hotkeys: Logger & Transient Service demo.
	app.Keys().Zone("base").Bind("h", func() {
		var logger contracts.Logger
		if err := app.Make(&logger); err == nil && logger != nil {
			logger.Info("[Base Hotkey] 'h' pressed! Press [ctrl+f] to open search, [ctrl+c] to exit.")
		}
	})

	app.Keys().Zone("base").Bind("t", func() {
		var client *services.APIClient
		if err := app.Make(&client); err == nil && client != nil {
			var logger contracts.Logger
			if err := app.Make(&logger); err == nil && logger != nil {
				logger.Info("[Transient Demo] Generated API Client: " + client.ID)
			}
		}
	})

	// 5. Component and Dialog demonstrations.
	app.Keys().Zone("main").Bind("c", func() {
		// Show custom component (manages its own state and keys)
		app.Overlay().ShowComponent(components.NewCounterComponent())
	})

	app.Keys().Zone("base").Bind("d", func() {
		app.Overlay().Dialog().Confirm("Demo: press [y/enter] to confirm or [n/esc] to cancel", func(yes bool) {
			if yes {
				app.Emit("demo:confirmed", nil)
			} else {
				app.Emit("demo:cancelled", nil)
			}
		})
	})

	// 6. Demonstrate Lazy Storage: increment launch count.
	var store contracts.KVStore
	if err := app.Make(&store); err == nil && store != nil {
		var count int
		_ = store.Get("launch_count", &count) // Ignore error if not found
		count++
		_ = store.Set("launch_count", count)

		var logger contracts.Logger
		if err := app.Make(&logger); err == nil && logger != nil {
			// Log launch count to prove storage works across hot-reloads/invocations
			logger.Info(fmt.Sprintf("App launch count this session: %d", count))
		}
	}

	// 7. Register custom UI Renderer (Adapter as bridge, Layout in MiniApp).
	app.Container().Singleton(
		new(contracts.UIRenderer),
		bubbletea.NewAdapter(app.Context(), &layouts.AppLayout{}),
	)

	// 8. Mouse Router: register clickable zones for the demo.
	//    Hitboxes are updated on every render in layouts/mouse_demo.go.
	app.Mouse().Zone("left-panel").
		OnClick(func(x, y int) {
			itemRow := y - 4 // header(1)+gap(1)+border(1)+panel-header(1)
			if itemRow >= 0 && itemRow < len(layouts.MouseState.Items) {
				layouts.MouseState.Selected = itemRow
			}
			layouts.MouseState.LastEvent = fmt.Sprintf("click left-panel (%d,%d)", x, y)
		}).
		OnRightClick(func(x, y int) {
			if len(layouts.MouseState.Items) > 0 {
				layouts.MouseState.Selected = (layouts.MouseState.Selected + 1) % len(layouts.MouseState.Items)
			}
			layouts.MouseState.LastEvent = fmt.Sprintf("right-click left-panel (%d,%d)", x, y)
		}).
		OnScrollUp(func(x, y int) {
			if layouts.MouseState.Selected > 0 {
				layouts.MouseState.Selected--
			}
			layouts.MouseState.LastEvent = fmt.Sprintf("scroll-up left-panel (%d,%d)", x, y)
		}).
		OnScrollDown(func(x, y int) {
			if layouts.MouseState.Selected < len(layouts.MouseState.Items)-1 {
				layouts.MouseState.Selected++
			}
			layouts.MouseState.LastEvent = fmt.Sprintf("scroll-down left-panel (%d,%d)", x, y)
		}).
		OnDragStart(func(x, y int) {
			layouts.MouseState.LastEvent = fmt.Sprintf("drag-start left-panel (%d,%d)", x, y)
		})

	app.Mouse().Zone("right-panel").
		OnClick(func(x, y int) {
			layouts.MouseState.LastEvent = fmt.Sprintf("click right-panel (%d,%d)", x, y)
		}).
		OnDrop(func(fromZone string, x, y int) {
			if fromZone == "left-panel" && len(layouts.MouseState.Items) > 0 && layouts.MouseState.Selected < len(layouts.MouseState.Items) {
				moved := layouts.MouseState.Items[layouts.MouseState.Selected]
				layouts.MouseState.Items = append(layouts.MouseState.Items[:layouts.MouseState.Selected], layouts.MouseState.Items[layouts.MouseState.Selected+1:]...)
				layouts.MouseState.Dropped = append(layouts.MouseState.Dropped, moved)
				if layouts.MouseState.Selected >= len(layouts.MouseState.Items) && layouts.MouseState.Selected > 0 {
					layouts.MouseState.Selected--
				}
				layouts.MouseState.LastEvent = fmt.Sprintf("dropped '%s' from %s", moved, fromZone)
			}
		})

	// 7. Boot and Run the TUI application.
	if err := tako.Run(app); err != nil {
		log.Fatalf("Mini-app crashed: %v", err)
	}
}
