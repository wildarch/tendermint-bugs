package main

import (
	"flag"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	tendermintbugs "tendermint-bugs/tendermint-bugs"

	"github.com/netrixframework/netrix/config"
	"github.com/netrixframework/netrix/testlib"
	"github.com/netrixframework/tendermint-testing/common"
	"github.com/netrixframework/tendermint-testing/util"
)

var bug = flag.String("bug", "", "Testcase to run (dummy, bug001, bug002)")

func main() {
	flag.Parse()
	sysParams := common.NewSystemParams(4)
	var testcase *testlib.TestCase
	switch *bug {
	case "dummy":
		testcase = tendermintbugs.Dummy(sysParams)
	case "bug001":
		testcase = tendermintbugs.Bug001(sysParams)
	case "bug002":
		testcase = tendermintbugs.Bug002(sysParams)
	default:
		fmt.Println("-bug is a required flag")
		flag.Usage()
		os.Exit(1)
	}

	// Catch SIGTERM, and stop the server if we receive it
	termCh := make(chan os.Signal, 1)
	signal.Notify(termCh, os.Interrupt, syscall.SIGTERM)

	server, err := testlib.NewTestingServer(
		&config.Config{
			APIServerAddr: "192.167.0.1:7074",
			NumReplicas:   sysParams.N,
			LogConfig: config.LogConfig{
				Format: "json",
				Path:   "/tmp/tendermint/log/checker.log",
			},
		},
		&util.TMessageParser{},
		[]*testlib.TestCase{
			testcase,
		},
	)

	if err != nil {
		fmt.Printf("Failed to start server: %s\n", err.Error())
		os.Exit(1)
	}

	// Stop the server if we receive SIGTERM
	go func() {
		<-termCh
		server.Stop()
	}()

	server.Start()
}
