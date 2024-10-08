package ethereum

import (
	"context"
	"errors"
	"fmt"
	"math/big"
	"os"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
)

type EthereumClient struct {
	client *ethclient.Client
	auth   *bind.TransactOpts
}

func NewEthereumClient() (*EthereumClient, error) {

	// privateKey := os.Getenv("TREASURY_PRIVATE_KEY")
	// if privateKey == "" {
	// 	return nil, errors.New("TREASURY_PRIVATE_KEY not set in .env file")
	// }
	passphrase := os.Getenv("ACCOUNT_PASSWORD")
	if passphrase == "" {
		return nil, errors.New("ACCOUNT_PASSWORD not set in .env file")
	}

	// 连接到本地geth节点
	rpcURL := os.Getenv("ETHEREUM_RPC_URL")
	if rpcURL == "" {
		return nil, errors.New("ETHEREUM_RPC_URL cannot be empty in environment variables")
	}
	client, err := ethclient.Dial(rpcURL)
	if err != nil {
		return nil, err
	}

	keystoreFilePath := os.Getenv("KEYSTORE_FILE_PATH")
	if keystoreFilePath == "" {
		return nil, errors.New("KEYSTORE_FILE_PATH cannot be empty in environment variables")
	}
	keyFile, err := os.Open(keystoreFilePath)
	if err != nil {
		return nil, err
	}
	defer keyFile.Close()

	chainID, err := client.ChainID(context.Background())
	if err != nil {
		return nil, err
	}
	transactOpts, err := bind.NewTransactorWithChainID(keyFile, passphrase, chainID)
	if err != nil {
		return nil, err
	}

	return &EthereumClient{
		client: client,
		auth:   transactOpts,
	}, nil
}

func (ec *EthereumClient) CheckBalance(ctx context.Context, tokenAddress string) (*big.Int, error) {
	token, err := NewERC20(common.HexToAddress(tokenAddress), ec.client)
	if err != nil {
		return nil, err
	}

	// 从环境变量中读取账户地址
	account := os.Getenv("TREASURY_ADDRESS")
	if account == "" {
		return nil, errors.New("TREASURY_ADDRESS not set in environment variables")
	}

	balance, err := token.BalanceOf(&bind.CallOpts{Context: ctx}, common.HexToAddress(account))
	if err != nil {
		return nil, err
	}

	return balance, nil
}

func (ec *EthereumClient) TransferTokens(ctx context.Context, tokenAddress string, to string, amount *big.Int) error {
	token, err := NewERC20(common.HexToAddress(tokenAddress), ec.client)
	if err != nil {
		return err
	}

	// 检查发送者的余额
	balance, err := token.BalanceOf(&bind.CallOpts{Context: ctx}, ec.auth.From)
	if err != nil {
		return err
	}
	fmt.Println("balance:", balance)

	if balance.Cmp(amount) < 0 {
		return errors.New("insufficient treasury balance, cannot approve")
	}

	// 执行转账
	_, err = token.Transfer(ec.auth, common.HexToAddress(to), amount)
	return err
}

func (ec *EthereumClient) GetTokenDecimals(ctx context.Context, tokenAddress string) (uint8, error) {
	token, err := NewERC20(common.HexToAddress(tokenAddress), ec.client)
	if err != nil {
		return 0, err
	}

	decimals, err := token.Decimals(&bind.CallOpts{Context: ctx})
	if err != nil {
		return 0, err
	}

	return decimals, nil
}
