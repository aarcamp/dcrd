// Copyright (c) 2014-2016 The btcsuite developers
// Copyright (c) 2015-2019 The Decred developers
// Use of this source code is governed by an ISC
// license that can be found in the LICENSE file.

package rpcclient

import (
	"encoding/hex"
	"encoding/json"

	"github.com/decred/dcrd/chaincfg/chainhash"
	"github.com/decred/dcrd/dcrjson/v3"
	"github.com/decred/dcrd/dcrutil/v2"
	"github.com/decred/dcrd/hdkeychain/v2"
	chainjson "github.com/decred/dcrd/rpc/jsonrpc/types"
	"github.com/decred/dcrd/wire"
	walletjson "github.com/decred/dcrwallet/rpc/jsonrpc/types"
)

// *****************************
// Transaction Listing Functions
// *****************************

// FutureGetTransactionResult is a future promise to deliver the result
// of a GetTransactionAsync RPC invocation (or an applicable error).
type FutureGetTransactionResult chan *response

// Receive waits for the response promised by the future and returns detailed
// information about a wallet transaction.
func (r FutureGetTransactionResult) Receive() (*walletjson.GetTransactionResult, error) {
	res, err := receiveFuture(r)
	if err != nil {
		return nil, err
	}

	// Unmarshal result as a gettransaction result object
	var getTx walletjson.GetTransactionResult
	err = json.Unmarshal(res, &getTx)
	if err != nil {
		return nil, err
	}

	return &getTx, nil
}

// GetTransactionAsync returns an instance of a type that can be used to get the
// result of the RPC at some future time by invoking the Receive function on
// the returned instance.
//
// See GetTransaction for the blocking version and more details.
func (c *Client) GetTransactionAsync(txHash *chainhash.Hash) FutureGetTransactionResult {
	hash := ""
	if txHash != nil {
		hash = txHash.String()
	}

	cmd := walletjson.NewGetTransactionCmd(hash, nil)
	return c.sendCmd(cmd)
}

// GetTransaction returns detailed information about a wallet transaction.
//
// See GetRawTransaction to return the raw transaction instead.
func (c *Client) GetTransaction(txHash *chainhash.Hash) (*walletjson.GetTransactionResult, error) {
	return c.GetTransactionAsync(txHash).Receive()
}

// FutureListTransactionsResult is a future promise to deliver the result of a
// ListTransactionsAsync, ListTransactionsCountAsync, or
// ListTransactionsCountFromAsync RPC invocation (or an applicable error).
type FutureListTransactionsResult chan *response

// Receive waits for the response promised by the future and returns a list of
// the most recent transactions.
func (r FutureListTransactionsResult) Receive() ([]walletjson.ListTransactionsResult, error) {
	res, err := receiveFuture(r)
	if err != nil {
		return nil, err
	}

	// Unmarshal result as an array of listtransaction result objects.
	var transactions []walletjson.ListTransactionsResult
	err = json.Unmarshal(res, &transactions)
	if err != nil {
		return nil, err
	}

	return transactions, nil
}

// ListTransactionsAsync returns an instance of a type that can be used to get
// the result of the RPC at some future time by invoking the Receive function on
// the returned instance.
//
// See ListTransactions for the blocking version and more details.
func (c *Client) ListTransactionsAsync(account string) FutureListTransactionsResult {
	cmd := walletjson.NewListTransactionsCmd(&account, nil, nil, nil)
	return c.sendCmd(cmd)
}

// ListTransactions returns a list of the most recent transactions.
//
// See the ListTransactionsCount and ListTransactionsCountFrom to control the
// number of transactions returned and starting point, respectively.
func (c *Client) ListTransactions(account string) ([]walletjson.ListTransactionsResult, error) {
	return c.ListTransactionsAsync(account).Receive()
}

// ListTransactionsCountAsync returns an instance of a type that can be used to
// get the result of the RPC at some future time by invoking the Receive
// function on the returned instance.
//
// See ListTransactionsCount for the blocking version and more details.
func (c *Client) ListTransactionsCountAsync(account string, count int) FutureListTransactionsResult {
	cmd := walletjson.NewListTransactionsCmd(&account, &count, nil, nil)
	return c.sendCmd(cmd)
}

// ListTransactionsCount returns a list of the most recent transactions up
// to the passed count.
//
// See the ListTransactions and ListTransactionsCountFrom functions for
// different options.
func (c *Client) ListTransactionsCount(account string, count int) ([]walletjson.ListTransactionsResult, error) {
	return c.ListTransactionsCountAsync(account, count).Receive()
}

// ListTransactionsCountFromAsync returns an instance of a type that can be used
// to get the result of the RPC at some future time by invoking the Receive
// function on the returned instance.
//
// See ListTransactionsCountFrom for the blocking version and more details.
func (c *Client) ListTransactionsCountFromAsync(account string, count, from int) FutureListTransactionsResult {
	cmd := walletjson.NewListTransactionsCmd(&account, &count, &from, nil)
	return c.sendCmd(cmd)
}

// ListTransactionsCountFrom returns a list of the most recent transactions up
// to the passed count while skipping the first 'from' transactions.
//
// See the ListTransactions and ListTransactionsCount functions to use defaults.
func (c *Client) ListTransactionsCountFrom(account string, count, from int) ([]walletjson.ListTransactionsResult, error) {
	return c.ListTransactionsCountFromAsync(account, count, from).Receive()
}

// FutureListUnspentResult is a future promise to deliver the result of a
// ListUnspentAsync, ListUnspentMinAsync, ListUnspentMinMaxAsync, or
// ListUnspentMinMaxAddressesAsync RPC invocation (or an applicable error).
type FutureListUnspentResult chan *response

// Receive waits for the response promised by the future and returns all
// unspent wallet transaction outputs returned by the RPC call.  If the
// future was returned by a call to ListUnspentMinAsync, ListUnspentMinMaxAsync,
// or ListUnspentMinMaxAddressesAsync, the range may be limited by the
// parameters of the RPC invocation.
func (r FutureListUnspentResult) Receive() ([]walletjson.ListUnspentResult, error) {
	res, err := receiveFuture(r)
	if err != nil {
		return nil, err
	}

	// Unmarshal result as an array of listunspent results.
	var unspent []walletjson.ListUnspentResult
	err = json.Unmarshal(res, &unspent)
	if err != nil {
		return nil, err
	}

	return unspent, nil
}

// ListUnspentAsync returns an instance of a type that can be used to get
// the result of the RPC at some future time by invoking the Receive function
// on the returned instance.
//
// See ListUnspent for the blocking version and more details.
func (c *Client) ListUnspentAsync() FutureListUnspentResult {
	cmd := walletjson.NewListUnspentCmd(nil, nil, nil)
	return c.sendCmd(cmd)
}

// ListUnspentMinAsync returns an instance of a type that can be used to get
// the result of the RPC at some future time by invoking the Receive function
// on the returned instance.
//
// See ListUnspentMin for the blocking version and more details.
func (c *Client) ListUnspentMinAsync(minConf int) FutureListUnspentResult {
	cmd := walletjson.NewListUnspentCmd(&minConf, nil, nil)
	return c.sendCmd(cmd)
}

// ListUnspentMinMaxAsync returns an instance of a type that can be used to get
// the result of the RPC at some future time by invoking the Receive function
// on the returned instance.
//
// See ListUnspentMinMax for the blocking version and more details.
func (c *Client) ListUnspentMinMaxAsync(minConf, maxConf int) FutureListUnspentResult {
	cmd := walletjson.NewListUnspentCmd(&minConf, &maxConf, nil)
	return c.sendCmd(cmd)
}

// ListUnspentMinMaxAddressesAsync returns an instance of a type that can be
// used to get the result of the RPC at some future time by invoking the Receive
// function on the returned instance.
//
// See ListUnspentMinMaxAddresses for the blocking version and more details.
func (c *Client) ListUnspentMinMaxAddressesAsync(minConf, maxConf int, addrs []dcrutil.Address) FutureListUnspentResult {
	addrStrs := make([]string, 0, len(addrs))
	for _, a := range addrs {
		addrStrs = append(addrStrs, a.Address())
	}

	cmd := walletjson.NewListUnspentCmd(&minConf, &maxConf, &addrStrs)
	return c.sendCmd(cmd)
}

// ListUnspent returns all unspent transaction outputs known to a wallet, using
// the default number of minimum and maximum number of confirmations as a
// filter (1 and 9999999, respectively).
func (c *Client) ListUnspent() ([]walletjson.ListUnspentResult, error) {
	return c.ListUnspentAsync().Receive()
}

// ListUnspentMin returns all unspent transaction outputs known to a wallet,
// using the specified number of minimum conformations and default number of
// maximum confirmations (9999999) as a filter.
func (c *Client) ListUnspentMin(minConf int) ([]walletjson.ListUnspentResult, error) {
	return c.ListUnspentMinAsync(minConf).Receive()
}

// ListUnspentMinMax returns all unspent transaction outputs known to a wallet,
// using the specified number of minimum and maximum number of confirmations as
// a filter.
func (c *Client) ListUnspentMinMax(minConf, maxConf int) ([]walletjson.ListUnspentResult, error) {
	return c.ListUnspentMinMaxAsync(minConf, maxConf).Receive()
}

// ListUnspentMinMaxAddresses returns all unspent transaction outputs that pay
// to any of specified addresses in a wallet using the specified number of
// minimum and maximum number of confirmations as a filter.
func (c *Client) ListUnspentMinMaxAddresses(minConf, maxConf int, addrs []dcrutil.Address) ([]walletjson.ListUnspentResult, error) {
	return c.ListUnspentMinMaxAddressesAsync(minConf, maxConf, addrs).Receive()
}

// FutureListSinceBlockResult is a future promise to deliver the result of a
// ListSinceBlockAsync or ListSinceBlockMinConfAsync RPC invocation (or an
// applicable error).
type FutureListSinceBlockResult chan *response

// Receive waits for the response promised by the future and returns all
// transactions added in blocks since the specified block hash, or all
// transactions if it is nil.
func (r FutureListSinceBlockResult) Receive() (*walletjson.ListSinceBlockResult, error) {
	res, err := receiveFuture(r)
	if err != nil {
		return nil, err
	}

	// Unmarshal result as a listsinceblock result object.
	var listResult walletjson.ListSinceBlockResult
	err = json.Unmarshal(res, &listResult)
	if err != nil {
		return nil, err
	}

	return &listResult, nil
}

// ListSinceBlockAsync returns an instance of a type that can be used to get
// the result of the RPC at some future time by invoking the Receive function on
// the returned instance.
//
// See ListSinceBlock for the blocking version and more details.
func (c *Client) ListSinceBlockAsync(blockHash *chainhash.Hash) FutureListSinceBlockResult {
	var hash *string
	if blockHash != nil {
		hash = dcrjson.String(blockHash.String())
	}

	cmd := walletjson.NewListSinceBlockCmd(hash, nil, nil)
	return c.sendCmd(cmd)
}

// ListSinceBlock returns all transactions added in blocks since the specified
// block hash, or all transactions if it is nil, using the default number of
// minimum confirmations as a filter.
//
// See ListSinceBlockMinConf to override the minimum number of confirmations.
func (c *Client) ListSinceBlock(blockHash *chainhash.Hash) (*walletjson.ListSinceBlockResult, error) {
	return c.ListSinceBlockAsync(blockHash).Receive()
}

// ListSinceBlockMinConfAsync returns an instance of a type that can be used to
// get the result of the RPC at some future time by invoking the Receive
// function on the returned instance.
//
// See ListSinceBlockMinConf for the blocking version and more details.
func (c *Client) ListSinceBlockMinConfAsync(blockHash *chainhash.Hash, minConfirms int) FutureListSinceBlockResult {
	var hash *string
	if blockHash != nil {
		hash = dcrjson.String(blockHash.String())
	}

	cmd := walletjson.NewListSinceBlockCmd(hash, &minConfirms, nil)
	return c.sendCmd(cmd)
}

// ListSinceBlockMinConf returns all transactions added in blocks since the
// specified block hash, or all transactions if it is nil, using the specified
// number of minimum confirmations as a filter.
//
// See ListSinceBlock to use the default minimum number of confirmations.
func (c *Client) ListSinceBlockMinConf(blockHash *chainhash.Hash, minConfirms int) (*walletjson.ListSinceBlockResult, error) {
	return c.ListSinceBlockMinConfAsync(blockHash, minConfirms).Receive()
}

// **************************
// Transaction Send Functions
// **************************

// FutureLockUnspentResult is a future promise to deliver the error result of a
// LockUnspentAsync RPC invocation.
type FutureLockUnspentResult chan *response

// Receive waits for the response promised by the future and returns the result
// of locking or unlocking the unspent output(s).
func (r FutureLockUnspentResult) Receive() error {
	_, err := receiveFuture(r)
	return err
}

