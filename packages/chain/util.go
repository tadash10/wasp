package chain

import (
	"github.com/iotaledger/wasp/packages/chain/eventmessages"

	"github.com/iotaledger/hive.go/core/logger"
	"github.com/iotaledger/trie.go/trie"
	"github.com/iotaledger/wasp/packages/hashing"
	"github.com/iotaledger/wasp/packages/isc"
	"github.com/iotaledger/wasp/packages/publisher"
)

// LogStateTransition also used in testing
func LogStateTransition(blockIndex uint32, outputID string, rootCommitment trie.VCommitment, reqids []isc.RequestID, log *logger.Logger) {
	if blockIndex > 0 {
		log.Infof("STATE TRANSITION TO #%d. Requests: %d, chain output: %s", blockIndex, len(reqids), outputID)
		log.Debugf("STATE TRANSITION. Root commitment: %s", rootCommitment)
	} else {
		log.Infof("ORIGIN STATE SAVED. State output id: %s", outputID)
		log.Debugf("ORIGIN STATE SAVED. Root commitment: %s", rootCommitment)
	}
}

// LogGovernanceTransition
func LogGovernanceTransition(blockIndex uint32, outputID string, rootCommitment trie.VCommitment, log *logger.Logger) {
	log.Infof("GOVERNANCE TRANSITION. State index #%d, anchor output: %s", blockIndex, outputID)
	log.Debugf("GOVERNANCE TRANSITION. Root commitment: %s", rootCommitment)
}

func PublishRequestsSettled(chainID *isc.ChainID, stateIndex uint32, reqids []isc.RequestID) {
	for _, reqid := range reqids {
		message := eventmessages.RequestOut{
			RequestID:  reqid.String(),
			RequestIDs: len(reqids),
			StateIndex: stateIndex,
		}

		publisher.Publish(eventmessages.NewChainEventRequestOut(chainID.String(), message))
	}
}

func PublishStateTransition(chainID *isc.ChainID, stateOutput *isc.AliasOutputWithID, reqIDsLength int) {
	stateHash, _ := hashing.HashValueFromBytes(stateOutput.GetStateMetadata())

	message := eventmessages.StateUpdate{
		StateIndex:    stateOutput.GetStateIndex(),
		RequestIDs:    reqIDsLength,
		StateOutputID: isc.OID(stateOutput.ID()),
		StateHash:     stateHash.String(),
	}

	publisher.Publish(eventmessages.NewChainEventStateUpdate(chainID.String(), message))
}

func PublishGovernanceTransition(stateOutput *isc.AliasOutputWithID) {
	stateHash, _ := hashing.HashValueFromBytes(stateOutput.GetStateMetadata())
	chainID := isc.ChainIDFromAliasID(stateOutput.GetAliasID())

	message := eventmessages.Rotation{
		StateIndex:    stateOutput.GetStateIndex(),
		StateOutputID: isc.OID(stateOutput.ID()),
		StateHash:     stateHash.String(),
	}

	publisher.Publish(eventmessages.NewChainEventRotation(chainID.String(), message))
}
