package integration_test

import (
	"context"
	"time"

	"github.com/Vulcanize/ipld-eth-db-validator/pkg/validator"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/statediff/indexer/node"
	"github.com/ethereum/go-ethereum/statediff/indexer/postgres"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	integration "github.com/vulcanize/ipld-eth-server/test"
)

const trail = 2
const head = 1

var randomAddr = common.HexToAddress("0x1C3ab14BBaD3D99F4203bd7a11aCB94882050E6f")

var _ = Describe("Integration test", func() {
	ctx := context.Background()

	var contract *integration.ContractDeployed
	var contractErr error
	sleepInterval := 2 * time.Second

	Describe("Validate all blocks", func() {
		address := "0x1111111111111111111111111111111111111112"

		err := sendEthTransactions(address, sleepInterval, 5)
		Expect(err).ToNot(HaveOccurred())

		contract, contractErr = integration.DeployContract()
		Expect(contractErr).ToNot(HaveOccurred())
		time.Sleep(sleepInterval)

		err = sendEthTransactions(address, sleepInterval, 5)
		Expect(err).ToNot(HaveOccurred())

		_, err = integration.DestroyContract(contract.Address)
		Expect(err).ToNot(HaveOccurred())

		time.Sleep(sleepInterval)

		err = sendEthTransactions(address, sleepInterval, 5)
		Expect(err).ToNot(HaveOccurred())

		// Run validator
		db, _ := setupDB()
		srvc := validator.NewService(db, head, trail, validator.IntegrationTestChainConfig)
		_, err = srvc.Start(ctx)
		Expect(err).ToNot(HaveOccurred())

	})
})

func sendEthTransactions(address string, sleepInterval time.Duration, n int) error {
	for i := 0; i < n; i++ {
		if _, err := integration.SendEth(address, "0.01"); err != nil {
			return err
		}
		time.Sleep(sleepInterval)
	}
	return nil
}
func setupDB() (*postgres.DB, error) {
	uri := postgres.DbConnectionString(postgres.ConnectionParams{
		User:     "vdbm",
		Password: "password",
		Hostname: "localhost",
		Name:     "vulcanize_testing",
		Port:     8077,
	})
	return validator.NewDB(uri, postgres.ConnectionConfig{}, node.Info{})
}