// LockUnspentAsync returns an instance of a type that can be used to get the
// result of the RPC at some future time by invoking the Receive function on the
// returned instance.
//
// See LockUnspent for the blocking version and more details.
func (c *Client) LockUnspentAsync(unlock bool, ops []*wire.OutPoint) FutureLockUnspentResult {
	outputs := make([]chainjson.TransactionInput, len(ops))
	for i, op := range ops {
		outputs[i] = chainjson.TransactionInput{
			Txid: op.Hash.String(),
			Vout: op.Index,
			Tree: op.Tree,
		}
	}
	cmd := walletjson.NewLockUnspentCmd(unlock, outputs)
	return c.sendCmd(cmd)
}

// LockUnspent marks outputs as locked or unlocked, depending on the value of
// the unlock bool.  When locked, the unspent output will not be selected as
// input for newly created, non-raw transactions, and will not be returned in
// future ListUnspent results, until the output is marked unlocked again.
//
// If unlock is false, each outpoint in ops will be marked locked.  If unlocked
// is true and specific outputs are specified in ops (len != 0), exactly those
// outputs will be marked unlocked.  If unlocked is true and no outpoints are
// specified, all previous locked outputs are marked unlocked.
//
// The locked or unlocked state of outputs are not written to disk and after
// restarting a wallet process, this data will be reset (every output unlocked).
//
// NOTE: While this method would be a bit more readable if the unlock bool was
// reversed (that is, LockUnspent(true, ...) locked the outputs), it has been
// left as unlock to keep compatibility with the reference client API and to
// avoid confusion for those who are already familiar with the lockunspent RPC.
func (c *Client) LockUnspent(unlock bool, ops []*wire.OutPoint) error {
	return c.LockUnspentAsync(unlock, ops).Receive()
}

// FutureListLockUnspentResult is a future promise to deliver the result of a
// ListLockUnspentAsync RPC invocation (or an applicable error).
type FutureListLockUnspentResult chan *response

// Receive waits for the response promised by the future and returns the result
// of all currently locked unspent outputs.
func (r FutureListLockUnspentResult) Receive() ([]*wire.OutPoint, error) {
	res, err := receiveFuture(r)
	if err != nil {
		return nil, err
	}

	// Unmarshal as an array of transaction inputs.
	var inputs []chainjson.TransactionInput
	err = json.Unmarshal(res, &inputs)
	if err != nil {
		return nil, err
	}

	// Create a slice of outpoints from the transaction input structs.
	ops := make([]*wire.OutPoint, len(inputs))
	for i, input := range inputs {
		sha, err := chainhash.NewHashFromStr(input.Txid)
		if err != nil {
			return nil, err
		}
		ops[i] = wire.NewOutPoint(sha, input.Vout, input.Tree) // Decred TODO
	}

	return ops, nil
}

// ListLockUnspentAsync returns an instance of a type that can be used to get
// the result of the RPC at some future time by invoking the Receive function on
// the returned instance.
//
// See ListLockUnspent for the blocking version and more details.
func (c *Client) ListLockUnspentAsync() FutureListLockUnspentResult {
	cmd := walletjson.NewListLockUnspentCmd()
	return c.sendCmd(cmd)
}

// ListLockUnspent returns a slice of outpoints for all unspent outputs marked
// as locked by a wallet.  Unspent outputs may be marked locked using
// LockOutput.
func (c *Client) ListLockUnspent() ([]*wire.OutPoint, error) {
	return c.ListLockUnspentAsync().Receive()
}

// FutureSendToAddressResult is a future promise to deliver the result of a
// SendToAddressAsync RPC invocation (or an applicable error).
type FutureSendToAddressResult chan *response

// Receive waits for the response promised by the future and returns the hash
// of the transaction sending the passed amount to the given address.
func (r FutureSendToAddressResult) Receive() (*chainhash.Hash, error) {
	res, err := receiveFuture(r)
	if err != nil {
		return nil, err
	}

	// Unmarshal result as a string.
	var txHash string
	err = json.Unmarshal(res, &txHash)
	if err != nil {
		return nil, err
	}

	return chainhash.NewHashFromStr(txHash)
}

// SendToAddressAsync returns an instance of a type that can be used to get the
// result of the RPC at some future time by invoking the Receive function on the
// returned instance.
//
// See SendToAddress for the blocking version and more details.
func (c *Client) SendToAddressAsync(address dcrutil.Address, amount dcrutil.Amount) FutureSendToAddressResult {
	addr := address.Address()
	cmd := walletjson.NewSendToAddressCmd(addr, amount.ToCoin(), nil, nil)
	return c.sendCmd(cmd)
}

// SendToAddress sends the passed amount to the given address.
//
// See SendToAddressComment to associate comments with the transaction in the
// wallet.  The comments are not part of the transaction and are only internal
// to the wallet.
//
// NOTE: This function requires to the wallet to be unlocked.  See the
// WalletPassphrase function for more details.
func (c *Client) SendToAddress(address dcrutil.Address, amount dcrutil.Amount) (*chainhash.Hash, error) {
	return c.SendToAddressAsync(address, amount).Receive()
}

// SendToAddressCommentAsync returns an instance of a type that can be used to
// get the result of the RPC at some future time by invoking the Receive
// function on the returned instance.
//
// See SendToAddressComment for the blocking version and more details.
func (c *Client) SendToAddressCommentAsync(address dcrutil.Address,
	amount dcrutil.Amount, comment,
	commentTo string) FutureSendToAddressResult {

	addr := address.Address()
	cmd := walletjson.NewSendToAddressCmd(addr, amount.ToCoin(), &comment,
		&commentTo)
	return c.sendCmd(cmd)
}

// SendToAddressComment sends the passed amount to the given address and stores
// the provided comment and comment to in the wallet.  The comment parameter is
// intended to be used for the purpose of the transaction while the commentTo
// parameter is intended to be used for who the transaction is being sent to.
//
// The comments are not part of the transaction and are only internal
// to the wallet.
//
// See SendToAddress to avoid using comments.
//
// NOTE: This function requires to the wallet to be unlocked.  See the
// WalletPassphrase function for more details.
func (c *Client) SendToAddressComment(address dcrutil.Address, amount dcrutil.Amount, comment, commentTo string) (*chainhash.Hash, error) {
	return c.SendToAddressCommentAsync(address, amount, comment,
		commentTo).Receive()
}

// FutureSendFromResult is a future promise to deliver the result of a
// SendFromAsync, SendFromMinConfAsync, or SendFromCommentAsync RPC invocation
// (or an applicable error).
type FutureSendFromResult chan *response

// Receive waits for the response promised by the future and returns the hash
// of the transaction sending amount to the given address using the provided
// account as a source of funds.
func (r FutureSendFromResult) Receive() (*chainhash.Hash, error) {
	res, err := receiveFuture(r)
	if err != nil {
		return nil, err
	}

	// Unmarshal result as a string.
	var txHash string
	err = json.Unmarshal(res, &txHash)
	if err != nil {
		return nil, err
	}

	return chainhash.NewHashFromStr(txHash)
}

// SendFromAsync returns an instance of a type that can be used to get the
// result of the RPC at some future time by invoking the Receive function on the
// returned instance.
//
// See SendFrom for the blocking version and more details.
func (c *Client) SendFromAsync(fromAccount string, toAddress dcrutil.Address, amount dcrutil.Amount) FutureSendFromResult {
	addr := toAddress.Address()
	cmd := walletjson.NewSendFromCmd(fromAccount, addr, amount.ToCoin(), nil,
		nil, nil)
	return c.sendCmd(cmd)
}

// SendFrom sends the passed amount to the given address using the provided
// account as a source of funds.  Only funds with the default number of minimum
// confirmations will be used.
//
// See SendFromMinConf and SendFromComment for different options.
//
// NOTE: This function requires to the wallet to be unlocked.  See the
// WalletPassphrase function for more details.
func (c *Client) SendFrom(fromAccount string, toAddress dcrutil.Address, amount dcrutil.Amount) (*chainhash.Hash, error) {
	return c.SendFromAsync(fromAccount, toAddress, amount).Receive()
}

// SendFromMinConfAsync returns an instance of a type that can be used to get
// the result of the RPC at some future time by invoking the Receive function on
// the returned instance.
//
// See SendFromMinConf for the blocking version and more details.
func (c *Client) SendFromMinConfAsync(fromAccount string, toAddress dcrutil.Address, amount dcrutil.Amount, minConfirms int) FutureSendFromResult {
	addr := toAddress.Address()
	cmd := walletjson.NewSendFromCmd(fromAccount, addr, amount.ToCoin(),
		&minConfirms, nil, nil)
	return c.sendCmd(cmd)
}

// SendFromMinConf sends the passed amount to the given address using the
// provided account as a source of funds.  Only funds with the passed number of
// minimum confirmations will be used.
//
// See SendFrom to use the default number of minimum confirmations and
// SendFromComment for additional options.
//
// NOTE: This function requires to the wallet to be unlocked.  See the
// WalletPassphrase function for more details.
func (c *Client) SendFromMinConf(fromAccount string, toAddress dcrutil.Address, amount dcrutil.Amount, minConfirms int) (*chainhash.Hash, error) {
	return c.SendFromMinConfAsync(fromAccount, toAddress, amount,
		minConfirms).Receive()
}

// SendFromCommentAsync returns an instance of a type that can be used to get
// the result of the RPC at some future time by invoking the Receive function on
// the returned instance.
//
// See SendFromComment for the blocking version and more details.
func (c *Client) SendFromCommentAsync(fromAccount string,
	toAddress dcrutil.Address, amount dcrutil.Amount, minConfirms int,
	comment, commentTo string) FutureSendFromResult {

	addr := toAddress.Address()
	cmd := walletjson.NewSendFromCmd(fromAccount, addr, amount.ToCoin(),
		&minConfirms, &comment, &commentTo)
	return c.sendCmd(cmd)
}

// SendFromComment sends the passed amount to the given address using the
// provided account as a source of funds and stores the provided comment and
// comment to in the wallet.  The comment parameter is intended to be used for
// the purpose of the transaction while the commentTo parameter is intended to
// be used for who the transaction is being sent to.  Only funds with the passed
// number of minimum confirmations will be used.
//
// See SendFrom and SendFromMinConf to use defaults.
//
// NOTE: This function requires to the wallet to be unlocked.  See the
// WalletPassphrase function for more details.
func (c *Client) SendFromComment(fromAccount string, toAddress dcrutil.Address,
	amount dcrutil.Amount, minConfirms int,
	comment, commentTo string) (*chainhash.Hash, error) {

	return c.SendFromCommentAsync(fromAccount, toAddress, amount,
		minConfirms, comment, commentTo).Receive()
}

// FutureSendManyResult is a future promise to deliver the result of a
// SendManyAsync, SendManyMinConfAsync, or SendManyCommentAsync RPC invocation
// (or an applicable error).
type FutureSendManyResult chan *response

// Receive waits for the response promised by the future and returns the hash
// of the transaction sending multiple amounts to multiple addresses using the
// provided account as a source of funds.
func (r FutureSendManyResult) Receive() (*chainhash.Hash, error) {
	res, err := receiveFuture(r)
	if err != nil {
		return nil, err
	}

	// Unmarshal result as a string.
	var txHash string
	err = json.Unmarshal(res, &txHash)
	if err != nil {
		return nil, err
	}

	return chainhash.NewHashFromStr(txHash)
}

// SendManyAsync returns an instance of a type that can be used to get the
// result of the RPC at some future time by invoking the Receive function on the
// returned instance.
//
// See SendMany for the blocking version and more details.
func (c *Client) SendManyAsync(fromAccount string, amounts map[dcrutil.Address]dcrutil.Amount) FutureSendManyResult {
	convertedAmounts := make(map[string]float64, len(amounts))
	for addr, amount := range amounts {
		convertedAmounts[addr.Address()] = amount.ToCoin()
	}
	cmd := walletjson.NewSendManyCmd(fromAccount, convertedAmounts, nil, nil)
	return c.sendCmd(cmd)
}

// SendMany sends multiple amounts to multiple addresses using the provided
// account as a source of funds in a single transaction.  Only funds with the
// default number of minimum confirmations will be used.
//
// See SendManyMinConf and SendManyComment for different options.
//
// NOTE: This function requires to the wallet to be unlocked.  See the
// WalletPassphrase function for more details.
func (c *Client) SendMany(fromAccount string, amounts map[dcrutil.Address]dcrutil.Amount) (*chainhash.Hash, error) {
	return c.SendManyAsync(fromAccount, amounts).Receive()
}

// SendManyMinConfAsync returns an instance of a type that can be used to get
// the result of the RPC at some future time by invoking the Receive function on
// the returned instance.
//
// See SendManyMinConf for the blocking version and more details.
func (c *Client) SendManyMinConfAsync(fromAccount string,
	amounts map[dcrutil.Address]dcrutil.Amount,
	minConfirms int) FutureSendManyResult {

	convertedAmounts := make(map[string]float64, len(amounts))
	for addr, amount := range amounts {
		convertedAmounts[addr.Address()] = amount.ToCoin()
	}
	cmd := walletjson.NewSendManyCmd(fromAccount, convertedAmounts,
		&minConfirms, nil)
	return c.sendCmd(cmd)
}

