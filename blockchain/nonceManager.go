package blockchain

import (
	"context"
	"fmt"
	"math/big"
	"sync"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
)

type NonceManager struct {
	mutex       sync.Mutex
	current     uint64
	initialized bool
}

func NewNonceManager() *NonceManager {
	return &NonceManager{
		mutex:       sync.Mutex{},
		current:     0,
		initialized: false,
	}
}

func (nonceManager *NonceManager) Init(client *ethclient.Client, from common.Address, ctx context.Context) error {
	nonceManager.mutex.Lock()
	defer nonceManager.mutex.Unlock()

	if nonceManager.initialized {
		return nil
	}

	nonce, err := client.PendingNonceAt(ctx, from)
	if err != nil {
		return fmt.Errorf("failed to get initial nonce: %w", err)
	}

	nonceManager.current = nonce
	nonceManager.initialized = true
	return nil
}

func (nonceManager *NonceManager) Next() *big.Int {
	nonceManager.mutex.Lock()
	defer nonceManager.mutex.Unlock()

	nonce := nonceManager.current
	nonceManager.current++

	return big.NewInt(int64(nonce))
}
