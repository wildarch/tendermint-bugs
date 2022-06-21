package tendermintbugs

import (
	"github.com/netrixframework/netrix/log"
	"github.com/netrixframework/netrix/testlib"
	"github.com/netrixframework/netrix/types"
	"github.com/netrixframework/tendermint-testing/util"
)

func dropMessageLoudly(e *types.Event, c *testlib.Context) (message []*types.Message) {
	m, ok := util.GetMessageFromEvent(e, c)
	if ok {
		c.Logger().With(log.LogParams{
			"from":   getPartLabel(c, m.From),
			"to":     getPartLabel(c, m.To),
			"type":   m.Type,
			"height": m.Height(),
			"round":  m.Round(),
		}).Info("Dropping message")
	} else {
		c.Logger().Warn("Dropping message with unknown height/round")
	}
	return
}