// SendManyMinConf sends multiple amounts to multiple addresses using the
// provided account as a source of funds in a single transaction.  Only funds
// with the passed number of minimum confirmations will be used.
//
// See SendMany to use the default number of minimum confirmations and
// SendManyComment for additional options.
//
// NOTE: This function requires to the wallet to be unlocked.  See the
// WalletPassphrase function for more details.
func (c *Client) SendManyMinConf(fromAccount string,
	amounts map[dcrutil.Address]dcrutil.Amount,
	minConfirms int) (*chainhash.Hash, error) {

	return c.SendManyMinConfAsync(fromAccount, amounts, minConfirms).Receive()
}

// SendManyCommentAsync returns an instance of a type that can be used to get
// the result of the RPC at some future time by invoking the Receive function on
// the returned instance.
//
// See SendManyComment for the blocking version and more details.
func (c *Client) SendManyCommentAsync(fromAccount string,
	amounts map[dcrutil.Address]dcrutil.Amount, minConfirms int,
	comment string) FutureSendManyResult {

	convertedAmounts := make(map[string]float64, len(amounts))
	for addr, amount := range amounts {
		convertedAmounts[addr.Address()] = amount.ToCoin()
	}
	cmd := walletjson.NewSendManyCmd(fromAccount, convertedAmounts,
		&minConfirms, &comment)
	return c.sendCmd(cmd)
}

// SendManyComment sends multiple amounts to multiple addresses using the
// provided account as a source of funds in a single transaction and stores the
// provided comment in the wallet.  The comment parameter is intended to be used
// for the purpose of the transaction   Only funds with the passed number of
// minimum confirmations will be used.
//
// See SendMany and SendManyMinConf to use defaults.
//
// NOTE: This function requires to the wallet to be unlocked.  See the
// WalletPassphrase function for more details.
func (c *Client) SendManyComment(fromAccount string,
	amounts map[dcrutil.Address]dcrutil.Amount, minConfirms int,
	comment string) (*chainhash.Hash, error) {

	return c.SendManyCommentAsync(fromAccount, amounts, minConfirms,
		comment).Receive()
}

// Begin DECRED FUNCTIONS ---------------------------------------------------------
//
// SStx generation RPC call handling

// FuturePurchaseTicketResult a channel for the response promised by the future.
type FuturePurchaseTicketResult chan *response

// Receive waits for the response promised by the future and returns the hash
// of the transaction sending multiple amounts to multiple addresses using the
// provided account as a source of funds.
func (r FuturePurchaseTicketResult) Receive() ([]*chainhash.Hash, error) {
	res, err := receiveFuture(r)
	if err != nil {
		return nil, err
	}

	// Unmarshal result as a string slice.
	var txHashesStr []string
	err = json.Unmarshal(res, &txHashesStr)
	if err != nil {
		return nil, err
	}

	txHashes := make([]*chainhash.Hash, len(txHashesStr))
	for i := range txHashesStr {
		h, err := chainhash.NewHashFromStr(txHashesStr[i])
		if err != nil {
			return nil, err
		}
		txHashes[i] = h
	}

	return txHashes, nil
}

// PurchaseTicketAsync takes an account and a spending limit and returns a
// future chan to receive the transactions representing the purchased tickets
// when they come.
func (c *Client) PurchaseTicketAsync(fromAccount string,
	spendLimit dcrutil.Amount, minConf *int, ticketAddress dcrutil.Address,
	numTickets *int, poolAddress dcrutil.Address, poolFees *dcrutil.Amount,
	expiry *int, ticketFee *dcrutil.Amount) FuturePurchaseTicketResult {
	// An empty string is used to keep the sendCmd
	// passing of the command from accidentally
	// removing certain fields. We fill in the
	// default values of the other arguments as
	// well for the same reason.

	minConfVal := 1
	if minConf != nil {
		minConfVal = *minConf
	}

	ticketAddrStr := ""
	if ticketAddress != nil {
		ticketAddrStr = ticketAddress.Address()
	}

	numTicketsVal := 1
	if numTickets != nil {
		numTicketsVal = *numTickets
	}

	poolAddrStr := ""
	if poolAddress != nil {
		poolAddrStr = poolAddress.Address()
	}

	poolFeesFloat := 0.0
	if poolFees != nil {
		poolFeesFloat = poolFees.ToCoin()
	}

	expiryVal := 0
	if expiry != nil {
		expiryVal = *expiry
	}

	ticketFeeFloat := 0.0
	if ticketFee != nil {
		ticketFeeFloat = ticketFee.ToCoin()
	}

	cmd := walletjson.NewPurchaseTicketCmd(fromAccount, spendLimit.ToCoin(),
		&minConfVal, &ticketAddrStr, &numTicketsVal, &poolAddrStr,
		&poolFeesFloat, &expiryVal, dcrjson.String(""), &ticketFeeFloat)

	return c.sendCmd(cmd)
}

// PurchaseTicket takes an account and a spending limit and calls the async
// purchasetickets command.
func (c *Client) PurchaseTicket(fromAccount string,
	spendLimit dcrutil.Amount, minConf *int, ticketAddress dcrutil.Address,
	numTickets *int, poolAddress dcrutil.Address, poolFees *dcrutil.Amount,
	expiry *int, ticketChange *bool, ticketFee *dcrutil.Amount) ([]*chainhash.Hash, error) {

	return c.PurchaseTicketAsync(fromAccount, spendLimit, minConf, ticketAddress,
		numTickets, poolAddress, poolFees, expiry, ticketFee).Receive()
}

// END DECRED FUNCTIONS -----------------------------------------------------------

// *************************
// Address/Account Functions
// *************************

// FutureAddMultisigAddressResult is a future promise to deliver the result of a
// AddMultisigAddressAsync RPC invocation (or an applicable error).
type FutureAddMultisigAddressResult chan *response

// Receive waits for the response promised by the future and returns the
// multisignature address that requires the specified number of signatures for
// the provided addresses.
func (r FutureAddMultisigAddressResult) Receive(net dcrutil.AddressParams) (dcrutil.Address, error) {
	res, err := receiveFuture(r)
	if err != nil {
		return nil, err
	}

	// Unmarshal result as a string.
	var addr string
	err = json.Unmarshal(res, &addr)
	if err != nil {
		return nil, err
	}

	return dcrutil.DecodeAddress(addr, net)
}

// AddMultisigAddressAsync returns an instance of a type that can be used to get
// the result of the RPC at some future time by invoking the Receive function on
// the returned instance.
//
// See AddMultisigAddress for the blocking version and more details.
func (c *Client) AddMultisigAddressAsync(requiredSigs int, addresses []dcrutil.Address, account string) FutureAddMultisigAddressResult {
	addrs := make([]string, 0, len(addresses))
	for _, addr := range addresses {
		addrs = append(addrs, addr.String())
	}

	cmd := walletjson.NewAddMultisigAddressCmd(requiredSigs, addrs, &account)
	return c.sendCmd(cmd)
}

// AddMultisigAddress adds a multisignature address that requires the specified
// number of signatures for the provided addresses to the wallet.
func (c *Client) AddMultisigAddress(requiredSigs int, addresses []dcrutil.Address, account string, net dcrutil.AddressParams) (dcrutil.Address, error) {
	return c.AddMultisigAddressAsync(requiredSigs, addresses,
		account).Receive(net)
}

// FutureCreateMultisigResult is a future promise to deliver the result of a
// CreateMultisigAsync RPC invocation (or an applicable error).
type FutureCreateMultisigResult chan *response

// Receive waits for the response promised by the future and returns the
// multisignature address and script needed to redeem it.
func (r FutureCreateMultisigResult) Receive() (*walletjson.CreateMultiSigResult, error) {
	res, err := receiveFuture(r)
	if err != nil {
		return nil, err
	}

	// Unmarshal result as a createmultisig result object.
	var multisigRes walletjson.CreateMultiSigResult
	err = json.Unmarshal(res, &multisigRes)
	if err != nil {
		return nil, err
	}

	return &multisigRes, nil
}

// CreateMultisigAsync returns an instance of a type that can be used to get
// the result of the RPC at some future time by invoking the Receive function on
// the returned instance.
//
// See CreateMultisig for the blocking version and more details.
func (c *Client) CreateMultisigAsync(requiredSigs int, addresses []dcrutil.Address) FutureCreateMultisigResult {
	addrs := make([]string, 0, len(addresses))
	for _, addr := range addresses {
		addrs = append(addrs, addr.String())
	}

	cmd := walletjson.NewCreateMultisigCmd(requiredSigs, addrs)
	return c.sendCmd(cmd)
}

// CreateMultisig creates a multisignature address that requires the specified
// number of signatures for the provided addresses and returns the
// multisignature address and script needed to redeem it.
func (c *Client) CreateMultisig(requiredSigs int, addresses []dcrutil.Address) (*walletjson.CreateMultiSigResult, error) {
	return c.CreateMultisigAsync(requiredSigs, addresses).Receive()
}

// FutureCreateNewAccountResult is a future promise to deliver the result of a
// CreateNewAccountAsync RPC invocation (or an applicable error).
type FutureCreateNewAccountResult chan *response

// Receive waits for the response promised by the future and returns the
// result of creating new account.
func (r FutureCreateNewAccountResult) Receive() error {
	_, err := receiveFuture(r)
	return err
}

// CreateNewAccountAsync returns an instance of a type that can be used to get the
// result of the RPC at some future time by invoking the Receive function on the
// returned instance.
//
// See CreateNewAccount for the blocking version and more details.
func (c *Client) CreateNewAccountAsync(account string) FutureCreateNewAccountResult {
	cmd := walletjson.NewCreateNewAccountCmd(account)
	return c.sendCmd(cmd)
}

// CreateNewAccount creates a new wallet account.
func (c *Client) CreateNewAccount(account string) error {
	return c.CreateNewAccountAsync(account).Receive()
}

// FutureGetNewAddressResult is a future promise to deliver the result of a
// GetNewAddressAsync RPC invocation (or an applicable error).
type FutureGetNewAddressResult chan *response

// Receive waits for the response promised by the future and returns a new
// address.
func (r FutureGetNewAddressResult) Receive(net dcrutil.AddressParams) (dcrutil.Address, error) {
	res, err := receiveFuture(r)
	if err != nil {
		return nil, err
	}

	// Unmarshal result as a string.
	var addr string
	err = json.Unmarshal(res, &addr)
	if err != nil {
		return nil, err
	}

	return dcrutil.DecodeAddress(addr, net)
}

// GapPolicy defines the policy to use when the BIP0044 unused address gap limit
// would be violated by creating a new address.
type GapPolicy string

// Gap policies that are understood by a wallet JSON-RPC server.  These are
// defined for safety and convenience, but string literals can be used as well.
const (
	GapPolicyError  GapPolicy = "error"
	GapPolicyIgnore GapPolicy = "ignore"
	GapPolicyWrap   GapPolicy = "wrap"
)

// GetNewAddressAsync returns an instance of a type that can be used to get the
// result of the RPC at some future time by invoking the Receive function on the
// returned instance.
//
// See GetNewAddress for the blocking version and more details.
func (c *Client) GetNewAddressAsync(account string) FutureGetNewAddressResult {
	cmd := walletjson.NewGetNewAddressCmd(&account, nil)
	return c.sendCmd(cmd)
}

// GetNewAddressGapPolicyAsync returns an instance of a type that can be used to
// get the result of the RPC at some future time by invoking the Receive
// function on the returned instance.
//
// See GetNewAddressGapPolicy for the blocking version and more details.
func (c *Client) GetNewAddressGapPolicyAsync(account string, gapPolicy GapPolicy) FutureGetNewAddressResult {
	cmd := walletjson.NewGetNewAddressCmd(&account, (*string)(&gapPolicy))
	return c.sendCmd(cmd)
}

// GetNewAddress returns a new address.
func (c *Client) GetNewAddress(account string, net dcrutil.AddressParams) (dcrutil.Address, error) {
	return c.GetNewAddressAsync(account).Receive(net)
}

// GetNewAddressGapPolicy returns a new address while allowing callers to
// control the BIP0044 unused address gap limit policy.
func (c *Client) GetNewAddressGapPolicy(account string, gapPolicy GapPolicy, net dcrutil.AddressParams) (dcrutil.Address, error) {
	return c.GetNewAddressGapPolicyAsync(account, gapPolicy).Receive(net)
}

// FutureGetRawChangeAddressResult is a future promise to deliver the result of
// a GetRawChangeAddressAsync RPC invocation (or an applicable error).
type FutureGetRawChangeAddressResult chan *response

// Receive waits for the response promised by the future and returns a new
// address for receiving change that will be associated with the provided
// account.  Note that this is only for raw transactions and NOT for normal use.
func (r FutureGetRawChangeAddressResult) Receive(net dcrutil.AddressParams) (dcrutil.Address, error) {
	res, err := receiveFuture(r)
	if err != nil {
		return nil, err
	}

	// Unmarshal result as a string.
	var addr string
	err = json.Unmarshal(res, &addr)
	if err != nil {
		return nil, err
	}

	return dcrutil.DecodeAddress(addr, net)
}

