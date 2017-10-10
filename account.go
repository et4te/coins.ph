package main

import (
	"encoding/json"
	"fmt"
	"github.com/valyala/fasthttp"
)

//------------------------------------------------------------------------------
// Model
//------------------------------------------------------------------------------

type Account struct {
	Id       string  `json:"id"`
	Owner    string  `json:"owner"`
	Balance  float64 `json:"balance"`
	Currency string  `json:"currency"`
}

func (db *BoltState) GetAccount(k string) (*Account, error) {
	rec, err := db.GetRecord(ACCOUNT_BUCKET, []byte(k))
	if err != nil {
		return nil, err
	}
	return AccountFromRecord(rec)
}

//------------------------------------------------------------------------------
// Routes
//------------------------------------------------------------------------------

func (db *BoltState) AccountsGET() fasthttp.RequestHandler {
	return func(ctx *fasthttp.RequestCtx) {
		ctx.SetContentType("application/json")

		records, err := db.GetRecords(ACCOUNT_BUCKET)
		if err != nil {
			fmt.Fprint(ctx, err)
			return
		}

		var accounts []*Account
		for i := 0; i < len(records); i++ {
			account, err := AccountFromRecord(records[i])
			if err != nil {
				fmt.Fprint(ctx, err)
				return
			}
			accounts = append(accounts, []*Account{account}...)
		}

		if len(accounts) > 0 {
			if bytes, err := AccountsToJSON(accounts); err != nil {
				fmt.Fprint(ctx, err)
				return
			} else {
				fmt.Fprintf(ctx, "%s", bytes)
				return
			}
		} else {
			fmt.Fprintf(ctx, "[]")
			return
		}
	}
}

//------------------------------------------------------------------------------
// JSON Encoding / Decoding
//------------------------------------------------------------------------------

func AccountsToJSON(accounts []*Account) ([]byte, error) {
	return json.Marshal(accounts)
}

func AccountsFromJSON(bytes []byte) ([]*Account, error) {
	var a []*Account
	if err := json.Unmarshal(bytes, &a); err != nil {
		return nil, err
	}
	return a, nil
}

func AccountFromJSON(bytes []byte) (*Account, error) {
	var a Account
	if err := json.Unmarshal(bytes, &a); err != nil {
		return nil, err
	}
	return &a, nil
}

func (a *Account) AsJSON() ([]byte, error) {
	return json.Marshal(a)
}

//------------------------------------------------------------------------------
// From / To persistent record
//------------------------------------------------------------------------------

func AccountFromRecord(r *Record) (*Account, error) {
	return AccountFromJSON(r.Value)
}

func (a *Account) AsRecord() (*Record, error) {
	bytes, err := a.AsJSON()
	if err != nil {
		return nil, err
	}
	return &Record{Key: []byte(a.Id), Value: bytes}, nil
}
