package main

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/gettako/tako/contracts"
	"github.com/gettako/tako/internal/console/cli"
	"github.com/gettako/tako/internal/event"
	"github.com/gettako/tako/internal/recorder"
	"github.com/gettako/tako/internal/state"
	"github.com/gettako/tako/pkg/foundation/commands/replayplayer"
)

// PaymentProcessor adalah contoh service kompleks yang memiliki "Bug" tersembunyi.
type PaymentProcessor struct {
	store contracts.StateManager
	bus   contracts.EventBus
}

func (p *PaymentProcessor) Start(ctx context.Context) {
	// Mengamati antrean event pembelian
	p.bus.Subscribe(ctx, "transaction.purchase", func(evt contracts.Event) {
		payload, ok := evt.Data.(map[string]any)
		if !ok {
			return
		}

		// Parsing payload (bisa int saat live, float64 saat di-replay dari JSON)
		var amount int
		switch v := payload["amount"].(type) {
		case int:
			amount = v
		case float64:
			amount = int(v)
		}
		item := payload["item"].(string)

		// SIMULASI BUG KRITIKAL:
		// Developer tidak memvalidasi angka negatif. Jika negatif, memicu Panic.
		if amount < 0 {
			crashMsg := fmt.Sprintf("FATAL BUG: Transaksi manipulatif terdeteksi! Amount: %d", amount)
			p.store.Key("system.status").Value("CRASHED").Broadcast()
			panic(crashMsg) // Aplikasi meledak!
		}

		// Memproses transaksi normal
		currentBalanceRaw := p.store.Get("user.balance")
		var currentBalance int
		if currentBalanceRaw != nil {
			switch v := currentBalanceRaw.(type) {
			case int:
				currentBalance = v
			case float64:
				currentBalance = int(v)
			}
		}

		newBalance := currentBalance - amount
		p.store.Key("user.balance").Value(newBalance).Broadcast()

		fmt.Printf("✅ Pembelian '%s' sukses. Terpotong: $%d (Sisa Saldo: $%d)\n", item, amount, newBalance)
	})
}