// GetRawChangeAddressAsync returns an instance of a type that can be used to
// get the result of the RPC at some future time by invoking the Receive
// function on the returned instance.
//
// See GetRawChangeAddress for the blocking version and more details.
func (c *Client) GetRawChangeAddressAsync(account string) FutureGetRawChangeAddressResult {
	cmd := walletjson.NewGetRawChangeAddressCmd(&account)
	return c.sendCmd(cmd)
}

// GetRawChangeAddress returns a new address for receiving change that will be
// associated with the provided account.  Note that this is only for raw
// transactions and NOT for normal use.
func (c *Client) GetRawChangeAddress(account string, net dcrutil.AddressParams) (dcrutil.Address, error) {
	return c.GetRawChangeAddressAsync(account).Receive(net)
}

// FutureGetAccountAddressResult is a future promise to deliver the result of a
// GetAccountAddressAsync RPC invocation (or an applicable error).
type FutureGetAccountAddressResult chan *response

// Receive waits for the response promised by the future and returns the current
// Decred address for receiving payments to the specified account.
func (r FutureGetAccountAddressResult) Receive(net dcrutil.AddressParams) (dcrutil.Address, error) {
	res, err := receiveFuture(r)
	if err != nil {
		return nil, err
	}

	// Unmarshal result as a string.
	var addr string
	err = json.Unmarshal(res, &addr)
	if err != nil {
		return nil, err
	}

	return dcrutil.DecodeAddress(addr, net)
}

// GetAccountAddressAsync returns an instance of a type that can be used to get
// the result of the RPC at some future time by invoking the Receive function on
// the returned instance.
//
// See GetAccountAddress for the blocking version and more details.
func (c *Client) GetAccountAddressAsync(account string) FutureGetAccountAddressResult {
	cmd := walletjson.NewGetAccountAddressCmd(account)
	return c.sendCmd(cmd)
}

// GetAccountAddress returns the current Decred address for receiving payments
// to the specified account.
func (c *Client) GetAccountAddress(account string, net dcrutil.AddressParams) (dcrutil.Address, error) {
	return c.GetAccountAddressAsync(account).Receive(net)
}

// FutureGetAccountResult is a future promise to deliver the result of a
// GetAccountAsync RPC invocation (or an applicable error).
type FutureGetAccountResult chan *response

// Receive waits for the response promised by the future and returns the account
// associated with the passed address.
func (r FutureGetAccountResult) Receive() (string, error) {
	res, err := receiveFuture(r)
	if err != nil {
		return "", err
	}

	// Unmarshal result as a string.
	var account string
	err = json.Unmarshal(res, &account)
	if err != nil {
		return "", err
	}

	return account, nil
}

// GetAccountAsync returns an instance of a type that can be used to get the
// result of the RPC at some future time by invoking the Receive function on the
// returned instance.
//
// See GetAccount for the blocking version and more details.
func (c *Client) GetAccountAsync(address dcrutil.Address) FutureGetAccountResult {
	addr := address.Address()
	cmd := walletjson.NewGetAccountCmd(addr)
	return c.sendCmd(cmd)
}

// GetAccount returns the account associated with the passed address.
func (c *Client) GetAccount(address dcrutil.Address) (string, error) {
	return c.GetAccountAsync(address).Receive()
}

// FutureGetAddressesByAccountResult is a future promise to deliver the result
// of a GetAddressesByAccountAsync RPC invocation (or an applicable error).
type FutureGetAddressesByAccountResult chan *response

// Receive waits for the response promised by the future and returns the list of
// addresses associated with the passed account.
func (r FutureGetAddressesByAccountResult) Receive(net dcrutil.AddressParams) ([]dcrutil.Address, error) {
	res, err := receiveFuture(r)
	if err != nil {
		return nil, err
	}

	// Unmarshal result as an array of string.
	var addrStrings []string
	err = json.Unmarshal(res, &addrStrings)
	if err != nil {
		return nil, err
	}

	addrs := make([]dcrutil.Address, 0, len(addrStrings))
	for _, addrStr := range addrStrings {
		addr, err := dcrutil.DecodeAddress(addrStr, net)
		if err != nil {
			return nil, err
		}
		addrs = append(addrs, addr)
	}

	return addrs, nil
}

// GetAddressesByAccountAsync returns an instance of a type that can be used to
// get the result of the RPC at some future time by invoking the Receive
// function on the returned instance.
//
// See GetAddressesByAccount for the blocking version and more details.
func (c *Client) GetAddressesByAccountAsync(account string) FutureGetAddressesByAccountResult {
	cmd := walletjson.NewGetAddressesByAccountCmd(account)
	return c.sendCmd(cmd)
}

// GetAddressesByAccount returns the list of addresses associated with the
// passed account.
func (c *Client) GetAddressesByAccount(account string, net dcrutil.AddressParams) ([]dcrutil.Address, error) {
	return c.GetAddressesByAccountAsync(account).Receive(net)
}

// FutureMoveResult is a future promise to deliver the result of a MoveAsync,
// MoveMinConfAsync, or MoveCommentAsync RPC invocation (or an applicable
// error).
type FutureMoveResult chan *response

// Receive waits for the response promised by the future and returns the result
// of the move operation.
func (r FutureMoveResult) Receive() (bool, error) {
	res, err := receiveFuture(r)
	if err != nil {
		return false, err
	}

	// Unmarshal result as a boolean.
	var moveResult bool
	err = json.Unmarshal(res, &moveResult)
	if err != nil {
		return false, err
	}

	return moveResult, nil
}

// FutureRenameAccountResult is a future promise to deliver the result of a
// RenameAccountAsync RPC invocation (or an applicable error).
type FutureRenameAccountResult chan *response

// Receive waits for the response promised by the future and returns the
// result of creating new account.
func (r FutureRenameAccountResult) Receive() error {
	_, err := receiveFuture(r)
	return err
}

// RenameAccountAsync returns an instance of a type that can be used to get the
// result of the RPC at some future time by invoking the Receive function on the
// returned instance.
//
// See RenameAccount for the blocking version and more details.
func (c *Client) RenameAccountAsync(oldAccount, newAccount string) FutureRenameAccountResult {
	cmd := walletjson.NewRenameAccountCmd(oldAccount, newAccount)
	return c.sendCmd(cmd)
}

// RenameAccount creates a new wallet account.
func (c *Client) RenameAccount(oldAccount, newAccount string) error {
	return c.RenameAccountAsync(oldAccount, newAccount).Receive()
}

// FutureValidateAddressResult is a future promise to deliver the result of a
// ValidateAddressAsync RPC invocation (or an applicable error).
type FutureValidateAddressResult chan *response

// Receive waits for the response promised by the future and returns information
// about the given Decred address.
func (r FutureValidateAddressResult) Receive() (*walletjson.ValidateAddressWalletResult, error) {
	res, err := receiveFuture(r)
	if err != nil {
		return nil, err
	}

	// Unmarshal result as a validateaddress result object.
	var addrResult walletjson.ValidateAddressWalletResult
	err = json.Unmarshal(res, &addrResult)
	if err != nil {
		return nil, err
	}

	return &addrResult, nil
}

// ValidateAddressAsync returns an instance of a type that can be used to get
// the result of the RPC at some future time by invoking the Receive function on
// the returned instance.
//
// See ValidateAddress for the blocking version and more details.
func (c *Client) ValidateAddressAsync(address dcrutil.Address) FutureValidateAddressResult {
	addr := address.Address()
	cmd := chainjson.NewValidateAddressCmd(addr)
	return c.sendCmd(cmd)
}

// ValidateAddress returns information about the given Decred address.
func (c *Client) ValidateAddress(address dcrutil.Address) (*walletjson.ValidateAddressWalletResult, error) {
	return c.ValidateAddressAsync(address).Receive()
}

// FutureKeyPoolRefillResult is a future promise to deliver the result of a
// KeyPoolRefillAsync RPC invocation (or an applicable error).
type FutureKeyPoolRefillResult chan *response

// Receive waits for the response promised by the future and returns the result
// of refilling the key pool.
func (r FutureKeyPoolRefillResult) Receive() error {
	_, err := receiveFuture(r)
	return err
}

// KeyPoolRefillAsync returns an instance of a type that can be used to get the
// result of the RPC at some future time by invoking the Receive function on the
// returned instance.
//
// See KeyPoolRefill for the blocking version and more details.
func (c *Client) KeyPoolRefillAsync() FutureKeyPoolRefillResult {
	cmd := walletjson.NewKeyPoolRefillCmd(nil)
	return c.sendCmd(cmd)
}

// KeyPoolRefill fills the key pool as necessary to reach the default size.
//
// See KeyPoolRefillSize to override the size of the key pool.
func (c *Client) KeyPoolRefill() error {
	return c.KeyPoolRefillAsync().Receive()
}

// KeyPoolRefillSizeAsync returns an instance of a type that can be used to get
// the result of the RPC at some future time by invoking the Receive function on
// the returned instance.
//
// See KeyPoolRefillSize for the blocking version and more details.
func (c *Client) KeyPoolRefillSizeAsync(newSize uint) FutureKeyPoolRefillResult {
	cmd := walletjson.NewKeyPoolRefillCmd(&newSize)
	return c.sendCmd(cmd)
}

// KeyPoolRefillSize fills the key pool as necessary to reach the specified
// size.
func (c *Client) KeyPoolRefillSize(newSize uint) error {
	return c.KeyPoolRefillSizeAsync(newSize).Receive()
}

// ************************
// Amount/Balance Functions
// ************************

// FutureListAccountsResult is a future promise to deliver the result of a
// ListAccountsAsync or ListAccountsMinConfAsync RPC invocation (or an
// applicable error).
type FutureListAccountsResult chan *response

// Receive waits for the response promised by the future and returns a
// map of account names and their associated balances.
func (r FutureListAccountsResult) Receive() (map[string]dcrutil.Amount, error) {
	res, err := receiveFuture(r)
	if err != nil {
		return nil, err
	}

	// Unmarshal result as a json object.
	var accounts map[string]float64
	err = json.Unmarshal(res, &accounts)
	if err != nil {
		return nil, err
	}

	accountsMap := make(map[string]dcrutil.Amount)
	for k, v := range accounts {
		amount, err := dcrutil.NewAmount(v)
		if err != nil {
			return nil, err
		}

		accountsMap[k] = amount
	}

	return accountsMap, nil
}

// ListAccountsAsync returns an instance of a type that can be used to get the
// result of the RPC at some future time by invoking the Receive function on the
// returned instance.
//
// See ListAccounts for the blocking version and more details.
func (c *Client) ListAccountsAsync() FutureListAccountsResult {
	cmd := walletjson.NewListAccountsCmd(nil)
	return c.sendCmd(cmd)
}

// ListAccounts returns a map of account names and their associated balances
// using the default number of minimum confirmations.
//
// See ListAccountsMinConf to override the minimum number of confirmations.
func (c *Client) ListAccounts() (map[string]dcrutil.Amount, error) {
	return c.ListAccountsAsync().Receive()
}

// ListAccountsMinConfAsync returns an instance of a type that can be used to
// get the result of the RPC at some future time by invoking the Receive
// function on the returned instance.
//
// See ListAccountsMinConf for the blocking version and more details.
func (c *Client) ListAccountsMinConfAsync(minConfirms int) FutureListAccountsResult {
	cmd := walletjson.NewListAccountsCmd(&minConfirms)
	return c.sendCmd(cmd)
}

// ListAccountsMinConf returns a map of account names and their associated
// balances using the specified number of minimum confirmations.
//
// See ListAccounts to use the default minimum number of confirmations.
func (c *Client) ListAccountsMinConf(minConfirms int) (map[string]dcrutil.Amount, error) {
	return c.ListAccountsMinConfAsync(minConfirms).Receive()
}

// FutureGetMasterPubkeyResult is a future promise to deliver the result of a
// GetMasterPubkeyAsync RPC invocation (or an applicable error).
type FutureGetMasterPubkeyResult chan *response

// Receive waits for the response promised by the future and returns a pointer to
// the master extended public key for account and the network's hierarchical
// deterministic extended key magic versions (e.g. MainNetParams)
func (r FutureGetMasterPubkeyResult) Receive(net hdkeychain.NetworkParams) (*hdkeychain.ExtendedKey, error) {
	res, err := receiveFuture(r)
	if err != nil {
		return nil, err
	}

	var pubkeystr string
	err = json.Unmarshal(res, &pubkeystr)
	if err != nil {
		return nil, err
	}

	// pubkey is a hierarchical deterministic master public key
	pubkey, err := hdkeychain.NewKeyFromString(pubkeystr, net)
	if err != nil {
		return nil, err
	}

	return pubkey, nil
}

// GetMasterPubkeyAsync returns an instance of a type that can be used to get the
// result of the RPC at some future time by invoking the Receive function on the
// returned instance.
//
// See GetMasterPubkey for the blocking version and more details.
func (c *Client) GetMasterPubkeyAsync(account string) FutureGetMasterPubkeyResult {
	cmd := walletjson.NewGetMasterPubkeyCmd(&account)
	return c.sendCmd(cmd)
}

