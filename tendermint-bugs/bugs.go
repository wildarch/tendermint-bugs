package tendermintbugs

import (
	"fmt"
	"strings"
	"time"

	"github.com/netrixframework/netrix/log"
	"github.com/netrixframework/netrix/testlib"
	"github.com/netrixframework/netrix/types"
	"github.com/netrixframework/tendermint-testing/common"
	"github.com/netrixframework/tendermint-testing/util"
)

type Partition = [][]string

// Does not inject any faults, provided to verify communication between test server and nodes
func Dummy(sysParams *common.SystemParams) *testlib.TestCase {
	sm := testlib.NewStateMachine()
	sm.Builder().On(common.HeightReached(5), testlib.SuccessStateLabel)
	filters := testlib.NewFilterSet()

	testCase := testlib.NewTestCase(
		"Dummy",
		30*time.Second,
		sm,
		filters,
	)
	testCase.SetupFunc(common.Setup(sysParams, labelNodes))
	return testCase
}

// Isolates two nodes in the first prevote step.
// This will prevent the nodes from consensus.
// After 1 minute the network is healed, and we expect them to recover (which in our tests they do not).
func Bug001(sysParams *common.SystemParams) *testlib.TestCase {
	sm := testlib.NewStateMachine()
	// Expect to be able to commit at least one block.
	// This may happen after the liveness timer expires and the network is healed.
	sm.Builder().On(common.HeightReached(2), testlib.SuccessStateLabel)

	filters := testlib.NewFilterSet()
	filters.AddFilter(trackNodeHeightRound)

	// In first prevote stage, create the following partitioning
	// {0,2} {1} {3}
	partition := Partition{
		{"node0", "node2"},
		{"node1"},
		{"node3"},
	}
	// Drop the prevotes at H=1/R=0 (the very first round)
	filters.AddFilter(testlib.If(
		common.IsMessageType(util.Prevote).
			And(IsMessageWithSenderHeightRound(1, 0).
				And(FromToIsolated(partition)).
				And(IsBeforeLivenessCheck)),
	).Then(testlib.DropMessage()))

	testCase := testlib.NewTestCase(
		"Bug001",
		2*time.Minute,
		sm,
		filters,
	)
	testCase.SetupFunc(common.Setup(sysParams, labelNodes, SetupLivenessTimer(time.Minute)))
	return testCase
}

func nodeLabel(idx int) string {
	return fmt.Sprintf("node%d", idx)
}

func labelNodes(c *testlib.Context) {
	parts := make([]*util.Part, len(c.Replicas.Iter()))
	for i, replica := range c.Replicas.Iter() {
		replicaSet := util.NewReplicaSet()
		replicaSet.Add(replica)
		parts[i] = &util.Part{
			ReplicaSet: replicaSet,
			Label:      nodeLabel(i),
		}
	}
	partition := util.NewPartition(parts...)
	c.Vars.Set("partition", partition)
	c.Logger().With(log.LogParams{
		"partition": partition.String(),
	}).Info("Partitioned replicas")
}

func getPartLabel(ctx *testlib.Context, id types.ReplicaID) string {
	partitionR, ok := ctx.Vars.Get("partition")
	if !ok {
		panic("No partition found")
	}
	partition := partitionR.(*util.Partition)
	for _, p := range partition.Parts {
		if strings.HasPrefix(p.Label, "node") && p.Contains(id) {
			return p.Label
		}
	}
	panic("Replica not found")
}

func FromToIsolated(p Partition) testlib.Condition {
	return func(e *types.Event, c *testlib.Context) bool {
		message, ok := c.GetMessage(e)
		if !ok {
			return false
		}
		from := getPartLabel(c, message.From)
		to := getPartLabel(c, message.To)

		return isolates(p, from, to)
	}
}

func isolates(p Partition, a string, b string) bool {
	for _, part := range p {
		if partContains(part, a) && partContains(part, b) {
			return false
		}
	}
	return true
}

func partContains(part []string, i string) bool {
	for _, v := range part {
		if v == i {
			return true
		}
	}

	return false
}
