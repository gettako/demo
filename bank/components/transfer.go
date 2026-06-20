// transfer.go
// Description: Transfer form with native bubbletea integration and RPC validation.

package components

import (
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"time"

	"charm.land/bubbles/v2/textinput"
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"

	"banking/utils"

	"gettako.dev/tako/contracts"
)

type Transfer struct {
	app        contracts.Application
	destInput  textinput.Model
	amtInput   textinput.Model
	focusIndex int // 0 = dest, 1 = amt
}

func NewTransfer(app contracts.Application) *Transfer {
	return &Transfer{app: app}
}

func (c *Transfer) ID() string { return "transfer" }

func (c *Transfer) Init(ctx contracts.Context) {
	c.destInput = textinput.New()
	c.destInput.Placeholder = "1234567890"

	c.amtInput = textinput.New()
	c.amtInput.Placeholder = "50000"
	c.focusIndex = -1
}

//nolint:funlen
func (c *Transfer) Keys(km contracts.KeyManager) {
	km.Bind("tab", func() {
		switch c.focusIndex {
		case -1:
			c.focusIndex = 0
			c.destInput.Focus()
		case 0:
			c.focusIndex = 1
			c.destInput.Blur()
			c.amtInput.Focus()
		default:
			c.focusIndex = -1
			c.amtInput.Blur()
		}
	})

	km.Bind("down", func() {
		switch c.focusIndex {
		case -1:
			c.focusIndex = 0
			c.destInput.Focus()
		case 0:
			c.focusIndex = 1
			c.destInput.Blur()
			c.amtInput.Focus()
		}
	})

	km.Bind("up", func() {
		switch c.focusIndex {
		case 1:
			c.focusIndex = 0
			c.amtInput.Blur()
			c.destInput.Focus()
		case 0:
			c.focusIndex = -1
			c.destInput.Blur()
		}
	})

	km.Bind("enter", func() {
		// Validate and submit
		dest := c.destInput.Value()
		amtStr := c.amtInput.Value()
		amt, err := strconv.ParseFloat(amtStr, 64)
		if dest == "" {
			c.app.UI().Dialog().Alert("Error", "Destination cannot be empty.")
			return
		}
		if err != nil || amt <= 0 {
			c.app.UI().Dialog().Alert("Error", "Amount must be a valid positive number.")
			return
		}

		lang := c.app.Lang()
		msg := lang.T("transfer.confirm", "amount", utils.FormatIDR(amt), "dest", dest)

		c.app.UI().Dialog().Confirm(msg, func(yes bool) {
			if yes {
				c.processTransfer(dest, amt)
			}
		})
	})
}

// Native Update handles Bubbletea messages directly for text input typing.
func (c *Transfer) UpdateNative(msg any) (any, bool) {
	var cmd tea.Cmd
	var cmds []tea.Cmd
	consumed := false

	if teaMsg, ok := msg.(tea.Msg); ok {
		switch c.focusIndex {
		case 0:
			c.destInput, cmd = c.destInput.Update(teaMsg)
			if cmd != nil {
				cmds = append(cmds, cmd)
			}
		case 1:
			c.amtInput, cmd = c.amtInput.Update(teaMsg)
			if cmd != nil {
				cmds = append(cmds, cmd)
			}
		}

		if keyMsg, isKey := msg.(tea.KeyMsg); isKey {
			if c.focusIndex == -1 {
				// Navigation mode: Do not consume any keys!
				consumed = false
			} else {
				// Form input mode: Only pass through hotkeys
				switch keyMsg.String() {
				case "tab", "shift+tab", "enter", "up", "down", "f1", "f2", "f3", "f4", "alt+1", "alt+2", "alt+3", "alt+4":
					consumed = false
				default:
					consumed = true
				}
			}
		}
	}

	c.app.UI().RequestRender()
	return tea.Batch(cmds...), consumed
}

func (c *Transfer) processTransfer(dest string, amount float64) {
	// 1. Call RPC Fraud Check
	resResp, err := c.app.RPC().Call("fraud.check").WithPayload(amount).WithContext(c.app.Context().StdCtx()).Await()
	res := resResp.Data
	if err != nil {
		c.app.Logger().Error("RPC Error", "err", err)
		return
	}

	resMap := res.(map[string]any)
	if safe, ok := resMap["safe"].(bool); ok && !safe {
		c.app.Logger().Warn("Transfer Blocked!", "dest", dest, "amount", amount)
		return // Blocked by fraud
	}

	// 2. Modify State Balance
	balFloat := c.app.State().Get("balance", 0.0).Float64()
	newBal := balFloat - amount
	c.app.State().Mutate("balance").Value(newBal).Broadcast()

	// 3. Save to Storage (BoltDB Persistence)
	txID := time.Now().Format("20060102150405")
	record := Transaction{
		ID:        txID,
		Dest:      dest,
		Amount:    int(amount),
		Status:    "SUCCESS",
		Timestamp: time.Now(),
	}

	// Save persistent history record using the Context().KV() interface
	var history []Transaction
	_ = c.app.KV().Get("history.all", &history)
	history = append(history, record)
	_ = c.app.KV().Set("history.all", history)

	// Keep receipt for legacy purposes just in case
	recJSON, _ := json.Marshal(record)
	_ = os.WriteFile(fmt.Sprintf("receipt_%s.txt", dest), recJSON, 0644)

	// 4. Emit Event
	c.app.Events().Dispatch("transaction:success", record)

	// Reset form
	c.destInput.SetValue("")
	c.amtInput.SetValue("")
	// RequestRender is called automatically by Broadcast() above.
}

func (c *Transfer) Render() string {
	width, _ := c.app.UI().Dimensions()
	lang := c.app.Lang()

	boxStyle := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(PrimaryColor).
		Padding(1, 4).
		MarginBottom(1).
		Width(width - 4)

	// Update text input widths dynamically
	c.destInput.SetWidth(width - 14)
	c.amtInput.SetWidth(width - 14)

	title := lipgloss.NewStyle().Foreground(PrimaryColor).Bold(true).Render(lang.T("transfer.title"))

	destLabel := lipgloss.NewStyle().Foreground(lipgloss.Color("#888888")).Render(lang.T("transfer.dest"))
	amtLabel := lipgloss.NewStyle().Foreground(lipgloss.Color("#888888")).Render(lang.T("transfer.amount"))

	destView := c.destInput.View()
	amtView := c.amtInput.View()

	form := lipgloss.JoinVertical(lipgloss.Center,
		title,
		"",
		destLabel,
		destView,
		"",
		amtLabel,
		amtView,
		"",
		lipgloss.NewStyle().Foreground(MutedColor).Render("[tab] next field | [enter] submit"),
	)

	// Add success/danger border styling
	box := boxStyle.Render(form)
	return lipgloss.Place(0, 0, lipgloss.Center, lipgloss.Center, box)
}