// GetMasterPubkey returns a pointer to the master extended public key for account
// and the network's hierarchical deterministic extended key magic versions
// (e.g. MainNetParams)
func (c *Client) GetMasterPubkey(account string, net hdkeychain.NetworkParams) (*hdkeychain.ExtendedKey, error) {
	return c.GetMasterPubkeyAsync(account).Receive(net)
}

// FutureGetBalanceResult is a future promise to deliver the result of a
// GetBalanceAsync or GetBalanceMinConfAsync RPC invocation (or an applicable
// error).
type FutureGetBalanceResult chan *response

// Receive waits for the response promised by the future and returns the
// available balance from the server for the specified account.
func (r FutureGetBalanceResult) Receive() (*walletjson.GetBalanceResult, error) {
	res, err := receiveFuture(r)
	if err != nil {
		return nil, err
	}

	// Unmarshal result as a floating point number.
	var balance walletjson.GetBalanceResult
	err = json.Unmarshal(res, &balance)
	if err != nil {
		return nil, err
	}

	return &balance, nil
}

// GetBalanceAsync returns an instance of a type that can be used to get the
// result of the RPC at some future time by invoking the Receive function on the
// returned instance.
//
// See GetBalance for the blocking version and more details.
func (c *Client) GetBalanceAsync(account string) FutureGetBalanceResult {
	cmd := walletjson.NewGetBalanceCmd(&account, nil)
	return c.sendCmd(cmd)
}

// GetBalance returns the available balance from the server for the specified
// account using the default number of minimum confirmations.  The account may
// be "*" for all accounts.
//
// See GetBalanceMinConf to override the minimum number of confirmations.
func (c *Client) GetBalance(account string) (*walletjson.GetBalanceResult, error) {
	return c.GetBalanceAsync(account).Receive()
}

// GetBalanceMinConfAsync returns an instance of a type that can be used to get
// the result of the RPC at some future time by invoking the Receive function on
// the returned instance.
//
// See GetBalanceMinConf for the blocking version and more details.
func (c *Client) GetBalanceMinConfAsync(account string, minConfirms int) FutureGetBalanceResult {
	cmd := walletjson.NewGetBalanceCmd(&account, &minConfirms)
	return c.sendCmd(cmd)
}

// GetBalanceMinConf returns the available balance from the server for the
// specified account using the specified number of minimum confirmations.  The
// account may be "*" for all accounts.
//
// See GetBalance to use the default minimum number of confirmations.
func (c *Client) GetBalanceMinConf(account string, minConfirms int) (*walletjson.GetBalanceResult, error) {
	return c.GetBalanceMinConfAsync(account, minConfirms).Receive()
}

// FutureGetReceivedByAccountResult is a future promise to deliver the result of
// a GetReceivedByAccountAsync or GetReceivedByAccountMinConfAsync RPC
// invocation (or an applicable error).
type FutureGetReceivedByAccountResult chan *response

// Receive waits for the response promised by the future and returns the total
// amount received with the specified account.
func (r FutureGetReceivedByAccountResult) Receive() (dcrutil.Amount, error) {
	res, err := receiveFuture(r)
	if err != nil {
		return 0, err
	}

	// Unmarshal result as a floating point number.
	var balance float64
	err = json.Unmarshal(res, &balance)
	if err != nil {
		return 0, err
	}

	amount, err := dcrutil.NewAmount(balance)
	if err != nil {
		return 0, err
	}

	return amount, nil
}

// GetReceivedByAccountAsync returns an instance of a type that can be used to
// get the result of the RPC at some future time by invoking the Receive
// function on the returned instance.
//
// See GetReceivedByAccount for the blocking version and more details.
func (c *Client) GetReceivedByAccountAsync(account string) FutureGetReceivedByAccountResult {
	cmd := walletjson.NewGetReceivedByAccountCmd(account, nil)
	return c.sendCmd(cmd)
}

// GetReceivedByAccount returns the total amount received with the specified
// account with at least the default number of minimum confirmations.
//
// See GetReceivedByAccountMinConf to override the minimum number of
// confirmations.
func (c *Client) GetReceivedByAccount(account string) (dcrutil.Amount, error) {
	return c.GetReceivedByAccountAsync(account).Receive()
}

// GetReceivedByAccountMinConfAsync returns an instance of a type that can be
// used to get the result of the RPC at some future time by invoking the Receive
// function on the returned instance.
//
// See GetReceivedByAccountMinConf for the blocking version and more details.
func (c *Client) GetReceivedByAccountMinConfAsync(account string, minConfirms int) FutureGetReceivedByAccountResult {
	cmd := walletjson.NewGetReceivedByAccountCmd(account, &minConfirms)
	return c.sendCmd(cmd)
}

// GetReceivedByAccountMinConf returns the total amount received with the
// specified account with at least the specified number of minimum
// confirmations.
//
// See GetReceivedByAccount to use the default minimum number of confirmations.
func (c *Client) GetReceivedByAccountMinConf(account string, minConfirms int) (dcrutil.Amount, error) {
	return c.GetReceivedByAccountMinConfAsync(account, minConfirms).Receive()
}

// FutureGetUnconfirmedBalanceResult is a future promise to deliver the result
// of a GetUnconfirmedBalanceAsync RPC invocation (or an applicable error).
type FutureGetUnconfirmedBalanceResult chan *response

// Receive waits for the response promised by the future and returns the
// unconfirmed balance from the server for the specified account.
func (r FutureGetUnconfirmedBalanceResult) Receive() (dcrutil.Amount, error) {
	res, err := receiveFuture(r)
	if err != nil {
		return 0, err
	}

	// Unmarshal result as a floating point number.
	var balance float64
	err = json.Unmarshal(res, &balance)
	if err != nil {
		return 0, err
	}

	amount, err := dcrutil.NewAmount(balance)
	if err != nil {
		return 0, err
	}

	return amount, nil
}

// GetUnconfirmedBalanceAsync returns an instance of a type that can be used to
// get the result of the RPC at some future time by invoking the Receive
// function on the returned instance.
//
// See GetUnconfirmedBalance for the blocking version and more details.
func (c *Client) GetUnconfirmedBalanceAsync(account string) FutureGetUnconfirmedBalanceResult {
	cmd := walletjson.NewGetUnconfirmedBalanceCmd(&account)
	return c.sendCmd(cmd)
}

// GetUnconfirmedBalance returns the unconfirmed balance from the server for
// the specified account.
func (c *Client) GetUnconfirmedBalance(account string) (dcrutil.Amount, error) {
	return c.GetUnconfirmedBalanceAsync(account).Receive()
}

// FutureGetReceivedByAddressResult is a future promise to deliver the result of
// a GetReceivedByAddressAsync or GetReceivedByAddressMinConfAsync RPC
// invocation (or an applicable error).
type FutureGetReceivedByAddressResult chan *response

// Receive waits for the response promised by the future and returns the total
// amount received by the specified address.
func (r FutureGetReceivedByAddressResult) Receive() (dcrutil.Amount, error) {
	res, err := receiveFuture(r)
	if err != nil {
		return 0, err
	}

	// Unmarshal result as a floating point number.
	var balance float64
	err = json.Unmarshal(res, &balance)
	if err != nil {
		return 0, err
	}

	amount, err := dcrutil.NewAmount(balance)
	if err != nil {
		return 0, err
	}

	return amount, nil
}

// GetReceivedByAddressAsync returns an instance of a type that can be used to
// get the result of the RPC at some future time by invoking the Receive
// function on the returned instance.
//
// See GetReceivedByAddress for the blocking version and more details.
func (c *Client) GetReceivedByAddressAsync(address dcrutil.Address) FutureGetReceivedByAddressResult {
	addr := address.Address()
	cmd := walletjson.NewGetReceivedByAddressCmd(addr, nil)
	return c.sendCmd(cmd)

}

// GetReceivedByAddress returns the total amount received by the specified
// address with at least the default number of minimum confirmations.
//
// See GetReceivedByAddressMinConf to override the minimum number of
// confirmations.
func (c *Client) GetReceivedByAddress(address dcrutil.Address) (dcrutil.Amount, error) {
	return c.GetReceivedByAddressAsync(address).Receive()
}

// GetReceivedByAddressMinConfAsync returns an instance of a type that can be
// used to get the result of the RPC at some future time by invoking the Receive
// function on the returned instance.
//
// See GetReceivedByAddressMinConf for the blocking version and more details.
func (c *Client) GetReceivedByAddressMinConfAsync(address dcrutil.Address, minConfirms int) FutureGetReceivedByAddressResult {
	addr := address.Address()
	cmd := walletjson.NewGetReceivedByAddressCmd(addr, &minConfirms)
	return c.sendCmd(cmd)
}

// GetReceivedByAddressMinConf returns the total amount received by the specified
// address with at least the specified number of minimum confirmations.
//
// See GetReceivedByAddress to use the default minimum number of confirmations.
func (c *Client) GetReceivedByAddressMinConf(address dcrutil.Address, minConfirms int) (dcrutil.Amount, error) {
	return c.GetReceivedByAddressMinConfAsync(address, minConfirms).Receive()
}

// FutureListReceivedByAccountResult is a future promise to deliver the result
// of a ListReceivedByAccountAsync, ListReceivedByAccountMinConfAsync, or
// ListReceivedByAccountIncludeEmptyAsync RPC invocation (or an applicable
// error).
type FutureListReceivedByAccountResult chan *response

// Receive waits for the response promised by the future and returns a list of
// balances by account.
func (r FutureListReceivedByAccountResult) Receive() ([]walletjson.ListReceivedByAccountResult, error) {
	res, err := receiveFuture(r)
	if err != nil {
		return nil, err
	}

	// Unmarshal as an array of listreceivedbyaccount result objects.
	var received []walletjson.ListReceivedByAccountResult
	err = json.Unmarshal(res, &received)
	if err != nil {
		return nil, err
	}

	return received, nil
}

// ListReceivedByAccountAsync returns an instance of a type that can be used to
// get the result of the RPC at some future time by invoking the Receive
// function on the returned instance.
//
// See ListReceivedByAccount for the blocking version and more details.
func (c *Client) ListReceivedByAccountAsync() FutureListReceivedByAccountResult {
	cmd := walletjson.NewListReceivedByAccountCmd(nil, nil, nil)
	return c.sendCmd(cmd)
}

// ListReceivedByAccount lists balances by account using the default number
// of minimum confirmations and including accounts that haven't received any
// payments.
//
// See ListReceivedByAccountMinConf to override the minimum number of
// confirmations and ListReceivedByAccountIncludeEmpty to filter accounts that
// haven't received any payments from the results.
func (c *Client) ListReceivedByAccount() ([]walletjson.ListReceivedByAccountResult, error) {
	return c.ListReceivedByAccountAsync().Receive()
}

// ListReceivedByAccountMinConfAsync returns an instance of a type that can be
// used to get the result of the RPC at some future time by invoking the Receive
// function on the returned instance.
//
// See ListReceivedByAccountMinConf for the blocking version and more details.
func (c *Client) ListReceivedByAccountMinConfAsync(minConfirms int) FutureListReceivedByAccountResult {
	cmd := walletjson.NewListReceivedByAccountCmd(&minConfirms, nil, nil)
	return c.sendCmd(cmd)
}

// ListReceivedByAccountMinConf lists balances by account using the specified
// number of minimum confirmations not including accounts that haven't received
// any payments.
//
// See ListReceivedByAccount to use the default minimum number of confirmations
// and ListReceivedByAccountIncludeEmpty to also filter accounts that haven't
// received any payments from the results.
func (c *Client) ListReceivedByAccountMinConf(minConfirms int) ([]walletjson.ListReceivedByAccountResult, error) {
	return c.ListReceivedByAccountMinConfAsync(minConfirms).Receive()
}

// ListReceivedByAccountIncludeEmptyAsync returns an instance of a type that can
// be used to get the result of the RPC at some future time by invoking the
// Receive function on the returned instance.
//
// See ListReceivedByAccountIncludeEmpty for the blocking version and more details.
func (c *Client) ListReceivedByAccountIncludeEmptyAsync(minConfirms int, includeEmpty bool) FutureListReceivedByAccountResult {
	cmd := walletjson.NewListReceivedByAccountCmd(&minConfirms, &includeEmpty,
		nil)
	return c.sendCmd(cmd)
}

// ListReceivedByAccountIncludeEmpty lists balances by account using the
// specified number of minimum confirmations and including accounts that
// haven't received any payments depending on specified flag.
//
// See ListReceivedByAccount and ListReceivedByAccountMinConf to use defaults.
func (c *Client) ListReceivedByAccountIncludeEmpty(minConfirms int, includeEmpty bool) ([]walletjson.ListReceivedByAccountResult, error) {
	return c.ListReceivedByAccountIncludeEmptyAsync(minConfirms,
		includeEmpty).Receive()
}

// FutureListReceivedByAddressResult is a future promise to deliver the result
// of a ListReceivedByAddressAsync, ListReceivedByAddressMinConfAsync, or
// ListReceivedByAddressIncludeEmptyAsync RPC invocation (or an applicable
// error).
type FutureListReceivedByAddressResult chan *response

