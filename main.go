package main

import (
	"flag"
	"fmt"
	"os"
	"time"

	"github.com/pkg/errors"
)

var bqCli = &BigQueryImpl{}

func main() {
	targetYmd := flag.String("targetYmd", "", "target date [yyyyMMdd]")
	flag.Parse()

	if *targetYmd == "" {
		fmt.Println("targetYmd should be set")
		os.Exit(1)
	}
	err := checkTargetYmd(*targetYmd)
	if err != nil {
		fmt.Printf("%v", err)
		os.Exit(1)
	}

	q, err := bqCli.CreateQuery(*targetYmd)
	if err != nil {
		fmt.Println("failed to create buffer of query")
		os.Exit(1)
	}

	res, err := bqCli.SelectTable(q)
	if err != nil {
		fmt.Println("failed to select table")
		os.Exit(1)
	}

	err = sendToSlack(*targetYmd, res)
	if err != nil {
		fmt.Printf("failed to send to slack: %v", err)
		os.Exit(1)
	}
}

func checkTargetYmd(targetYmd string) error {
	_, err := time.Parse("20060102", targetYmd)
	if err != nil {
		err := errors.New("targetYmd is bad format")
		return err
	}

	fmt.Printf("Start Reporting: %s\n", targetYmd)

	return nil
}
