package tendermintbugs

import (
	"fmt"
	"strconv"

	"github.com/netrixframework/netrix/testlib"
	"github.com/netrixframework/netrix/types"
)

func trackNodeHeightRound(e *types.Event, c *testlib.Context) (messages []*types.Message, handled bool) {
	// Parse event to retrieve height and round
	eType, ok := e.Type.(*types.GenericEventType)
	if !ok {
		return
	}
	if eType.T != "newStep" {
		return
	}
	heightS, ok := eType.Params["height"]
	if !ok {
		return
	}
	height, err := strconv.Atoi(heightS)
	if err != nil {
		return
	}
	roundS, ok := eType.Params["round"]
	if !ok {
		return
	}
	round, err := strconv.Atoi(roundS)
	if err != nil {
		return
	}

	c.Vars.Set(nodeHeightKey(e.Replica), height)
	c.Vars.Set(nodeRoundKey(e.Replica), round)
	return
}

func IsMessageWithSenderHeightRound(height, round int) testlib.Condition {
	return func(e *types.Event, c *testlib.Context) bool {
		if !e.IsMessageSend() {
			return false
		}
		senderHeight, ok := c.Vars.GetInt(nodeHeightKey(e.Replica))
		if !ok {
			height = 1
		}
		senderRound, ok := c.Vars.GetInt(nodeRoundKey(e.Replica))
		if !ok {
			round = 0
		}
		return height == senderHeight && round == senderRound
	}
}

func nodeHeightKey(e types.ReplicaID) string {
	return fmt.Sprintf("node_height_%s", e)
}

func nodeRoundKey(e types.ReplicaID) string {
	return fmt.Sprintf("node_round_%s", e)
}
