package main

import (
	"fmt"
	"os"

	"github.com/fatih/color"
)

func printStartMessage() {
	color.Blue("--- Requests Source ---")

	fmt.Printf("%-20s\t: %s ~ %s \n",
		color.BlueString("REQUEST TIME"),
		oldest.StringOriginTime(),
		newest.StringOriginTime())

	fmt.Printf("%-20s\t: %-10d reqs\n",
		color.BlueString("REQUESTS"),
		len(readreqs))

	fmt.Printf("%-20s\t: %-10d reqs\n",
		color.BlueString("IGNORED"),
		ignoredLine)

	fmt.Printf("%-20s\t: %-10d reqs\n",
		color.BlueString("NON SUPPORTED"),
		nonSuportedLine)

	fmt.Printf("%-20s\t: %-10d reqs\n",
		color.BlueString("PARSE ERROR"),
		parseErrLine)

	fmt.Printf("%-20s\t: %s\n",
		color.BlueString("DryRun"),
		fmt.Sprint(isDryrun))
	fmt.Println()

	color.Green("--- Repeat Plan ---")
	fmt.Printf("%-20s\t: %s ~ %s \n",
		color.GreenString("REPEAT TIME"),
		oldest.StringPlanTime(),
		newest.StringPlanTime())
	fmt.Printf("%-20s\t: %s:%s\n",
		color.GreenString("REPEAT TARGET"),
		host,
		port)
	color.Green("Repeat Samples\t:")
	for i, r := range readreqs[0:5] {
		fmt.Printf("%s %s\n", color.GreenString(" %d:", i+1), r.String())
	}
	fmt.Println("...and more")
	fmt.Println()
}

func waitPrompt() {
	var key string
	var ok bool
	for !ok {
		fmt.Print(color.MagentaString("Start/Cancel>"))
		fmt.Scanf("%s", &key)
		switch key {
		case "S", "s", "Start", "start":
			ok = true
		case "C", "c", "Cancel", "cancel":
			fmt.Println("canceled.")
			os.Exit(1)
		default:
			continue
		}
	}
}
