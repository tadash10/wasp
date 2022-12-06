package client

import (
	"net/http"
	"github.com/iotaledger/wasp/packages/webapi/v1/routes"

	"github.com/iotaledger/wasp/packages/isc"
)

// ActivateChain sends a request to activate a chain in the wasp node
func (c *WaspClient) ActivateChain(chID *isc.ChainID) error {
	return c.do(http.MethodPost, routes.ActivateChain(chID.String()), nil, nil)
}

// DeactivateChain sends a request to deactivate a chain in the wasp node
func (c *WaspClient) DeactivateChain(chID *isc.ChainID) error {
	return c.do(http.MethodPost, routes.DeactivateChain(chID.String()), nil, nil)
}
