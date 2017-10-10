package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/valyala/fasthttp"
)

type Payment struct {
	Account     string  `json:"account"`
	Amount      float64 `json:"amount"`
	FromAccount string  `json:"from_account"`
	Direction   string  `json:"direction"`
}

func (db *BoltState) GetPayments() fasthttp.RequestHandler {
	return func(ctx *fasthttp.RequestCtx) {
		ctx.SetContentType("application/json")

		records, err := db.GetRecords(PAYMENT_BUCKET)
		if err != nil {
			fmt.Fprint(ctx, err)
			return
		}

		var payments []*Payment
		for i := 0; i < len(records); i++ {
			payment, err := PaymentFromRecord(records[i])
			if err != nil {
				fmt.Fprint(ctx, err)
				return
			}
			payments = append(payments, []*Payment{payment}...)
		}

		if len(payments) > 0 {
			if bytes, err := PaymentsToJSON(payments); err != nil {
				fmt.Fprint(ctx, err)
				return
			} else {
				fmt.Fprintf(ctx, "%s", bytes)
				return
			}
		} else {
			fmt.Fprint(ctx, "[]")
			return
		}
	}
}

func (db *BoltState) PostPayment() fasthttp.RequestHandler {
	return func(ctx *fasthttp.RequestCtx) {
		ctx.SetContentType("text/plain; charset=utf8")

		body := ctx.Request.Body()
		payment, err := PaymentFromJSON(body)
		if err != nil {
			fmt.Fprint(ctx, err)
			return
		}

		rec, err := payment.AsRecord(db)
		if err != nil {
			fmt.Fprint(ctx, err)
			return
		}

		if err := db.PutRecord(PAYMENT_BUCKET, rec); err != nil {
			fmt.Fprint(ctx, err)
		}

		return
	}
}

//------------------------------------------------------------------------------
// JSON encoding / decoding
//------------------------------------------------------------------------------

func PaymentsToJSON(payments []*Payment) ([]byte, error) {
	return json.Marshal(payments)
}

func PaymentsFromJSON(bytes []byte) ([]*Payment, error) {
	var ps []*Payment
	if err := json.Unmarshal(bytes, &ps); err != nil {
		return nil, err
	}
	return ps, nil
}

func PaymentFromJSON(bytes []byte) (*Payment, error) {
	var a Payment
	if err := json.Unmarshal(bytes, &a); err != nil {
		return nil, err
	}
	return &a, nil
}

func (a *Payment) AsJSON() ([]byte, error) {
	return json.Marshal(a)
}

//------------------------------------------------------------------------------
// From / to persistent record
//------------------------------------------------------------------------------

func PaymentFromRecord(r *Record) (*Payment, error) {
	if r.Verify([]byte("payment")) {
		return PaymentFromJSON(r.Value)
	} else {
		return nil, errors.New("Failed to verify record integrity.")
	}
}

func (a *Payment) AsRecord(db *BoltState) (*Record, error) {
	// Check that to_account exists
	toAccount, err := db.GetAccount(a.Account)
	if err != nil {
		return nil, err
	}
	if toAccount == nil {
		return nil, errors.New("to_account does not exist.")
	}

	// Check that from_account exists
	fromAccount, err := db.GetAccount(a.FromAccount)
	if err != nil {
		return nil, err
	}
	if fromAccount == nil {
		return nil, errors.New("from_account does not exist.")
	}

	// Check that to_account != from_account
	if toAccount == fromAccount {
		return nil, errors.New("Cannot make a payment to self.")
	}

	// Check that from_account balance > amount
	if fromAccount.Balance <= a.Amount {
		return nil, errors.New("Payment amount exceeds available balance.")
	}

	// Fields are valid
	bytes, err := a.AsJSON()
	if err != nil {
		return nil, err
	}

	h, err := HashBlake2b([]byte("payment"), bytes)
	if err != nil {
		return nil, err
	}

	return &Record{Key: h[:], Value: bytes}, nil
}
