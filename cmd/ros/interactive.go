package main

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"strconv"
	"strings"

	"ros/internal/domain/entity"
	"ros/internal/usecase"
)

type InteractiveCLI struct {
	app    *usecase.ROSUsecase
	reader *bufio.Reader
}

func NewInteractiveCLI(app *usecase.ROSUsecase) *InteractiveCLI {
	return &InteractiveCLI{
		app:    app,
		reader: bufio.NewReader(os.Stdin),
	}
}

func (cli *InteractiveCLI) Run(ctx context.Context) {
	for {
		fmt.Println("\n==============================================")
		fmt.Println("  ROS (Relation Optimizer System) Console")
		fmt.Println("==============================================")
		fmt.Println("1) Dashboard (List Nodes)")
		fmt.Println("2) Register Node (Add Person)")
		fmt.Println("3) Record Observation (Log Event)")
		fmt.Println("q) Quit")
		fmt.Print("\nSelect an option [1-3, q]: ")

		choice := cli.readLine()
		switch strings.ToLower(choice) {
		case "1":
			cli.showDashboard(ctx)
		case "2":
			cli.addPerson(ctx)
		case "3":
			cli.recordObservation(ctx)
		case "q", "exit":
			fmt.Println("[*] Exiting.")
			return
		default:
			fmt.Println("[!] Invalid selection.")
		}
	}
}

func (cli *InteractiveCLI) showDashboard(ctx context.Context) {
	list, err := cli.app.Dashboard(ctx)
	if err != nil {
		fmt.Printf("[!] Error: %v\n", err)
		return
	}
	if len(list) == 0 {
		fmt.Println("\n[i] No nodes registered yet. Please register a node first.")
		return
	}

	fmt.Println("\n--- Current Effective Odds ---")
	fmt.Printf("%-10s %-16s %-20s %-16s %-10s %-8s %-15s\n", "ID", "NAME", "TYPE", "STATUS", "SCORE(dB)", "EVENTS", "VERDICT")
	fmt.Println(strings.Repeat("-", 102))
	for _, item := range list {
		fmt.Printf("%-10s %-16s %-20s %-16s %+9.2f  %-8d %-15s\n",
			item.ID, item.Name, item.CurrentType, item.Status, item.CurrentLogOddsDB, item.TotalEvents, item.Verdict())
	}
}

func (cli *InteractiveCLI) addPerson(ctx context.Context) {
	fmt.Println("\n--- Register New Node ---")
	fmt.Print("Display Name: ")
	name := cli.readLine()
	if name == "" {
		fmt.Println("[!] Name cannot be empty.")
		return
	}

	fmt.Print("Node ID (leave empty for auto-generation): ")
	id := cli.readLine()
	if id == "" {
		id = generateShortID("p")
	}

	types := []struct {
		val  string
		desc string
	}{
		{string(entity.TypeAdultCustomer), "Adult / Customer (contractual / paying)"},
		{string(entity.TypeAdultNonCustomer), "Adult / Non-customer (social / acquaintance)"},
		{string(entity.TypeEmployee), "Employee / Internal resource / Partner"},
		{string(entity.TypeChildCustomer), "Child / Customer (minor in structured relation)"},
		{string(entity.TypeChildNonCustomer), "Child / Non-customer (minor baseline tolerance)"},
	}

	fmt.Println("\nSelect Persona Type:")
	for i, t := range types {
		fmt.Printf("  %d) %s\n", i+1, t.desc)
	}
	idx := cli.readChoice(len(types))
	chosenType := types[idx-1].val

	if err := cli.app.AddPerson(ctx, id, name, chosenType); err != nil {
		fmt.Printf("[!] Registration failed: %v\n", err)
		return
	}
	fmt.Printf("[+] Node registered: %s (ID: %s, Type: %s)\n", name, id, chosenType)
}

func (cli *InteractiveCLI) recordObservation(ctx context.Context) {
	list, err := cli.app.Dashboard(ctx)
	if err != nil || len(list) == 0 {
		fmt.Println("\n[!] No nodes available for logging. Please register a node first.")
		return
	}

	fmt.Println("\nSelect Target Node:")
	for i, item := range list {
		fmt.Printf("  %d) %s (%s) [%s / %.1f dB]\n", i+1, item.Name, item.ID, item.Status, item.CurrentLogOddsDB)
	}
	pIdx := cli.readChoice(len(list))
	targetPerson := list[pIdx-1]

	fmt.Println("\nSelect Polarity (Sign):")
	fmt.Println("  1) pos (+) Ideal alignment / discipline / integrity")
	fmt.Println("  2) neg (-) Noise / exploitation / breach")
	sIdx := cli.readChoice(2)
	chosenSign := "pos"
	if sIdx == 2 {
		chosenSign = "neg"
	}

	inevChoices := []struct {
		val  string
		desc string
	}{
		{"daily", "daily (Habitual baseline / Half-life: 30d)"},
		{"critical_moment", "critical_moment (Crucible test / Half-life: INFINITE)"},
		{"transient_noise", "transient_noise (Superficial noise / Half-life: 7d)"},
	}
	fmt.Println("\nSelect Inevitability Tier:")
	for i, in := range inevChoices {
		fmt.Printf("  %d) %s\n", i+1, in.desc)
	}
	iIdx := cli.readChoice(len(inevChoices))
	chosenInev := inevChoices[iIdx-1].val

	fmt.Print("\nObserved Episode (Fact-based memo): ")
	memo := cli.readLine()
	if memo == "" {
		fmt.Println("[!] Memo cannot be empty.")
		return
	}

	if err := cli.app.RecordObservation(ctx, targetPerson.ID, memo, chosenSign, chosenInev); err != nil {
		fmt.Printf("[!] Failed to record observation: %v\n", err)
		return
	}
	fmt.Printf("[+] Recorded event: [%s] -> %s / %s applied.\n", targetPerson.Name, chosenSign, chosenInev)
}

func (cli *InteractiveCLI) readLine() string {
	text, _ := cli.reader.ReadString('\n')
	return strings.TrimSpace(text)
}

func (cli *InteractiveCLI) readChoice(max int) int {
	for {
		fmt.Printf("Choice [1-%d]: ", max)
		input := cli.readLine()
		num, err := strconv.Atoi(input)
		if err == nil && num >= 1 && num <= max {
			return num
		}
		fmt.Printf("Invalid input. Please enter a number between 1 and %d.\n", max)
	}
}
