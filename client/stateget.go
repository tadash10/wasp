package client

import (
	"encoding/hex"
	"github.com/iotaledger/wasp/packages/webapi/v1/routes"
	"net/http"

	"github.com/iotaledger/wasp/packages/isc"
)

// StateGet fetches the raw value associated with the given key in the chain state
func (c *WaspClient) StateGet(chainID *isc.ChainID, key string) ([]byte, error) {
	var res []byte
	if err := c.do(http.MethodGet, routes.StateGet(chainID.String(), hex.EncodeToString([]byte(key))), nil, &res); err != nil {
		return nil, err
	}
	return res, nil
}
