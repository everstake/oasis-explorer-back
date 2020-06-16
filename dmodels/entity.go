package dmodels

import "time"

const (
	EntityNodesView                   = "entity_nodes_view"
	EntityActiveDepositorsCounterView = "entity_active_depositors_counter_view"
)

type EntityNodesContainer struct {
	Nodes []EntityNode
}

func (c EntityNodesContainer) IsEmpty() bool {
	return len(c.Nodes) == 0
}

func (c EntityNodesContainer) IsNode(accountID string) bool {
	return len(c.Nodes) == 1 && c.Nodes[0].NodeID == accountID
}

func (c EntityNodesContainer) IsEntity(accountID string) bool {
	return len(c.Nodes) > 0 && c.Nodes[0].EntityID == accountID
}

func (c EntityNodesContainer) GetEntityAddress() string {
	if len(c.Nodes) == 0 {
		return ""
	}

	return c.Nodes[0].EntityID
}

func (c EntityNodesContainer) GetEntity() (resp EntityNode) {
	if len(c.Nodes) == 0 {
		return resp
	}

	var blocksCount, blockSignaturesCount uint64
	var lastProposedBlockTime, lastBlockSignatureTime time.Time

	for i := range c.Nodes {
		blocksCount += c.Nodes[i].BlocksCount
		blockSignaturesCount += c.Nodes[i].BlockSignaturesCount

		if c.Nodes[i].LastSignatureTime.After(lastBlockSignatureTime) {
			lastBlockSignatureTime = c.Nodes[i].LastSignatureTime
		}

		if c.Nodes[i].LastBlockTime.After(lastProposedBlockTime) {
			lastProposedBlockTime = c.Nodes[i].LastBlockTime
		}
	}

	resp = c.Nodes[len(c.Nodes)-1]
	resp.BlocksCount = blocksCount
	resp.BlockSignaturesCount = blockSignaturesCount
	resp.LastBlockTime = lastBlockSignatureTime
	resp.LastSignatureTime = lastBlockSignatureTime

	return resp
}

type EntityNode struct {
	EntityID             string    `db="reg_entity_id"`
	NodeID               string    `db="reg_id"`
	ConsensusAddress     string    `db="reg_consensus_address"`
	LastRegBlock         uint64    `db="blk_lvl"`
	CreatedTime          time.Time `db="created_time"`
	Expiration           uint64    `db="reg_expiration"`
	BlocksCount          uint64    `db="blocks"`
	LastBlockTime        time.Time `db="last_block_time"`
	BlockSignaturesCount uint64    `db="signatures"`
	LastSignatureTime    time.Time `db="last_signature_time"`
}

func (n EntityNode) GetLastActiveTime() time.Time {
	if n.LastSignatureTime.After(n.LastBlockTime) {
		return n.LastSignatureTime
	}

	return n.LastBlockTime
}
