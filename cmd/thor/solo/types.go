// Copyright (c) 2018 The VeChainThor developers

// Distributed under the GNU Lesser General Public License v3.0 software license, see the accompanying
// file LICENSE or <https://www.gnu.org/licenses/lgpl-3.0.html>

package solo

import (
	"github.com/vechain/thor/comm"
	"github.com/vechain/thor/thor"
)

// Communicator in solo is a fake one just for api handler.
type Communicator struct {
}

// PeersStats returns nil solo doesn't join p2p network.
func (comm *Communicator) PeersStats() []*comm.PeerStats {
	return nil
}

// BFTEngine is a fake bft engine for solo.
type BFTEngine struct {
}

func (engine *BFTEngine) Finalized() thor.Bytes32 {
	return thor.Bytes32{}
}
