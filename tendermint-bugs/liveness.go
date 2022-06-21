package tendermintbugs

import (
	"time"

	"github.com/netrixframework/netrix/testlib"
	"github.com/netrixframework/netrix/types"
	"github.com/netrixframework/tendermint-testing/common"
)

func SetupLivenessTimer(timeout time.Duration) common.SetupOption {
	return func(ctx *testlib.Context) {
		go func() {
			ctx.Logger().Info("Waiting for timeout to expire")
			time.Sleep(timeout)
			ctx.Logger().Info("Healing network, checking liveness")
			ctx.Vars.Set(testFinishedKey, true)
		}()
	}
}

func IsBeforeLivenessCheck(e *types.Event, ctx *testlib.Context) bool {
	finished, ok := ctx.Vars.GetBool(testFinishedKey)
	if !ok {
		return true
	}
	return !finished
}

const testFinishedKey = "test_finished"
