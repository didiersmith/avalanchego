// Copyright (C) 2019-2021, Ava Labs, Inc. All rights reserved.
// See the file LICENSE for licensing terms.

package transactions

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/ava-labs/avalanchego/ids"
	"github.com/ava-labs/avalanchego/utils/constants"
	"github.com/ava-labs/avalanchego/vms/platformvm/transactions/signed"
	"github.com/ava-labs/avalanchego/vms/platformvm/transactions/unsigned"

	p_validator "github.com/ava-labs/avalanchego/vms/platformvm/validator"
)

func TestPrimaryValidatorSet(t *testing.T) {
	// Initialize the chain state
	nodeID0 := ids.GenerateTestNodeID()
	node0Weight := uint64(1)
	vdr0 := &currentValidatorImpl{
		addValidator: signed.ValidatorAndID{
			UnsignedAddValidatorTx: &unsigned.AddValidatorTx{
				Validator: p_validator.Validator{
					Wght: node0Weight,
				},
			},
			TxID: ids.ID{'v', 'd', 'r', '0', 'I', 'D'},
		},
	}

	nodeID1 := ids.GenerateTestNodeID()
	node1Weight := uint64(2)
	vdr1 := &currentValidatorImpl{
		addValidator: signed.ValidatorAndID{
			UnsignedAddValidatorTx: &unsigned.AddValidatorTx{
				Validator: p_validator.Validator{
					Wght: node1Weight,
				},
			},
			TxID: ids.ID{'v', 'd', 'r', '1', 'I', 'D'},
		},
	}

	nodeID2 := ids.GenerateTestNodeID()
	node2Weight := uint64(2)
	vdr2 := &currentValidatorImpl{
		addValidator: signed.ValidatorAndID{
			UnsignedAddValidatorTx: &unsigned.AddValidatorTx{
				Validator: p_validator.Validator{
					Wght: node2Weight,
				},
			},
		},
	}

	cs := &currentStaker{
		ValidatorsByNodeID: map[ids.NodeID]*currentValidatorImpl{
			nodeID0: vdr0,
			nodeID1: vdr1,
			nodeID2: vdr2,
		},
	}
	nodeID3 := ids.GenerateTestNodeID()

	{
		// Apply the on-chain validator set to [vdrs]
		vdrs, err := cs.ValidatorSet(constants.PrimaryNetworkID)
		assert.NoError(t, err)

		// Validate that the state was applied and the old state was cleared
		assert.EqualValues(t, 3, vdrs.Len())
		assert.EqualValues(t, node0Weight+node1Weight+node2Weight, vdrs.Weight())
		gotNode0Weight, exists := vdrs.GetWeight(nodeID0)
		assert.True(t, exists)
		assert.EqualValues(t, node0Weight, gotNode0Weight)
		gotNode1Weight, exists := vdrs.GetWeight(nodeID1)
		assert.True(t, exists)
		assert.EqualValues(t, node1Weight, gotNode1Weight)
		gotNode2Weight, exists := vdrs.GetWeight(nodeID2)
		assert.True(t, exists)
		assert.EqualValues(t, node2Weight, gotNode2Weight)
		_, exists = vdrs.GetWeight(nodeID3)
		assert.False(t, exists)
	}

	{
		// Apply the on-chain validator set again
		vdrs, err := cs.ValidatorSet(constants.PrimaryNetworkID)
		assert.NoError(t, err)

		// The state should be the same
		assert.EqualValues(t, 3, vdrs.Len())
		assert.EqualValues(t, node0Weight+node1Weight+node2Weight, vdrs.Weight())
		gotNode0Weight, exists := vdrs.GetWeight(nodeID0)
		assert.True(t, exists)
		assert.EqualValues(t, node0Weight, gotNode0Weight)
		gotNode1Weight, exists := vdrs.GetWeight(nodeID1)
		assert.True(t, exists)
		assert.EqualValues(t, node1Weight, gotNode1Weight)
		gotNode2Weight, exists := vdrs.GetWeight(nodeID2)
		assert.True(t, exists)
		assert.EqualValues(t, node2Weight, gotNode2Weight)
	}
}

