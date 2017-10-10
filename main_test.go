package main

import (
	"github.com/stretchr/testify/assert"
	"github.com/valyala/fasthttp"
	"testing"
)

type TestClient struct {
	HTTPClient *fasthttp.Client
}

func (cli *TestClient) GET(url string) ([]byte, error) {
	req := fasthttp.AcquireRequest()
	req.SetRequestURI(url)
	req.Header.SetMethod("GET")
	rsp := fasthttp.AcquireResponse()
	cli.HTTPClient.Do(req, rsp)
	return rsp.Body(), nil
}

func (cli *TestClient) POST(url string, data []byte) ([]byte, error) {
	req := fasthttp.AcquireRequest()
	req.SetRequestURI(url)
	req.Header.SetMethod("POST")
	req.SetBody(data)
	rsp := fasthttp.AcquireResponse()
	cli.HTTPClient.Do(req, rsp)
	return rsp.Body(), nil
}

func TestMain(t *testing.T) {
	assert := assert.New(t)

	db := OpenDB(*db)
	db.Initialise()

	go (func() { db.runServer() })()

	cli := TestClient{&fasthttp.Client{}}

	// Accounts testing
	records, err := db.GetRecords(ACCOUNT_BUCKET)
	if err != nil {
		t.Fatal(err)
		return
	}

	var accounts []*Account
	for i := 0; i < len(records); i++ {
		account, err := AccountFromRecord(records[i])
		if err != nil {
			t.Fatal(err)
			return
		}
		accounts = append(accounts, []*Account{account}...)
	}

	testAccountsJSON, err := cli.GET("http://127.0.0.1:8080/v1/accounts")
	if err != nil {
		t.Fatal(err)
	}

	testAccounts, err := AccountsFromJSON(testAccountsJSON)
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(testAccounts, accounts)

	// Payments testing

	// POST a new (valid payment)
	payment := Payment{
		Account:     "bob123",
		Amount:      0.0001,
		FromAccount: "alice456",
		Direction:   "outgoing",
	}

	bytes, err := payment.AsJSON()
	if err != nil {
		t.Fatal(err)
	}

	_, err = cli.POST("http://localhost:8080/v1/payments", bytes)
	if err != nil {
		t.Fatal(err)
	}

	//
	records, err = db.GetRecords(PAYMENT_BUCKET)
	if err != nil {
		t.Fatal(err)
	}

	var payments []*Payment
	for i := 0; i < len(records); i++ {
		payment, err := PaymentFromRecord(records[i])
		if err != nil {
			t.Fatal(err)
		}
		payments = append(payments, []*Payment{payment}...)
	}

	testPaymentsJSON, err := cli.GET("http://127.0.0.1:8080/v1/payments")
	if err != nil {
		t.Fatal(err)
	}

	testPayments, err := PaymentsFromJSON(testPaymentsJSON)
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(testPayments, payments)
}
