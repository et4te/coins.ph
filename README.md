# coins.ph golang exercise

This exercise is an implementation of the Coins Django Task using golang.

The libraries used are fasthttp, fasthttprouter, boltdb, blake2b and testify/assert.

The accounts are stored keyed by the id of the account.

The payments are stored keyed by the blake2b hash of the payment encoded as JSON and the integrity of the data is verified upon retrieval.

To generate the executable run `go build` relative to your gopath source tree.

To run the tests run `go test`.

To run the executable run `coins.ph`.
