package main

import (
	"fmt"

	"github.com/fatih/color"
)

func printStartMessage() {
	color.Blue("--- Requests Source ---")

	fmt.Printf("%-20s\t: %s ~ %s \n",
		color.BlueString("REQUEST TIME"),
		oldest.StringOriginTime(),
		newest.StringOriginTime())

	fmt.Printf("%-20s\t: %-10d reqs\n",
		color.BlueString("REQUEST COUNT"),
		len(readreqs))

	fmt.Printf("%-20s\t: %-10d reqs\n",
		color.BlueString("IGNORED COUNT"),
		ignoredLine)

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
		fmt.Printf("%s %s\n", color.GreenString(" %d:", i+1), r.URL)
	}
	fmt.Println("...and more")
	fmt.Println()
}
