package main

import (
	"context"
	"crypto/ecdsa"
	"fmt"
	"math/big"

	"github.com/athanorlabs/atomic-swap/ethereum/block"
	"github.com/athanorlabs/go-relayer/common"
	mock "github.com/athanorlabs/go-relayer/examples/mock_recipient"
	"github.com/athanorlabs/go-relayer/impls/mforwarder"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	ethcommon "github.com/ethereum/go-ethereum/common"
	ethtypes "github.com/ethereum/go-ethereum/core/types"
	ethcrypto "github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"

	"github.com/athanorlabs/go-relayer-client"
)

func main() {
	auth, ec, pk, chainID := setup()
	forwarderAddr := ethcommon.HexToAddress("0xe78A0F7E598Cc8b0Bb87894B0F60dD2a88d6a8Ab")
	recipientAddress, err := deployMockRecipient(auth, ec, forwarderAddr)
	if err != nil {
		panic(err)
	}

	// transfer to recipient - only for setup, not needed if contract already is funded
	value := big.NewInt(1000000)
	fee := big.NewInt(10000)

	transferTx := ethtypes.NewTransaction(
		0,
		recipientAddress,
		value,
		100000,
		big.NewInt(679639582),
		nil,
	)

	transferTx, err = ethtypes.SignTx(transferTx, ethtypes.LatestSignerForChainID(chainID), pk)
	if err != nil {
		panic(err)
	}
	err = ec.SendTransaction(context.Background(), transferTx)
	if err != nil {
		panic(err)
	}
	_, err = block.WaitForReceipt(context.Background(), ec, transferTx.Hash())
	if err != nil {
		panic(err)
	}

	// form withdraw transaction to relay
	abi, err := mock.MockMetaData.GetAbi()
	if err != nil {
		panic(err)
	}

	calldata, err := abi.Pack("withdraw", value, fee)
	if err != nil {
		panic(err)
	}

	forwarder, err := mforwarder.NewMinimalForwarder(forwarderAddr, ec)
	if err != nil {
		panic(err)
	}

	// generate fresh address
	key, err := common.GenerateKey()
	if err != nil {
		panic(err)
	}

	nonce, err := forwarder.GetNonce(&bind.CallOpts{}, key.Address())
	if err != nil {
		panic(err)
	}

	req := &mforwarder.IMinimalForwarderForwardRequest{
		From:  key.Address(),
		To:    recipientAddress,
		Value: big.NewInt(0),
		Gas:   big.NewInt(679639582), // TODO: fetch from ethclient
		Nonce: nonce,
		Data:  calldata,
	}

	name := "MinimalForwarder"
	version := "0.0.1"

	domainSeparator, err := common.GetEIP712DomainSeparator(name, version, chainID, forwarderAddr)
	if err != nil {
		panic(err)
	}

	digest, err := common.GetForwardRequestDigestToSign(
		req,
		domainSeparator,
		nil,
	)
	if err != nil {
		panic(err)
	}

	sig, err := key.Sign(digest)
	if err != nil {
		panic(err)
	}

	rpcReq := &common.SubmitTransactionRequest{
		From:      req.From,
		To:        req.To,
		Value:     req.Value,
		Gas:       req.Gas,
		Nonce:     req.Nonce,
		Data:      req.Data,
		Signature: sig,
	}

	// submit transaction to relayer
	c := client.NewClient(client.DefaultLocalRelayerEndpoint)
	resp, err := c.SubmitTransaction(rpcReq)
	if err != nil {
		panic(err)
	}

	fmt.Println("sent tx to relayer: hash", resp.TxHash)
	receipt, err := block.WaitForReceipt(context.Background(), ec, resp.TxHash)
	if err != nil {
		panic(err)
	}

	if receipt.Status == 1 {
		fmt.Println("tx successful!")
	} else {
		fmt.Println("tx failed :(")
	}
}

func setup() (*bind.TransactOpts, *ethclient.Client, *ecdsa.PrivateKey, *big.Int) {
	const ethEndpoint = "http://localhost:8545"

	ec, err := ethclient.Dial(ethEndpoint)
	if err != nil {
		panic(err)
	}

	chainID, err := ec.ChainID(context.Background())
	if err != nil {
		panic(err)
	}

	// NODE_OPTIONS="--max_old_space_size=8192" ganache --deterministic --accounts=50
	pk, err := ethcrypto.HexToECDSA("6cbed15c793ce57650b9877cf6fa156fbef513c4e6134f022a85b1ffdd59b2a1")
	if err != nil {
		panic(err)
	}

	auth, err := bind.NewKeyedTransactorWithChainID(pk, chainID)
	if err != nil {
		panic(err)
	}

	return auth, ec, pk, chainID
}

func deployMockRecipient(
	auth *bind.TransactOpts,
	conn *ethclient.Client,
	forwarderAddr ethcommon.Address,
) (ethcommon.Address, error) {
	address, tx, _, err := mock.DeployMock(auth, conn, forwarderAddr)
	if err != nil {
		panic(err)
	}

	_, err = block.WaitForReceipt(context.Background(), conn, tx.Hash())
	if err != nil {
		panic(err)
	}

	return address, nil
}