func TestSubnetValidatorSet(t *testing.T) {
	subnetID := ids.GenerateTestID()

	// Initialize the chain state
	nodeID0 := ids.GenerateTestNodeID()
	node0Weight := uint64(1)
	vdr0 := &currentValidatorImpl{
		validatorImpl: validatorImpl{
			subnets: map[ids.ID]signed.SubnetValidatorAndID{
				subnetID: {
					UnsignedAddSubnetValidator: &unsigned.AddSubnetValidatorTx{
						Validator: p_validator.SubnetValidator{
							Validator: p_validator.Validator{
								Wght: node0Weight,
							},
						},
					},
					TxID: ids.ID{'s', 'u', 'b', 'v', 'd', 'r', '0', 'I', 'D'},
				},
			},
		},
	}

	nodeID1 := ids.GenerateTestNodeID()
	node1Weight := uint64(2)
	vdr1 := &currentValidatorImpl{
		validatorImpl: validatorImpl{
			subnets: map[ids.ID]signed.SubnetValidatorAndID{
				subnetID: {
					UnsignedAddSubnetValidator: &unsigned.AddSubnetValidatorTx{
						Validator: p_validator.SubnetValidator{
							Validator: p_validator.Validator{
								Wght: node1Weight,
							},
						},
					},
					TxID: ids.ID{'s', 'u', 'b', 'v', 'd', 'r', '1', 'I', 'D'},
				},
			},
		},
	}

	nodeID2 := ids.GenerateTestNodeID()
	node2Weight := uint64(2)
	vdr2 := &currentValidatorImpl{
		validatorImpl: validatorImpl{
			subnets: map[ids.ID]signed.SubnetValidatorAndID{
				subnetID: {
					UnsignedAddSubnetValidator: &unsigned.AddSubnetValidatorTx{
						Validator: p_validator.SubnetValidator{
							Validator: p_validator.Validator{
								Wght: node2Weight,
							},
						},
					},
					TxID: ids.ID{'s', 'u', 'b', 'v', 'd', 'r', '2', 'I', 'D'},
				},
			},
		},
	}

	cs := &currentStaker{
		ValidatorsByNodeID: map[ids.NodeID]*currentValidatorImpl{
			nodeID0: vdr0,
			nodeID1: vdr1,
			nodeID2: vdr2,
		},
	}

	nodeID3 := ids.GenerateTestNodeID()

	{
		// Apply the on-chain validator set to [vdrs]
		vdrs, err := cs.ValidatorSet(subnetID)
		assert.NoError(t, err)

		// Validate that the state was applied and the old state was cleared
		assert.EqualValues(t, 3, vdrs.Len())
		assert.EqualValues(t, node0Weight+node1Weight+node2Weight, vdrs.Weight())
		gotNode0Weight, exists := vdrs.GetWeight(nodeID0)
		assert.True(t, exists)
		assert.EqualValues(t, node0Weight, gotNode0Weight)
		gotNode1Weight, exists := vdrs.GetWeight(nodeID1)
		assert.True(t, exists)
		assert.EqualValues(t, node1Weight, gotNode1Weight)
		gotNode2Weight, exists := vdrs.GetWeight(nodeID2)
		assert.True(t, exists)
		assert.EqualValues(t, node2Weight, gotNode2Weight)
		_, exists = vdrs.GetWeight(nodeID3)
		assert.False(t, exists)
	}

	{
		// Apply the on-chain validator set again
		vdrs, err := cs.ValidatorSet(subnetID)
		assert.NoError(t, err)

		// The state should be the same
		assert.EqualValues(t, 3, vdrs.Len())
		assert.EqualValues(t, node0Weight+node1Weight+node2Weight, vdrs.Weight())
		gotNode0Weight, exists := vdrs.GetWeight(nodeID0)
		assert.True(t, exists)
		assert.EqualValues(t, node0Weight, gotNode0Weight)
		gotNode1Weight, exists := vdrs.GetWeight(nodeID1)
		assert.True(t, exists)
		assert.EqualValues(t, node1Weight, gotNode1Weight)
		gotNode2Weight, exists := vdrs.GetWeight(nodeID2)
		assert.True(t, exists)
		assert.EqualValues(t, node2Weight, gotNode2Weight)
	}
}
