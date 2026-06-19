// export_command.go
// Description: Headless CLI command for exporting history.

package console

import (
	"fmt"

	"gettako.dev/tako/contracts"
)

type ExportCommand struct {
	app contracts.Application
}

func NewExportCommand(app contracts.Application) *ExportCommand {
	return &ExportCommand{app: app}
}

func (c *ExportCommand) Signature() string {
	return "export:history {--month= : The month to export}"
}

func (c *ExportCommand) Description() string {
	return "Export transaction history to a CSV file"
}

func (c *ExportCommand) Handle(ctx contracts.CommandContext) error {
	month := ctx.OptionValue("month")
	if month == "" {
		month = "all"
	}

	fmt.Printf("Exporting transaction history for month: %s\n", month)
	fmt.Println("Connecting to Storage Manager...")

	storage := c.app.Storage()

	// Simulate reading history and generating CSV
	csvContent := "ID,Dest,Amount,Status\n1,Budi,50000,OK\n2,Andi,120000,OK\n3,Susi,15000000,BLOCKED\n"
	filename := fmt.Sprintf("export_%s.csv", month)

	err := storage.Set(filename, []byte(csvContent))
	if err != nil {
		return fmt.Errorf("failed to write export: %w", err)
	}

	fmt.Printf("✅ Export completed! Saved to .tako/storage/%s\n", filename)
	return nil
}
