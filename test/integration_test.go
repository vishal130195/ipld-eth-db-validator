package integration_test

import (
	"context"
	"time"

	"github.com/ethereum/go-ethereum/statediff/indexer/node"
	"github.com/ethereum/go-ethereum/statediff/indexer/postgres"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/Vulcanize/ipld-eth-db-validator/pkg/validator"

	integration "github.com/vulcanize/ipld-eth-server/test"
)

const (
	trail        = 2
	head         = 1
	toBlock      = 0
	sendEthCount = 5
)

var _ = Describe("Integration test", func() {
	ctx := context.Background()

	var contract *integration.ContractDeployed
	var contractErr error
	sleepInterval := 2 * time.Second

	Describe("Validate all blocks", func() {
		address := "0x1111111111111111111111111111111111111112"

		err := sendMultipleEth(address, sleepInterval, sendEthCount)
		Expect(err).ToNot(HaveOccurred())

		contract, contractErr = integration.DeployContract()
		Expect(contractErr).ToNot(HaveOccurred())
		time.Sleep(sleepInterval)

		err = sendMultipleEth(address, sleepInterval, sendEthCount)
		Expect(err).ToNot(HaveOccurred())

		_, err = integration.DestroyContract(contract.Address)
		Expect(err).ToNot(HaveOccurred())

		time.Sleep(sleepInterval)

		err = sendMultipleEth(address, sleepInterval, sendEthCount)
		Expect(err).ToNot(HaveOccurred())

		// Run validator
		db, _ := setupDB()
		srvc := validator.NewService(db, head, trail, toBlock, validator.IntegrationTestChainConfig)
		_, err = srvc.Start(ctx)
		Expect(err).ToNot(HaveOccurred())

	})
})

func sendMultipleEth(address string, sleepInterval time.Duration, n int) error {
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
