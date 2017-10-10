package main

import (
	"flag"
	"github.com/buaazp/fasthttprouter"
	"github.com/valyala/fasthttp"
	"log"
)

var (
	addr = flag.String("addr", ":8080", "TCP address to listen on.")
	db   = flag.String("db", "test.db", "The database file to use.")

	ACCOUNT_BUCKET = []byte("ACCOUNT")
	PAYMENT_BUCKET = []byte("PAYMENT")
)

//------------------------------------------------------------------------------
// Model
//------------------------------------------------------------------------------

func (db *BoltState) Initialise() error {
	// Ensure buckets exist.
	if err := db.SyncBucket(ACCOUNT_BUCKET); err != nil {
		return err
	}
	if err := db.SyncBucket(PAYMENT_BUCKET); err != nil {
		return err
	}

	// Create fixtures
	bob := Account{
		Id:       "bob123",
		Owner:    "bob",
		Balance:  100,
		Currency: "PHP",
	}
	alice := Account{
		Id:       "alice456",
		Owner:    "alice",
		Balance:  0.01,
		Currency: "PHP",
	}

	bobRecord, err := bob.AsRecord()
	if err != nil {
		return err
	}
	aliceRecord, err := alice.AsRecord()
	if err != nil {
		return err
	}

	db.PutRecord(ACCOUNT_BUCKET, bobRecord)
	db.PutRecord(ACCOUNT_BUCKET, aliceRecord)

	return nil
}

//------------------------------------------------------------------------------

func (db *BoltState) runServer() *BoltState {
	router := fasthttprouter.New()
	router.GET("/v1/accounts", db.AccountsGET())
	router.GET("/v1/payments", db.GetPayments())
	router.POST("/v1/payments", db.PostPayment())

	if err := fasthttp.ListenAndServe(*addr, router.Handler); err != nil {
		log.Fatalf("Error in ListenAndServe: %s", err)
	}

	return nil
}

func main() {
	flag.Parse()

	db := OpenDB(*db)
	db.Initialise()

	db.runServer()
}
