/*
 * Copyright (C) 2018 eeonevision
 *
 * Permission is hereby granted, free of charge, to any person obtaining a copy of
 * this software and associated documentation files (the "Software"), to deal in
 * the Software without restriction, including without limitation the rights to
 * use, copy, modify, merge, publish, distribute, sublicense, and/or sell copies of
 * the Software, and to permit persons to whom the Software is furnished to do so,
 * subject to the following conditions:
 *
 * The above copyright notice and this permission notice shall be included in all
 * copies or substantial portions of the Software.
 *
 * THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
 * IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY, FITNESS
 * FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR
 * COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER
 * IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN
 * CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.
 */

package state

import (
	"errors"

	"github.com/eeonevision/anychaindb/crypto"
)

//go:generate msgp

// Account struct keeps account related fields.
type Account struct {
	ID     string `msg:"_id" json:"_id" mapstructure:"_id" bson:"_id"`
	PubKey string `msg:"public_key" json:"public_key" mapstructure:"public_key" bson:"public_key"`
}

const accountsCollection = "accounts"

// AddAccount method adds new account if all checks were passed.
func (s *State) AddAccount(account *Account) error {
	if s.HasAccount(account.ID) {
		return errors.New("account exists")
	}
	return s.SetAccount(account)
}

// SetAccount method adds account in state.
func (s *State) SetAccount(account *Account) error {
	return s.DB.C(accountsCollection).Insert(account)
}

// HasAccount method checks if account exists in state or not exists.
func (s *State) HasAccount(id string) bool {
	if res, _ := s.GetAccount(id); res != nil {
		return true
	}
	return false
}

// GetAccount method returns account from accounts collection by given accoutn id.
func (s *State) GetAccount(id string) (*Account, error) {
	var result *Account
	return result, s.DB.C(accountsCollection).FindId(id).One(&result)
}

// GetAccountPubKey method returns public key by given account id.
func (s *State) GetAccountPubKey(id string) (*crypto.Key, error) {
	acc, err := s.GetAccount(id)
	if err != nil {
		return nil, err
	}
	return crypto.NewFromStrings(acc.PubKey, "")
}

// ListAccounts method returns all accounts from the state.
func (s *State) ListAccounts() (result []*Account, err error) {
	return result, s.DB.C(accountsCollection).Find(nil).All(&result)
}

// SearchAccounts method returns accounts by given search query, limit and offset parameters.
func (s *State) SearchAccounts(query interface{}, limit, offset int) (result []*Account, err error) {
	return result, s.DB.C(accountsCollection).Find(query).Skip(offset).Limit(limit).All(&result)
}