// Receive waits for the response promised by the future and returns a list of
// balances by address.
func (r FutureListReceivedByAddressResult) Receive() ([]walletjson.ListReceivedByAddressResult, error) {
	res, err := receiveFuture(r)
	if err != nil {
		return nil, err
	}

	// Unmarshal as an array of listreceivedbyaddress result objects.
	var received []walletjson.ListReceivedByAddressResult
	err = json.Unmarshal(res, &received)
	if err != nil {
		return nil, err
	}

	return received, nil
}

// ListReceivedByAddressAsync returns an instance of a type that can be used to
// get the result of the RPC at some future time by invoking the Receive
// function on the returned instance.
//
// See ListReceivedByAddress for the blocking version and more details.
func (c *Client) ListReceivedByAddressAsync() FutureListReceivedByAddressResult {
	cmd := walletjson.NewListReceivedByAddressCmd(nil, nil, nil)
	return c.sendCmd(cmd)
}

// ListReceivedByAddress lists balances by address using the default number
// of minimum confirmations not including addresses that haven't received any
// payments or watching only addresses.
//
// See ListReceivedByAddressMinConf to override the minimum number of
// confirmations and ListReceivedByAddressIncludeEmpty to also include addresses
// that haven't received any payments in the results.
func (c *Client) ListReceivedByAddress() ([]walletjson.ListReceivedByAddressResult, error) {
	return c.ListReceivedByAddressAsync().Receive()
}

// ListReceivedByAddressMinConfAsync returns an instance of a type that can be
// used to get the result of the RPC at some future time by invoking the Receive
// function on the returned instance.
//
// See ListReceivedByAddressMinConf for the blocking version and more details.
func (c *Client) ListReceivedByAddressMinConfAsync(minConfirms int) FutureListReceivedByAddressResult {
	cmd := walletjson.NewListReceivedByAddressCmd(&minConfirms, nil, nil)
	return c.sendCmd(cmd)
}

// ListReceivedByAddressMinConf lists balances by address using the specified
// number of minimum confirmations not including addresses that haven't received
// any payments.
//
// See ListReceivedByAddress to use the default minimum number of confirmations
// and ListReceivedByAddressIncludeEmpty to also include addresses that haven't
// received any payments in the results.
func (c *Client) ListReceivedByAddressMinConf(minConfirms int) ([]walletjson.ListReceivedByAddressResult, error) {
	return c.ListReceivedByAddressMinConfAsync(minConfirms).Receive()
}

// ListReceivedByAddressIncludeEmptyAsync returns an instance of a type that can
// be used to get the result of the RPC at some future time by invoking the
// Receive function on the returned instance.
//
// See ListReceivedByAccountIncludeEmpty for the blocking version and more details.
func (c *Client) ListReceivedByAddressIncludeEmptyAsync(minConfirms int, includeEmpty bool) FutureListReceivedByAddressResult {
	cmd := walletjson.NewListReceivedByAddressCmd(&minConfirms, &includeEmpty,
		nil)
	return c.sendCmd(cmd)
}

// ListReceivedByAddressIncludeEmpty lists balances by address using the
// specified number of minimum confirmations and including addresses that
// haven't received any payments depending on specified flag.
//
// See ListReceivedByAddress and ListReceivedByAddressMinConf to use defaults.
func (c *Client) ListReceivedByAddressIncludeEmpty(minConfirms int, includeEmpty bool) ([]walletjson.ListReceivedByAddressResult, error) {
	return c.ListReceivedByAddressIncludeEmptyAsync(minConfirms,
		includeEmpty).Receive()
}

// ************************
// Wallet Locking Functions
// ************************

// FutureWalletLockResult is a future promise to deliver the result of a
// WalletLockAsync RPC invocation (or an applicable error).
type FutureWalletLockResult chan *response

// Receive waits for the response promised by the future and returns the result
// of locking the wallet.
func (r FutureWalletLockResult) Receive() error {
	_, err := receiveFuture(r)
	return err
}

// WalletLockAsync returns an instance of a type that can be used to get the
// result of the RPC at some future time by invoking the Receive function on the
// returned instance.
//
// See WalletLock for the blocking version and more details.
func (c *Client) WalletLockAsync() FutureWalletLockResult {
	cmd := walletjson.NewWalletLockCmd()
	return c.sendCmd(cmd)
}

// WalletLock locks the wallet by removing the encryption key from memory.
//
// After calling this function, the WalletPassphrase function must be used to
// unlock the wallet prior to calling any other function which requires the
// wallet to be unlocked.
func (c *Client) WalletLock() error {
	return c.WalletLockAsync().Receive()
}

// WalletPassphrase unlocks the wallet by using the passphrase to derive the
// decryption key which is then stored in memory for the specified timeout
// (in seconds).
func (c *Client) WalletPassphrase(passphrase string, timeoutSecs int64) error {
	cmd := walletjson.NewWalletPassphraseCmd(passphrase, timeoutSecs)
	_, err := c.sendCmdAndWait(cmd)
	return err
}

// FutureWalletPassphraseChangeResult is a future promise to deliver the result
// of a WalletPassphraseChangeAsync RPC invocation (or an applicable error).
type FutureWalletPassphraseChangeResult chan *response

// Receive waits for the response promised by the future and returns the result
// of changing the wallet passphrase.
func (r FutureWalletPassphraseChangeResult) Receive() error {
	_, err := receiveFuture(r)
	return err
}

// WalletPassphraseChangeAsync returns an instance of a type that can be used to
// get the result of the RPC at some future time by invoking the Receive
// function on the returned instance.
//
// See WalletPassphraseChange for the blocking version and more details.
func (c *Client) WalletPassphraseChangeAsync(old, new string) FutureWalletPassphraseChangeResult {
	cmd := walletjson.NewWalletPassphraseChangeCmd(old, new)
	return c.sendCmd(cmd)
}

// WalletPassphraseChange changes the wallet passphrase from the specified old
// to new passphrase.
func (c *Client) WalletPassphraseChange(old, new string) error {
	return c.WalletPassphraseChangeAsync(old, new).Receive()
}

// *************************
// Message Signing Functions
// *************************

// FutureSignMessageResult is a future promise to deliver the result of a
// SignMessageAsync RPC invocation (or an applicable error).
type FutureSignMessageResult chan *response

// Receive waits for the response promised by the future and returns the message
// signed with the private key of the specified address.
func (r FutureSignMessageResult) Receive() (string, error) {
	res, err := receiveFuture(r)
	if err != nil {
		return "", err
	}

	// Unmarshal result as a string.
	var b64 string
	err = json.Unmarshal(res, &b64)
	if err != nil {
		return "", err
	}

	return b64, nil
}

// SignMessageAsync returns an instance of a type that can be used to get the
// result of the RPC at some future time by invoking the Receive function on the
// returned instance.
//
// See SignMessage for the blocking version and more details.
func (c *Client) SignMessageAsync(address dcrutil.Address, message string) FutureSignMessageResult {
	addr := address.Address()
	cmd := walletjson.NewSignMessageCmd(addr, message)
	return c.sendCmd(cmd)
}

// SignMessage signs a message with the private key of the specified address.
//
// NOTE: This function requires to the wallet to be unlocked.  See the
// WalletPassphrase function for more details.
func (c *Client) SignMessage(address dcrutil.Address, message string) (string, error) {
	return c.SignMessageAsync(address, message).Receive()
}

// FutureVerifyMessageResult is a future promise to deliver the result of a
// VerifyMessageAsync RPC invocation (or an applicable error).
type FutureVerifyMessageResult chan *response

// Receive waits for the response promised by the future and returns whether or
// not the message was successfully verified.
func (r FutureVerifyMessageResult) Receive() (bool, error) {
	res, err := receiveFuture(r)
	if err != nil {
		return false, err
	}

	// Unmarshal result as a boolean.
	var verified bool
	err = json.Unmarshal(res, &verified)
	if err != nil {
		return false, err
	}

	return verified, nil
}

// VerifyMessageAsync returns an instance of a type that can be used to get the
// result of the RPC at some future time by invoking the Receive function on the
// returned instance.
//
// See VerifyMessage for the blocking version and more details.
func (c *Client) VerifyMessageAsync(address dcrutil.Address, signature, message string) FutureVerifyMessageResult {
	addr := address.Address()
	cmd := chainjson.NewVerifyMessageCmd(addr, signature, message)
	return c.sendCmd(cmd)
}

// VerifyMessage verifies a signed message.
//
// NOTE: This function requires to the wallet to be unlocked.  See the
// WalletPassphrase function for more details.
func (c *Client) VerifyMessage(address dcrutil.Address, signature, message string) (bool, error) {
	return c.VerifyMessageAsync(address, signature, message).Receive()
}

// *********************
// Dump/Import Functions
// *********************

// FutureDumpPrivKeyResult is a future promise to deliver the result of a
// DumpPrivKeyAsync RPC invocation (or an applicable error).
type FutureDumpPrivKeyResult chan *response

// Receive waits for the response promised by the future and returns the private
// key corresponding to the passed address encoded in the wallet import format
// (WIF)
func (r FutureDumpPrivKeyResult) Receive(net [2]byte) (*dcrutil.WIF, error) {
	res, err := receiveFuture(r)
	if err != nil {
		return nil, err
	}

	// Unmarshal result as a string.
	var privKeyWIF string
	err = json.Unmarshal(res, &privKeyWIF)
	if err != nil {
		return nil, err
	}

	return dcrutil.DecodeWIF(privKeyWIF, net)
}

// DumpPrivKeyAsync returns an instance of a type that can be used to get the
// result of the RPC at some future time by invoking the Receive function on the
// returned instance.
//
// See DumpPrivKey for the blocking version and more details.
func (c *Client) DumpPrivKeyAsync(address dcrutil.Address) FutureDumpPrivKeyResult {
	addr := address.Address()
	cmd := walletjson.NewDumpPrivKeyCmd(addr)
	return c.sendCmd(cmd)
}

// DumpPrivKey gets the private key corresponding to the passed address encoded
// in the wallet import format (WIF).
//
// NOTE: This function requires to the wallet to be unlocked.  See the
// WalletPassphrase function for more details.
func (c *Client) DumpPrivKey(address dcrutil.Address, net [2]byte) (*dcrutil.WIF, error) {
	return c.DumpPrivKeyAsync(address).Receive(net)
}

// FutureImportPrivKeyResult is a future promise to deliver the result of an
// ImportPrivKeyAsync RPC invocation (or an applicable error).
type FutureImportPrivKeyResult chan *response

// Receive waits for the response promised by the future and returns the result
// of importing the passed private key which must be the wallet import format
// (WIF).
func (r FutureImportPrivKeyResult) Receive() error {
	_, err := receiveFuture(r)
	return err
}

// ImportPrivKeyAsync returns an instance of a type that can be used to get the
// result of the RPC at some future time by invoking the Receive function on the
// returned instance.
//
// See ImportPrivKey for the blocking version and more details.
func (c *Client) ImportPrivKeyAsync(privKeyWIF *dcrutil.WIF) FutureImportPrivKeyResult {
	wif := ""
	if privKeyWIF != nil {
		wif = privKeyWIF.String()
	}

	cmd := walletjson.NewImportPrivKeyCmd(wif, nil, nil, nil)
	return c.sendCmd(cmd)
}

// ImportPrivKey imports the passed private key which must be the wallet import
// format (WIF).
func (c *Client) ImportPrivKey(privKeyWIF *dcrutil.WIF) error {
	return c.ImportPrivKeyAsync(privKeyWIF).Receive()
}

// ImportPrivKeyLabelAsync returns an instance of a type that can be used to get the
// result of the RPC at some future time by invoking the Receive function on the
// returned instance.
//
// See ImportPrivKey for the blocking version and more details.
func (c *Client) ImportPrivKeyLabelAsync(privKeyWIF *dcrutil.WIF, label string) FutureImportPrivKeyResult {
	wif := ""
	if privKeyWIF != nil {
		wif = privKeyWIF.String()
	}

	cmd := walletjson.NewImportPrivKeyCmd(wif, &label, nil, nil)
	return c.sendCmd(cmd)
}

// ImportPrivKeyLabel imports the passed private key which must be the wallet import
// format (WIF). It sets the account label to the one provided.
func (c *Client) ImportPrivKeyLabel(privKeyWIF *dcrutil.WIF, label string) error {
	return c.ImportPrivKeyLabelAsync(privKeyWIF, label).Receive()
}

// ImportPrivKeyRescanAsync returns an instance of a type that can be used to get the
// result of the RPC at some future time by invoking the Receive function on the
// returned instance.
//
// See ImportPrivKey for the blocking version and more details.
func (c *Client) ImportPrivKeyRescanAsync(privKeyWIF *dcrutil.WIF, label string, rescan bool) FutureImportPrivKeyResult {
	wif := ""
	if privKeyWIF != nil {
		wif = privKeyWIF.String()
	}

	cmd := walletjson.NewImportPrivKeyCmd(wif, &label, &rescan, nil)
	return c.sendCmd(cmd)
}

// ImportPrivKeyRescan imports the passed private key which must be the wallet import
// format (WIF). It sets the account label to the one provided. When rescan is true,
// the block history is scanned for transactions addressed to provided privKey.
func (c *Client) ImportPrivKeyRescan(privKeyWIF *dcrutil.WIF, label string, rescan bool) error {
	return c.ImportPrivKeyRescanAsync(privKeyWIF, label, rescan).Receive()
}

