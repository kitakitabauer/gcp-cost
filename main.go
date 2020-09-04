package main

import (
	"flag"
	"fmt"
	"os"
	"time"

	"github.com/pkg/errors"
	"go.uber.org/zap"
)

func main() {
	targetYmd := flag.String("targetYmd", "", "target date [yyyyMMdd]")
	channel := flag.String("channel", "", "slack channel")
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

	if *channel == "" {
		fmt.Println("channel should be set")
		os.Exit(1)
	}

	logger, _ := zap.NewDevelopment()
	bqCli := &BigQueryImpl{log: logger}

	q, err := bqCli.CreateQuery(*targetYmd)
	if err != nil {
		logger.Error("failed to create buffer of query", zap.Error(err))
		os.Exit(1)
	}

	res, err := bqCli.SelectTable(q)
	if err != nil {
		os.Exit(1)
	}

	err = sendToSlack(*targetYmd, *channel, res)
	if err != nil {
		logger.Error("failed to send to slack", zap.Error(err))
		os.Exit(1)
	}
}

func checkTargetYmd(targetYmd string) error {
	_, err := time.Parse("2006/01/02", targetYmd)
	if err != nil {
		err := errors.New("targetYmd is bad format")
		return err
	}

	fmt.Printf("Start Reporting: %s\n", targetYmd)

	return nil
}
