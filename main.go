package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/silveiralexandre/logstash-exporter/tower"
)

const (
	releaseNumber = "v0.0.2-rc"
	helpMsg       = (`Generic exporter to feed JSON data from input of a source API into a target Logstash instance

Usage:
	logstash-exporter -n [ <report_name> ] [ -t <timeout> ] [ -l <limit> ] [ -r <retries> ] [ -ls ]

Examples:
	$ logstash-exporter -n tower-inventorydata -l 50 -t 90 -r 3
	$ logstash-exporter -ls

This will execute the informed script/command 50 times, limiting concurrency to 10 at the time.
Options:`)
)

var (
	reportName = flag.String("n", "", "Name of report to be executed")
	reportList = flag.Bool("ls", false, "List of currently supported reports")
	timeout    = flag.Int("t", 90, "Timeout in seconds for requests (default is 90)")
	retries    = flag.Int("r", 3, "Number of retries for HTTP requests (default is 3)")
	limit      = flag.Int("l", 50, "Limit of concurrent requests to be executed (default is 50)")
	version    = flag.Bool("v", false, "Prints version information")

	reports = map[string]string{
		"tower-inventorydata": "Ansible Inventory Variable Data presents metadata as predefined by users",
	}
)

func main() {
	flag.Usage = func() {
		flagSet := flag.CommandLine
		fmt.Println(helpMsg)
		order := []string{"n", "ls", "l", "t", "r", "v"}
		for _, name := range order {
			flag := flagSet.Lookup(name)
			fmt.Printf("   -%s\t", flag.Name)
			fmt.Printf("     %s\n", flag.Usage)
		}
	}
	flag.Parse()

	if *version {
		printReleaseNumber()
	}
	if *reportList {
		listExistingReports()
	}

	run(reportName)
}

func run(reportName *string) {
	switch *reportName {
	case "tower-inventorydata":
		tower.ReportInventoryMetadata(limit, timeout, retries)
	default:
		flag.Usage()
	}
}

func listExistingReports() {
	for name, description := range reports {
		fmt.Printf("%v\t%v\n", name, description)
	}
	os.Exit(0)
}

func printReleaseNumber() {
	fmt.Printf("logstash-exporter:%v\n", releaseNumber)
	os.Exit(0)
}
