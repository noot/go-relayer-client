# go-relayer-client

Client library for [go-relayer](https://github.com/AthanorLabs/go-relayer), used for submitting transactions to the relayer.

## Requirements

- go1.18+

## Usage

```go
// ... code omitted
req := &common.SubmitTransactionRequest{
	From:  from,
	To:    to,
	Value: big.NewInt(0),
	Gas:   big.NewInt(679639582),
	Nonce: nonce,
	Data:  calldata,
}

c := client.NewClient("http://localhost:8545")
resp, err := c.SubmitTransaction(req)
if err != nil {
	panic(err)
}

fmt.Println("sent tx to relayer: hash", resp.TxHash)
```

For a full example, see `examples/main.go`.

### Run local example 

The example in `examples/main.go` works with [go-relayer](https://github.com/AthanorLabs/go-relayer) locally.

To try it:

1. Install and run ganache: 
```bash
npm i -g ganache
ganache --deterministic --accounts=50
```

2. Clone and build `go-relayer`:
```bash
git clone https://github.com/AthanorLabs/go-relayer
cd go-relayer
make build
./bin/relayer --dev
```

3. Clone and build `go-relayer-client`:
```bash
git clone https://github.com/AthanorLabs/go-relayer-client
cd go-relayer-client
make build
./bin/example
```

You should see logs in the client as follows:
```bash
sent tx to relayer: hash 0xaaa5b2a84d1c4e4e5c251fe1ccb6059115267c6432d031615d20f5dae2771ddf
tx successful!
```

You should see logs in the relayer as follows:
```bash
2022-09-29T14:08:03.117-0400	INFO	cmd	cmd/main.go:132	starting relayer with ethereum endpoint http://localhost:8545 and chain ID 1337
2022-09-29T14:08:03.232-0400	INFO	cmd	cmd/main.go:208	deployed MinimalForwarder.sol to 0xe78A0F7E598Cc8b0Bb87894B0F60dD2a88d6a8Ab
2022-09-29T14:08:03.233-0400	INFO	rpc	rpc/server.go:62	starting RPC server on http://localhost:7799
2022-09-29T14:08:16.685-0400	INFO	relayer	relayer/relayer.go:109	submitted transaction 0xaaa5b2a84d1c4e4e5c251fe1ccb6059115267c6432d031615d20f5dae2771ddf
```