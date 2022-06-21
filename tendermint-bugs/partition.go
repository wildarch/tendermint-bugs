package tendermintbugs

import (
	"fmt"
	"strings"

	"github.com/netrixframework/netrix/log"
	"github.com/netrixframework/netrix/testlib"
	"github.com/netrixframework/netrix/types"
	"github.com/netrixframework/tendermint-testing/util"
)

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
