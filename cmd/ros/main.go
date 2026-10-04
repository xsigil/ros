package main

import (
	"context"
	"crypto/rand"
	"flag"
	"fmt"
	"math/big"
	"os"
	"strings"

	"ros/internal/domain/entity"
	"ros/internal/infrastructure/sqlite3"
	"ros/internal/usecase"
)

func main() {
	ctx := context.Background()
	db, err := sqlite3.NewDB(ctx, "ros.db")
	if err != nil {
		fmt.Fprintf(os.Stderr, "DB init error: %v\n", err)
		os.Exit(1)
	}
	defer db.Close()

	txManager := sqlite3.NewTxManager(db)
	rosRepo := sqlite3.NewROSRepository(db)
	app := usecase.NewROSUsecase(rosRepo, txManager)

	// 引数なし、または 'i' / 'interactive' の場合は対話モードへ
	if len(os.Args) < 2 || os.Args[1] == "i" || os.Args[1] == "interactive" {
		interactiveCLI := NewInteractiveCLI(app)
		interactiveCLI.Run(ctx)
		return
	}

	switch os.Args[1] {
	case "add":
		addCmd := flag.NewFlagSet("add", flag.ExitOnError)
		id := addCmd.String("id", "", "Unique node ID (optional: auto-generated if omitted)")
		name := addCmd.String("name", "", "Display Name (required)")
		pType := addCmd.String("type", string(entity.TypeAdultCustomer), "Persona type")
		addCmd.Parse(os.Args[2:])

		if *name == "" {
			printAddUsage()
			return
		}

		// -id が省略された場合は短いランダムID (例: p_a1b2) を自動生成
		nodeID := *id
		if nodeID == "" {
			nodeID = generateShortID("p")
		}

		if err := app.AddPerson(ctx, nodeID, *name, *pType); err != nil {
			fmt.Fprintf(os.Stderr, "\n[!] Add Failed: %v\n", err)
			printPersonaTypes()
			return
		}
		fmt.Printf("[+] Registered: %s (ID: %s, Type: %s)\n", *name, nodeID, *pType)

	case "log":
		logCmd := flag.NewFlagSet("log", flag.ExitOnError)
		id := logCmd.String("id", "", "Target Person ID (required)")
		memo := logCmd.String("memo", "", "Observed Episode (required)")
		sign := logCmd.String("sign", "", "pos | neg (required)")
		inev := logCmd.String("inev", "", "daily | critical_moment | transient_noise (required)")
		logCmd.Parse(os.Args[2:])

		if *id == "" || *memo == "" || *sign == "" || *inev == "" {
			printLogUsage()
			return
		}

		if err := app.RecordObservation(ctx, *id, *memo, *sign, *inev); err != nil {
			fmt.Fprintf(os.Stderr, "\n[!] Log Failed: %v\n", err)
			printLogChoices()
			return
		}
		fmt.Printf("[+] Recorded event for [%s] -> %s / %s (\"%s\")\n", *id, *sign, *inev, *memo)

	case "ls":
		list, err := app.Dashboard(ctx)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			return
		}
		if len(list) == 0 {
			fmt.Println("No nodes registered yet. Run 'ros add' first.")
			return
		}
		fmt.Printf("%-10s %-16s %-20s %-16s %-10s %-8s %-15s\n", "ID", "NAME", "TYPE", "STATUS", "SCORE(dB)", "EVENTS", "VERDICT")
		fmt.Println(strings.Repeat("-", 102))
		for _, item := range list {
			fmt.Printf("%-10s %-16s %-20s %-16s %+9.2f  %-8d %-15s\n",
				item.ID, item.Name, item.CurrentType, item.Status, item.CurrentLogOddsDB, item.TotalEvents, item.Verdict())
		}

	case "help", "-h", "--help":
		printUsage()

	default:
		fmt.Fprintf(os.Stderr, "Unknown command: %s\n\n", os.Args[1])
		printUsage()
	}
}

// -------------------------------------------------------------
// ヘルプ & 選択肢の表示
// -------------------------------------------------------------

func printUsage() {
	fmt.Println(`ROS (Relation Optimizer System) CLI
A low-pass Bayesian signal detector for personal bandwidth defense.

Commands:
  add    Register a new person/node
  log    Record an observed behavioral event
  ls     Display dashboard ranked by effective log-odds (dB)

Run 'ros <command> -h' or see details below.`)
	fmt.Println()
	printAddUsage()
	fmt.Println()
	printLogUsage()
}

func printAddUsage() {
	fmt.Println(`Usage: ros add -name <name> [-id <id>] [-type <type>]
-name  Display name of the person (required)
  -id    Custom identifier (optional, auto-generated if omitted)
  -type  Persona classification (optional, default: adult_customer)`)
	printPersonaTypes()
}

func printPersonaTypes() {
	fmt.Println(`
Available -type choices:
  * adult_customer     : Client, paying partner, contractual adult
  * adult_non_customer : Acquaintance, stranger, social connection
  * employee           : Team member, internal resource, partner
  * child_customer     : Minor with transactional/structured relation
  * child_non_customer : Minor (broad tolerance / zero baseline penalty)`)
}

func printLogUsage() {
	fmt.Println(`Usage: ros log -id <id> -memo <episode> -sign <sign> -inev <inevitability>
-id    Target person ID
  -memo  Fact-based note of the observation
  -sign  Observed polarity
  -inev  Inevitability layer`)
	printLogChoices()
}

func printLogChoices() {
	fmt.Println(`
Available -sign choices:
  * pos : Positive alignment with ideal protocol (H1)
  * neg : Negative noise / exploit / breach

Available -inev choices:
  * daily           : Habitual protocol, punctuality, responsiveness (Half-life: 30d)
  * critical_moment : Crucible test, liability acceptance, betrayal (Half-life: INFINITE)
  * transient_noise : Flattery, lavish treating, one-off friction (Half-life: 7d)`)
}

func generateShortID(prefix string) string {
	const chars = "abcdefghijklmnopqrstuvwxyz0123456789"
	b := make([]byte, 4)
	for i := range b {
		num, _ := rand.Int(rand.Reader, big.NewInt(int64(len(chars))))
		b[i] = chars[num.Int64()]
	}
	return fmt.Sprintf("%s_%s", prefix, string(b))
}
