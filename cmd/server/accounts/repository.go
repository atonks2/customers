// Copyright 2020 The Moov Authors
// Use of this source code is governed by an Apache License
// license that can be found in the LICENSE file.

package accounts

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/moov-io/base"
	"github.com/moov-io/customers/client"

	"github.com/go-kit/kit/log"
)

type Repository interface {
	getCustomerAccounts(customerID string) ([]*client.Account, error)
	createCustomerAccount(customerID, userID string, req *createAccountRequest) (*client.Account, error)
	deactivateCustomerAccount(accountID string) error
}

func NewRepo(logger log.Logger, db *sql.DB) *sqlAccountRepository {
	return &sqlAccountRepository{logger: logger, db: db}
}

type sqlAccountRepository struct {
	db     *sql.DB
	logger log.Logger
}

func (r *sqlAccountRepository) Close() error {
	return r.db.Close()
}

func (r *sqlAccountRepository) getCustomerAccounts(customerID string) ([]*client.Account, error) {
	query := `select account_id, masked_account_number, routing_number, type, holder_type from accounts
where customer_id = ? and deleted_at is null;`
	stmt, err := r.db.Prepare(query)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	rows, err := stmt.Query(customerID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var out []*client.Account
	for rows.Next() {
		var a client.Account
		if err := rows.Scan(&a.Id, &a.MaskedAccountNumber, &a.RoutingNumber, &a.Type, &a.HolderType); err != nil {
			return nil, err
		}
		out = append(out, &a)
	}
	return out, nil
}

func (r *sqlAccountRepository) createCustomerAccount(customerID, userID string, req *createAccountRequest) (*client.Account, error) {
	account := &client.Account{
		Id:                  base.ID(),
		MaskedAccountNumber: req.AccountNumber,
		RoutingNumber:       req.RoutingNumber,
		Type:                req.Type,
		HolderType:          req.HolderType,
	}
	query := `insert into accounts (account_id, customer_id, user_id, encrypted_account_number, hashed_account_number, masked_account_number, routing_number, type, holder_type, created_at) values (?, ?, ?, ?, ?, ?, ?, ?, ?, ?);`
	stmt, err := r.db.Prepare(query)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	_, err = stmt.Exec(account.Id, customerID, userID, "", "", "", account.RoutingNumber, account.Type, account.HolderType, time.Now())
	if err != nil {
		return nil, fmt.Errorf("problem creating account=%s: %v", account.Id, err)
	}
	return account, nil
}

func (r *sqlAccountRepository) deactivateCustomerAccount(accountID string) error {
	query := `update accounts set deleted_at = ? where account_id = ? and deleted_at is null;`
	stmt, err := r.db.Prepare(query)
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(time.Now(), accountID)
	if err == sql.ErrNoRows {
		return nil
	}
	return err
}
