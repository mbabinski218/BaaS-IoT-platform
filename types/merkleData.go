package types

import (
	"bytes"

	"github.com/cbergoon/merkletree"
)

type MerkleData struct {
	Hash [32]byte
}

func (m MerkleData) CalculateHash() ([]byte, error) {
	return m.Hash[:], nil
}

func (m MerkleData) Equals(other merkletree.Content) (bool, error) {
	o := other.(MerkleData)
	return bytes.Equal(m.Hash[:], o.Hash[:]), nil
}