// ImportPrivKeyRescanFromAsync returns an instance of a type that can be used to
// get the result of the RPC at some future time by invoking the Receive function o
// n the
// returned instance.
//
// See ImportPrivKey for the blocking version and more details.
func (c *Client) ImportPrivKeyRescanFromAsync(privKeyWIF *dcrutil.WIF, label string, rescan bool, scanFrom int) FutureImportPrivKeyResult {
	wif := ""
	if privKeyWIF != nil {
		wif = privKeyWIF.String()
	}

	cmd := walletjson.NewImportPrivKeyCmd(wif, &label, &rescan, &scanFrom)
	return c.sendCmd(cmd)
}

// ImportPrivKeyRescanFrom imports the passed private key which must be the wallet
// import format (WIF). It sets the account label to the one provided. When rescan
// is true, the block history from block scanFrom is scanned for transactions
// addressed to provided privKey.
func (c *Client) ImportPrivKeyRescanFrom(privKeyWIF *dcrutil.WIF, label string, rescan bool, scanFrom int) error {
	return c.ImportPrivKeyRescanFromAsync(privKeyWIF, label, rescan, scanFrom).Receive()
}

// FutureImportScriptResult is a future promise to deliver the result
// of a ImportScriptAsync RPC invocation (or an applicable error).
type FutureImportScriptResult chan *response

// Receive waits for the response promised by the future.
func (r FutureImportScriptResult) Receive() error {
	_, err := receiveFuture(r)
	return err
}

// ImportScriptAsync returns an instance of a type that can be used to get the
// result of the RPC at some future time by invoking the Receive function on
// the returned instance.
//
// See ImportScript for the blocking version and more details.
func (c *Client) ImportScriptAsync(script []byte) FutureImportScriptResult {
	scriptStr := ""
	if script != nil {
		scriptStr = hex.EncodeToString(script)
	}

	cmd := walletjson.NewImportScriptCmd(scriptStr, nil, nil)
	return c.sendCmd(cmd)
}

// ImportScript attempts to import a byte code script into wallet.
func (c *Client) ImportScript(script []byte) error {
	return c.ImportScriptAsync(script).Receive()
}

// ImportScriptRescanAsync returns an instance of a type that can be used to
// get the result of the RPC at some future time by invoking the Receive
// function on the returned instance.
//
// See ImportScript for the blocking version and more details.
func (c *Client) ImportScriptRescanAsync(script []byte, rescan bool) FutureImportScriptResult {
	scriptStr := ""
	if script != nil {
		scriptStr = hex.EncodeToString(script)
	}

	cmd := walletjson.NewImportScriptCmd(scriptStr, &rescan, nil)
	return c.sendCmd(cmd)
}

// ImportScriptRescan attempts to import a byte code script into wallet. It also
// allows the user to choose whether or not they do a rescan.
func (c *Client) ImportScriptRescan(script []byte, rescan bool) error {
	return c.ImportScriptRescanAsync(script, rescan).Receive()
}

// ImportScriptRescanFromAsync returns an instance of a type that can be used to
// get the result of the RPC at some future time by invoking the Receive
// function on the returned instance.
//
// See ImportScript for the blocking version and more details.
func (c *Client) ImportScriptRescanFromAsync(script []byte, rescan bool, scanFrom int) FutureImportScriptResult {
	scriptStr := ""
	if script != nil {
		scriptStr = hex.EncodeToString(script)
	}

	cmd := walletjson.NewImportScriptCmd(scriptStr, &rescan, &scanFrom)
	return c.sendCmd(cmd)
}

// ImportScriptRescanFrom attempts to import a byte code script into wallet. It
// also allows the user to choose whether or not they do a rescan, and which
// height to rescan from.
func (c *Client) ImportScriptRescanFrom(script []byte, rescan bool, scanFrom int) error {
	return c.ImportScriptRescanFromAsync(script, rescan, scanFrom).Receive()
}

// ***********************
// Miscellaneous Functions
// ***********************

// FutureAccountAddressIndexResult is a future promise to deliver the result of a
// AccountAddressIndexAsync RPC invocation (or an applicable error).
type FutureAccountAddressIndexResult chan *response

// Receive waits for the response promised by the future and returns the info
// provided by the server.
func (r FutureAccountAddressIndexResult) Receive() (int, error) {
	res, err := receiveFuture(r)
	if err != nil {
		return 0, err
	}

	// Unmarshal result as an accountaddressindex result object.
	var index int
	err = json.Unmarshal(res, &index)
	if err != nil {
		return 0, err
	}

	return index, nil
}

// AccountAddressIndexAsync returns an instance of a type that can be used
// to get the result of the RPC at some future time by invoking the Receive
// function on the returned instance.
//
// See AccountAddressIndex for the blocking version and more details.
func (c *Client) AccountAddressIndexAsync(account string, branch uint32) FutureAccountAddressIndexResult {
	cmd := walletjson.NewAccountAddressIndexCmd(account, int(branch))
	return c.sendCmd(cmd)
}

// AccountAddressIndex returns the address index for a given account's branch.
func (c *Client) AccountAddressIndex(account string, branch uint32) (int, error) {
	return c.AccountAddressIndexAsync(account, branch).Receive()
}

// FutureAccountSyncAddressIndexResult is a future promise to deliver the
// result of an AccountSyncAddressIndexAsync RPC invocation (or an
// applicable error).
type FutureAccountSyncAddressIndexResult chan *response

// Receive waits for the response promised by the future and returns the info
// provided by the server.
func (r FutureAccountSyncAddressIndexResult) Receive() error {
	_, err := receiveFuture(r)
	return err
}

// AccountSyncAddressIndexAsync returns an instance of a type that can be used
// to get the result of the RPC at some future time by invoking the Receive
// function on the returned instance.
//
// See AccountSyncAddressIndex for the blocking version and more details.
func (c *Client) AccountSyncAddressIndexAsync(account string, branch uint32, index int) FutureAccountSyncAddressIndexResult {
	cmd := walletjson.NewAccountSyncAddressIndexCmd(account, int(branch), index)
	return c.sendCmd(cmd)
}

// AccountSyncAddressIndex synchronizes an account branch to the passed address
// index.
func (c *Client) AccountSyncAddressIndex(account string, branch uint32, index int) error {
	return c.AccountSyncAddressIndexAsync(account, branch, index).Receive()
}

// FutureRevokeTicketsResult is a future promise to deliver the result of a
// RevokeTicketsAsync RPC invocation (or an applicable error).
type FutureRevokeTicketsResult chan *response

// Receive waits for the response promised by the future and returns the info
// provided by the server.
func (r FutureRevokeTicketsResult) Receive() error {
	_, err := receiveFuture(r)
	return err
}

// RevokeTicketsAsync returns an instance of a type that can be used to get the result
// of the RPC at some future time by invoking the Receive function on the
// returned instance.
//
// See RevokeTickets for the blocking version and more details.
func (c *Client) RevokeTicketsAsync() FutureRevokeTicketsResult {
	cmd := walletjson.NewRevokeTicketsCmd()
	return c.sendCmd(cmd)
}

// RevokeTickets triggers the wallet to issue revocations for any missed tickets that
// have not yet been revoked.
func (c *Client) RevokeTickets() error {
	return c.RevokeTicketsAsync().Receive()
}

// FutureAddTicketResult is a future promise to deliver the result of a
// AddTicketAsync RPC invocation (or an applicable error).
type FutureAddTicketResult chan *response

// Receive waits for the response promised by the future and returns the info
// provided by the server.
func (r FutureAddTicketResult) Receive() error {
	_, err := receiveFuture(r)
	return err
}

// AddTicketAsync returns an instance of a type that can be used to get the result
// of the RPC at some future time by invoking the Receive function on the
// returned instance.
//
// See AddTicket for the blocking version and more details.
func (c *Client) AddTicketAsync(rawHex string) FutureAddTicketResult {
	cmd := walletjson.NewAddTicketCmd(rawHex)
	return c.sendCmd(cmd)
}

// AddTicket manually adds a new ticket to the wallet stake manager. This is used
// to override normal security settings to insert tickets which would not
// otherwise be added to the wallet.
func (c *Client) AddTicket(ticket *dcrutil.Tx) error {
	ticketB, err := ticket.MsgTx().Bytes()
	if err != nil {
		return err
	}

	return c.AddTicketAsync(hex.EncodeToString(ticketB)).Receive()
}

// FutureFundRawTransactionResult is a future promise to deliver the result of a
// FundRawTransactionAsync RPC invocation (or an applicable error).
type FutureFundRawTransactionResult chan *response

// Receive waits for the response promised by the future and returns the unsigned
// transaction with the passed amount and the given address.
func (r FutureFundRawTransactionResult) Receive() (*walletjson.FundRawTransactionResult, error) {
	res, err := receiveFuture(r)
	if err != nil {
		return nil, err
	}

	// Unmarshal result as a string.
	var infoRes walletjson.FundRawTransactionResult
	err = json.Unmarshal(res, &infoRes)
	if err != nil {
		return nil, err
	}

	return &infoRes, nil
}

// FundRawTransactionAsync returns an instance of a type that can be used to get the
// result of the RPC at some future time by invoking the Receive function on the
// returned instance.
//
// See FundRawTransaction for the blocking version and more details.
func (c *Client) FundRawTransactionAsync(rawhex string, fundAccount string, options walletjson.FundRawTransactionOptions) FutureFundRawTransactionResult {
	cmd := walletjson.NewFundRawTransactionCmd(rawhex, fundAccount, &options)
	return c.sendCmd(cmd)
}

// FundRawTransaction Add inputs to a transaction until it has enough
// in value to meet its out value.
func (c *Client) FundRawTransaction(rawhex string, fundAccount string, options walletjson.FundRawTransactionOptions) (*walletjson.FundRawTransactionResult, error) {
	return c.FundRawTransactionAsync(rawhex, fundAccount, options).Receive()
}

// FutureGenerateVoteResult is a future promise to deliver the result of a
// GenerateVoteAsync RPC invocation (or an applicable error).
type FutureGenerateVoteResult chan *response

// Receive waits for the response promised by the future and returns the info
// provided by the server.
func (r FutureGenerateVoteResult) Receive() (*walletjson.GenerateVoteResult, error) {
	res, err := receiveFuture(r)
	if err != nil {
		return nil, err
	}

	// Unmarshal result as a getinfo result object.
	var infoRes walletjson.GenerateVoteResult
	err = json.Unmarshal(res, &infoRes)
	if err != nil {
		return nil, err
	}

	return &infoRes, nil
}

// GenerateVoteAsync returns an instance of a type that can be used to get the result
// of the RPC at some future time by invoking the Receive function on the
// returned instance.
//
// See GenerateVote for the blocking version and more details.
func (c *Client) GenerateVoteAsync(blockHash *chainhash.Hash, height int64, sstxHash *chainhash.Hash, voteBits uint16, voteBitsExt string) FutureGenerateVoteResult {
	cmd := walletjson.NewGenerateVoteCmd(blockHash.String(), height, sstxHash.String(), voteBits, voteBitsExt)
	return c.sendCmd(cmd)
}

// GenerateVote returns hex of an SSGen.
func (c *Client) GenerateVote(blockHash *chainhash.Hash, height int64, sstxHash *chainhash.Hash, voteBits uint16, voteBitsExt string) (*walletjson.GenerateVoteResult, error) {
	return c.GenerateVoteAsync(blockHash, height, sstxHash, voteBits, voteBitsExt).Receive()
}

// NOTE: While getinfo is implemented here (in wallet.go), a dcrd chain server
// will respond to getinfo requests as well, excluding any wallet information.

// FutureGetInfoResult is a future promise to deliver the result of a
// GetInfoAsync RPC invocation (or an applicable error).
type FutureGetInfoResult chan *response

// Receive waits for the response promised by the future and returns the info
// provided by the server.
func (r FutureGetInfoResult) Receive() (*walletjson.InfoWalletResult, error) {
	res, err := receiveFuture(r)
	if err != nil {
		return nil, err
	}

	// Unmarshal result as a getinfo result object.
	var infoRes walletjson.InfoWalletResult
	err = json.Unmarshal(res, &infoRes)
	if err != nil {
		return nil, err
	}

	return &infoRes, nil
}

// GetInfoAsync returns an instance of a type that can be used to get the result
// of the RPC at some future time by invoking the Receive function on the
// returned instance.
//
// See GetInfo for the blocking version and more details.
func (c *Client) GetInfoAsync() FutureGetInfoResult {
	cmd := chainjson.NewGetInfoCmd()
	return c.sendCmd(cmd)
}

// GetInfo returns miscellaneous info regarding the RPC server.  The returned
// info object may be void of wallet information if the remote server does
// not include wallet functionality.
func (c *Client) GetInfo() (*walletjson.InfoWalletResult, error) {
	return c.GetInfoAsync().Receive()
}

// FutureGetStakeInfoResult is a future promise to deliver the result of a
// GetStakeInfoAsync RPC invocation (or an applicable error).
type FutureGetStakeInfoResult chan *response

