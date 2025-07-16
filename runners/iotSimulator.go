package runners

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"time"

	"github.com/mbabinski218/BaaS-IoT-platform/blockchain"
	"github.com/mbabinski218/BaaS-IoT-platform/configs"
)

func RunIotSimulator(bc *blockchain.Client) error {
	for {
		blockNumber, err := bc.GetBlockNumber()
		if err != nil {
			log.Fatalf("failed to get block number: %v", err)
		}

		if blockNumber > 0 && bc.IsMining() {
			break
		}

		time.Sleep(1 * time.Second)
	}

	cmd := exec.Command(configs.Envs.IotSimulatorCommand, configs.Envs.IotSimulatorParams)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Start(); err != nil {
		return fmt.Errorf("failed to start IoT simulator: %w", err)
	}

	return nil
}
