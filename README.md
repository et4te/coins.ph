# coins.ph golang exercise

This exercise is an implementation of the Coins Django Task using golang.

The libraries used are fasthttp, fasthttprouter, boltdb, blake2b and testify/assert.

The accounts are stored keyed by the id of the account.

The payments are stored keyed by the blake2b hash of the payment encoded as JSON and the integrity of the data is verified upon retrieval.

More could be done in order to make the payments occur in a fault tolerant way but in order to keep the exercise simple the checks are done only when creating a payment (rather than as one atomic transaction).

To generate the executable run `go build` relative to your gopath source tree.

To run the tests run `go test`.

To run the executable run `coins.ph`.
