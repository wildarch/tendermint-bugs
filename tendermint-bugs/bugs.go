package tendermintbugs

import (
	"time"

	"github.com/netrixframework/netrix/testlib"
	"github.com/netrixframework/tendermint-testing/common"
	"github.com/netrixframework/tendermint-testing/util"
)

type Partition = [][]string

// Does not inject any faults, provided to verify communication between test server and nodes.
// Since there are no obstructions, and a new empty commit is created every second, we expect height to steadily increase.
// On the author's laptop, this test reaches height 21.
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

	// In first prevote step, create the following partitioning
	// {0,2} {1} {3}
	partition := Partition{
		{"node0", "node2"},
		{"node1"},
		{"node3"},
	}
	// Drop the prevotes sent across partitions at H=1/R=0 (the very first round)
	filters.AddFilter(testlib.If(
		common.IsMessageType(util.Prevote).
			And(IsMessageWithSenderHeightRound(1, 0).
				And(FromToIsolated(partition)).
				And(IsBeforeLivenessCheck)),
	).Then(dropMessageLoudly))

	testCase := testlib.NewTestCase(
		"Bug001",
		2*time.Minute,
		sm,
		filters,
	)
	testCase.SetupFunc(common.Setup(sysParams, labelNodes, SetupLivenessTimer(time.Minute)))
	return testCase
}

// Isolates two nodes in the first precommit step.
// This will prevent the nodes from consensus.
// After 1 minute the network is healed, and we expect them to recover (which in our tests they do not).
func Bug002(sysParams *common.SystemParams) *testlib.TestCase {
	sm := testlib.NewStateMachine()
	// Expect to be able to commit at least one block.
	// This may happen after the liveness timer expires and the network is healed.
	sm.Builder().On(common.HeightReached(2), testlib.SuccessStateLabel)

	filters := testlib.NewFilterSet()
	filters.AddFilter(trackNodeHeightRound)

	// In first precommit step, create the following partitioning
	// {0,2} {1} {3}
	partition := Partition{
		{"node0", "node2"},
		{"node1"},
		{"node3"},
	}
	// Drop the the precommits send across partitions at H=1/R=0 (the very first round)
	filters.AddFilter(testlib.If(
		common.IsMessageType(util.Precommit).
			And(IsMessageWithSenderHeightRound(1, 0).
				And(FromToIsolated(partition)).
				And(IsBeforeLivenessCheck)),
	).Then(dropMessageLoudly))

	testCase := testlib.NewTestCase(
		"Bug002",
		2*time.Minute,
		sm,
		filters,
	)
	testCase.SetupFunc(common.Setup(sysParams, labelNodes, SetupLivenessTimer(time.Minute)))
	return testCase
}