// Receive waits for the response promised by the future and returns the stake
// info provided by the server.
func (r FutureGetStakeInfoResult) Receive() (*walletjson.GetStakeInfoResult, error) {
	res, err := receiveFuture(r)
	if err != nil {
		return nil, err
	}

	// Unmarshal result as a getstakeinfo result object.
	var infoRes walletjson.GetStakeInfoResult
	err = json.Unmarshal(res, &infoRes)
	if err != nil {
		return nil, err
	}

	return &infoRes, nil
}

// GetStakeInfoAsync returns an instance of a type that can be used to get
// the result of the RPC at some future time by invoking the Receive function
// on the returned instance.
//
// See GetStakeInfo for the blocking version and more details.
func (c *Client) GetStakeInfoAsync() FutureGetStakeInfoResult {
	cmd := walletjson.NewGetStakeInfoCmd()
	return c.sendCmd(cmd)
}

// GetStakeInfo returns stake mining info from a given wallet. This includes
// various statistics on tickets it owns and votes it has produced.
func (c *Client) GetStakeInfo() (*walletjson.GetStakeInfoResult, error) {
	return c.GetStakeInfoAsync().Receive()
}

// FutureGetTicketsResult is a future promise to deliver the result of a
// GetTickets RPC invocation (or an applicable error).
type FutureGetTicketsResult chan *response

// Receive waits for the response promised by the future and returns the info
// provided by the server.
func (r FutureGetTicketsResult) Receive() ([]*chainhash.Hash, error) {
	res, err := receiveFuture(r)
	if err != nil {
		return nil, err
	}

	// Unmarshal result as a gettickets result object.
	var tixRes walletjson.GetTicketsResult
	err = json.Unmarshal(res, &tixRes)
	if err != nil {
		return nil, err
	}

	tickets := make([]*chainhash.Hash, len(tixRes.Hashes))
	for i := range tixRes.Hashes {
		h, err := chainhash.NewHashFromStr(tixRes.Hashes[i])
		if err != nil {
			return nil, err
		}

		tickets[i] = h
	}

	return tickets, nil
}

// GetTicketsAsync returns an instance of a type that can be used to
// get the result of the RPC at some future time by invoking the Receive
// function on the returned instance.
//
// See GetTicketsfor the blocking version and more details.
func (c *Client) GetTicketsAsync(includeImmature bool) FutureGetTicketsResult {
	cmd := walletjson.NewGetTicketsCmd(includeImmature)
	return c.sendCmd(cmd)
}

// GetTickets returns a list of the tickets owned by the wallet, partially
// or in full. The flag includeImmature is used to indicate if non mature
// tickets should also be returned.
func (c *Client) GetTickets(includeImmature bool) ([]*chainhash.Hash, error) {
	return c.GetTicketsAsync(includeImmature).Receive()
}

// FutureListScriptsResult is a future promise to deliver the result of a
// ListScriptsAsync RPC invocation (or an applicable error).
type FutureListScriptsResult chan *response

// Receive waits for the response promised by the future and returns the info
// provided by the server.
func (r FutureListScriptsResult) Receive() ([][]byte, error) {
	res, err := receiveFuture(r)
	if err != nil {
		return nil, err
	}

	// Unmarshal result as a listscripts result object.
	var resScr walletjson.ListScriptsResult
	err = json.Unmarshal(res, &resScr)
	if err != nil {
		return nil, err
	}

	// Convert the redeemscripts into byte slices and
	// store them.
	redeemScripts := make([][]byte, len(resScr.Scripts))
	for i := range resScr.Scripts {
		rs := resScr.Scripts[i].RedeemScript
		rsB, err := hex.DecodeString(rs)
		if err != nil {
			return nil, err
		}
		redeemScripts[i] = rsB
	}

	return redeemScripts, nil
}

// ListScriptsAsync returns an instance of a type that can be used to
// get the result of the RPC at some future time by invoking the Receive
// function on the returned instance.
//
// See ListScripts for the blocking version and more details.
func (c *Client) ListScriptsAsync() FutureListScriptsResult {
	cmd := walletjson.NewListScriptsCmd()
	return c.sendCmd(cmd)
}

// ListScripts returns a list of the currently known redeemscripts from the
// wallet as a slice of byte slices.
func (c *Client) ListScripts() ([][]byte, error) {
	return c.ListScriptsAsync().Receive()
}

// FutureSetTicketFeeResult is a future promise to deliver the result of a
// SetTicketFeeAsync RPC invocation (or an applicable error).
type FutureSetTicketFeeResult chan *response

// Receive waits for the response promised by the future and returns the info
// provided by the server.
func (r FutureSetTicketFeeResult) Receive() error {
	_, err := receiveFuture(r)
	return err
}

// SetTicketFeeAsync returns an instance of a type that can be used to
// get the result of the RPC at some future time by invoking the Receive
// function on the returned instance.
//
// See SetTicketFee for the blocking version and more details.
func (c *Client) SetTicketFeeAsync(fee dcrutil.Amount) FutureSetTicketFeeResult {
	cmd := walletjson.NewSetTicketFeeCmd(fee.ToCoin())
	return c.sendCmd(cmd)
}

// SetTicketFee sets the ticket fee per KB amount.
func (c *Client) SetTicketFee(fee dcrutil.Amount) error {
	return c.SetTicketFeeAsync(fee).Receive()
}

// FutureSetTxFeeResult is a future promise to deliver the result of a
// SetTxFeeAsync RPC invocation (or an applicable error).
type FutureSetTxFeeResult chan *response

// Receive waits for the response promised by the future.
func (r FutureSetTxFeeResult) Receive() error {
	_, err := receiveFuture(r)
	return err
}

// SetTxFeeAsync returns an instance of a type that can be used to
// get the result of the RPC at some future time by invoking the Receive
// function on the returned instance.
//
// See SetTxFee for the blocking version and more details.
func (c *Client) SetTxFeeAsync(fee dcrutil.Amount) FutureSetTxFeeResult {
	cmd := walletjson.NewSetTxFeeCmd(fee.ToCoin())
	return c.sendCmd(cmd)
}

// SetTxFee sets the transaction fee per KB amount.
func (c *Client) SetTxFee(fee dcrutil.Amount) error {
	return c.SetTxFeeAsync(fee).Receive()
}

// FutureGetVoteChoicesResult is a future promise to deliver the result of a
// GetVoteChoicesAsync RPC invocation (or an applicable error).
type FutureGetVoteChoicesResult chan *response

// Receive waits for the response promised by the future.
func (r FutureGetVoteChoicesResult) Receive() (*walletjson.GetVoteChoicesResult, error) {
	res, err := receiveFuture(r)
	if err != nil {
		return nil, err
	}

	// Unmarshal result as a getvotechoices result object.
	var choicesRes walletjson.GetVoteChoicesResult
	err = json.Unmarshal(res, &choicesRes)
	if err != nil {
		return nil, err
	}

	return &choicesRes, nil
}

// GetVoteChoicesAsync returns an instance of a type that can be used to get the
// result of the RPC at some future time by invoking the Receive function on the
// returned instance.
//
// See GetVoteChoices for the blocking version and more details.
func (c *Client) GetVoteChoicesAsync() FutureGetVoteChoicesResult {
	cmd := walletjson.NewGetVoteChoicesCmd()
	return c.sendCmd(cmd)
}

// GetVoteChoices returns the currently-set vote choices for each agenda in the
// latest supported stake version.
func (c *Client) GetVoteChoices() (*walletjson.GetVoteChoicesResult, error) {
	return c.GetVoteChoicesAsync().Receive()
}

// FutureSetVoteChoiceResult is a future promise to deliver the result of a
// SetVoteChoiceAsync RPC invocation (or an applicable error).
type FutureSetVoteChoiceResult chan *response

// Receive waits for the response promised by the future.
func (r FutureSetVoteChoiceResult) Receive() error {
	_, err := receiveFuture(r)
	return err
}

// SetVoteChoiceAsync returns an instance of a type that can be used to get the
// result of the RPC at some future time by invoking the Receive function on the
// returned instance.
//
// See SetVoteChoice for the blocking version and more details.
func (c *Client) SetVoteChoiceAsync(agendaID, choiceID string) FutureSetVoteChoiceResult {
	cmd := walletjson.NewSetVoteChoiceCmd(agendaID, choiceID)
	return c.sendCmd(cmd)
}

// SetVoteChoice sets a voting choice preference for an agenda.
func (c *Client) SetVoteChoice(agendaID, choiceID string) error {
	return c.SetVoteChoiceAsync(agendaID, choiceID).Receive()
}

// FutureStakePoolUserInfoResult is a future promise to deliver the result of a
// GetInfoAsync RPC invocation (or an applicable error).
type FutureStakePoolUserInfoResult chan *response

// Receive waits for the response promised by the future and returns the info
// provided by the server.
func (r FutureStakePoolUserInfoResult) Receive() (*walletjson.StakePoolUserInfoResult, error) {
	res, err := receiveFuture(r)
	if err != nil {
		return nil, err
	}

	// Unmarshal result as a stakepooluserinfo result object.
	var infoRes walletjson.StakePoolUserInfoResult
	err = json.Unmarshal(res, &infoRes)
	if err != nil {
		return nil, err
	}

	return &infoRes, nil
}

// StakePoolUserInfoAsync returns an instance of a type that can be used to
// get the result of the RPC at some future time by invoking the Receive
// function on the returned instance.
//
// See GetInfo for the blocking version and more details.
func (c *Client) StakePoolUserInfoAsync(addr dcrutil.Address) FutureStakePoolUserInfoResult {
	cmd := walletjson.NewStakePoolUserInfoCmd(addr.Address())
	return c.sendCmd(cmd)
}

// StakePoolUserInfo returns a list of tickets and information about them
// that are paying to the passed address.
func (c *Client) StakePoolUserInfo(addr dcrutil.Address) (*walletjson.StakePoolUserInfoResult, error) {
	return c.StakePoolUserInfoAsync(addr).Receive()
}

// FutureTicketsForAddressResult is a future promise to deliver the result of a
// GetInfoAsync RPC invocation (or an applicable error).
type FutureTicketsForAddressResult chan *response

// Receive waits for the response promised by the future and returns the info
// provided by the server.
func (r FutureTicketsForAddressResult) Receive() (*chainjson.TicketsForAddressResult, error) {
	res, err := receiveFuture(r)
	if err != nil {
		return nil, err
	}

	// Unmarshal result as a ticketsforaddress result object.
	var infoRes chainjson.TicketsForAddressResult
	err = json.Unmarshal(res, &infoRes)
	if err != nil {
		return nil, err
	}

	return &infoRes, nil
}

// TicketsForAddressAsync returns an instance of a type that can be used to
// get the result of the RPC at some future time by invoking the Receive
// function on the returned instance.
//
// See GetInfo for the blocking version and more details.
func (c *Client) TicketsForAddressAsync(addr dcrutil.Address) FutureTicketsForAddressResult {
	cmd := chainjson.NewTicketsForAddressCmd(addr.Address())
	return c.sendCmd(cmd)
}

// TicketsForAddress returns a list of tickets paying to the passed address.
// If the daemon server is queried, it returns a search of tickets in the
// live ticket pool. If the wallet server is queried, it searches all tickets
// owned by the wallet.
func (c *Client) TicketsForAddress(addr dcrutil.Address) (*chainjson.TicketsForAddressResult, error) {
	return c.TicketsForAddressAsync(addr).Receive()
}

// FutureWalletInfoResult is a future promise to deliver the result of a
// WalletInfoAsync RPC invocation (or an applicable error).
type FutureWalletInfoResult chan *response

// Receive waits for the response promised by the future and returns the stake
// info provided by the server.
func (r FutureWalletInfoResult) Receive() (*walletjson.WalletInfoResult, error) {
	res, err := receiveFuture(r)
	if err != nil {
		return nil, err
	}

	// Unmarshal result as a walletinfo result object.
	var infoRes walletjson.WalletInfoResult
	err = json.Unmarshal(res, &infoRes)
	if err != nil {
		return nil, err
	}

	return &infoRes, nil
}

// WalletInfoAsync returns an instance of a type that can be used to get
// the result of the RPC at some future time by invoking the Receive function
// on the returned instance.
//
// See WalletInfo for the blocking version and more details.
func (c *Client) WalletInfoAsync() FutureWalletInfoResult {
	cmd := walletjson.NewWalletInfoCmd()
	return c.sendCmd(cmd)
}

// WalletInfo returns wallet global state info for a given wallet.
func (c *Client) WalletInfo() (*walletjson.WalletInfoResult, error) {
	return c.WalletInfoAsync().Receive()
}

// TODO(davec): Implement
// backupwallet (NYI in dcrwallet)
// encryptwallet (Won't be supported by dcrwallet since it's always encrypted)
// getwalletinfo (NYI in dcrwallet or dcrjson)
// listaddressgroupings (NYI in dcrwallet)
// listreceivedbyaccount (NYI in dcrwallet)

// DUMP
// importwallet (NYI in dcrwallet)
// dumpwallet (NYI in dcrwallet)