func main() {
	bus := event.NewBus()
	store := state.NewManager(bus)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// ---------------------------------------------------------
	// MODE REPLAY
	// ---------------------------------------------------------
	if len(os.Args) > 1 && os.Args[1] == "replay" {
		cmd := replayplayer.NewCommand(bus)

		// Mendaftarkan Observer untuk melacak jejak state yang terekam
		// CATATAN: Kita tidak menjalankan PaymentProcessor di mode ini,
		// karena kita murni mereplay State yang terekam persis sebelum crash terjadi.
		store.Observe("user.balance").
			OnUpdate(func(old, new any) {
				fmt.Printf("[Time-Travel] Saldo ter-update: $%v\n", new)
			}).
			Subscribe(ctx)

		store.Observe("system.status").
			OnUpdate(func(old, new any) {
				fmt.Printf("[Time-Travel] Status Sistem: %v\n", new)
			}).
			Subscribe(ctx)

		fmt.Println("▶️  Memulai pemutaran ulang (Replay)...")
		err := cmd.Execute(&cli.Context{}, os.Args[2:])
		if err != nil {
			fmt.Printf("Replay Error: %v\n", err)
		}
		
		fmt.Println("\n==================================================")
		fmt.Println("🔍 ANALISIS TIME-TRAVEL:")
		fmt.Println("Berdasarkan log di atas, terlihat status sistem berubah menjadi CRASHED")
		fmt.Println("tepat setelah adanya event anomali. Anda bisa melihat history persisnya!")
		fmt.Println("==================================================")
		return
	}

	// Inisialisasi Service (Hanya dijalankan di mode normal/live)
	processor := &PaymentProcessor{store: store, bus: bus}
	processor.Start(ctx)

	// ---------------------------------------------------------
	// MODE RECORDING (Skenario Asli)
	// ---------------------------------------------------------
	fmt.Println("=== TOKO ONLINE (RECORDING MODE) ===")

	os.Setenv("TAKO_RECORD_TRACE", "1")
	rec := recorder.New(bus).
		MaskFields("credit_card", "cvv")

	if err := rec.Start(ctx); err != nil {
		fmt.Printf("Gagal memulai recorder: %v\n", err)
		return
	}
	defer rec.Stop()

	// Set Saldo Awal
	store.Key("user.balance").Value(1000).Broadcast()
	store.Key("system.status").Value("HEALTHY").Broadcast()
	fmt.Println("💳 Saldo awal disetel ke $1000")

	// Simulasi ratusan event background dan UI
	fmt.Println("🖱️ Mensimulasikan 100+ aktivitas UI dan Background...")
	for i := 0; i < 50; i++ {
		time.Sleep(10 * time.Millisecond)
		bus.Publish("mouse.move", map[string]any{"x": 10 + i*2, "y": 200 - i})
		bus.Publish("system.health_check", map[string]any{"cpu_usage": 10 + i%15, "memory_mb": 200 + i*2})
	}
	for i := 0; i < 30; i++ {
		time.Sleep(10 * time.Millisecond)
		bus.Publish("ui.scroll", map[string]any{"offsetY": i * 15})
	}

	// Rentetan transaksi normal
	transaksiNormal := []map[string]any{
		{"item": "Buku Golang", "amount": 150, "credit_card": "1234-5678"},
		{"item": "Kopi Susu", "amount": 50, "credit_card": "1234-5678"},
		{"item": "Keyboard Mekanik", "amount": 300, "credit_card": "1234-5678"},
	}

	for i, tx := range transaksiNormal {
		time.Sleep(200 * time.Millisecond)
		bus.Publish("mouse.move", map[string]any{"x": 45 + i*10, "y": 100 + i*5})
		
		time.Sleep(100 * time.Millisecond)
		bus.Publish("ui.button_hover", map[string]any{"id": "btn_buy", "state": "active"})
		
		time.Sleep(150 * time.Millisecond)
		bus.Publish("ui.button_click", map[string]any{"id": "btn_buy"})
		
		time.Sleep(100 * time.Millisecond)
		bus.Publish("network.request", map[string]any{"endpoint": "/api/checkout", "method": "POST"})
		
		time.Sleep(100 * time.Millisecond)
		bus.Publish("transaction.purchase", tx)
		
		time.Sleep(100 * time.Millisecond)
		bus.Publish("network.response", map[string]any{"status": 200, "message": "OK"})
	}

	// Jeda sebelum bencana, simulasi event background
	time.Sleep(400 * time.Millisecond)
	bus.Publish("system.sync", map[string]any{"target": "cloud_db", "status": "pending"})
	time.Sleep(200 * time.Millisecond)
	bus.Publish("system.sync", map[string]any{"target": "cloud_db", "status": "success"})
	time.Sleep(400 * time.Millisecond)

	// SIMULASI SERANGAN HACKER (Memasukkan quantity negatif agar saldo bertambah)
	// Ini akan memicu panic di PaymentProcessor
	hackerTx := map[string]any{
		"item": "Hacker Exploit", "amount": -5000, "credit_card": "0000-0000",
	}
	
	fmt.Println("\n⚠️  Peringatan: Terdapat transaksi anomali masuk...")
	time.Sleep(500 * time.Millisecond)
	
	// Sengaja kita jangan tangkap panic di sini agar simulasi crash nyata!
	bus.Publish("transaction.purchase", hackerTx)
	
	// Code below won't be reached due to panic, but ideally this is what we would print:
	fmt.Println("\n✅ Rekaman Selesai! File trace tersimpan di folder .tako/traces/")
	fmt.Println("\nUntuk melakukan pemutaran ulang (Replay), jalankan perintah berikut:")
	fmt.Println("Otomatis (file terbaru) : go run main.go replay")
	fmt.Println("Spesifik                : go run main.go replay .tako/traces/<nama-file>.jsonl")
}
