package transformer

import (
	"encoding/json"

	"github.com/dcb9/janus/pkg/eth"
	"github.com/dcb9/janus/pkg/qtum"
	"github.com/dcb9/janus/pkg/rpc"
	"github.com/pkg/errors"
)

func (m *Manager) sendtocontract(req *rpc.JSONRPCRequest, tx *eth.TransactionReq) error {
	gasLimit, gasPrice, err := EthGasToQtum(tx)
	if err != nil {
		return err
	}

	amount := 0.0
	if tx.Value != "" {
		var err error
		amount, err = EthValueToQtumAmount(tx.Value)
		if err != nil {
			return errors.Wrap(err, "EthValueToQtumAmount:")
		}
	}

	params := []interface{}{
		RemoveHexPrefix(tx.To),
		RemoveHexPrefix(tx.Data),
		amount,
		gasLimit,
		gasPrice,
	}

	if from := tx.From; from != "" {
		if IsEthHexAddress(from) {
			from, err = m.qtumClient.FromHexAddress(RemoveHexPrefix(from))
			if err != nil {
				return err
			}
		}
		params = append(params, from)
	}

	newParams, err := json.Marshal(params)
	if err != nil {
		return err
	}

	req.Params = newParams
	req.Method = qtum.MethodSendtocontract

	return nil
}

//  Eth RPC
//  params: [{
//    "from": "0xb60e8dd61c5d32be8058bb8eb970870f07233155",
//    "to": "0xd46e8dd67c5d32be8058bb8eb970870f07244567",
//    "gas": "0x76c0", // 30400
//    "gasPrice": "0x9184e72a000", // 10000000000000
//    "value": "0x9184e72a", // 2441406250
//    "data": "0xd46e8dd67c5d32be8d46e8dd67c5d32be8058bb8eb970870f072445675058bb8eb970870f072445675"
//  }]

//Qtum RPC
//  sendtocontract "contractaddress" "data" (amount gaslimit gasprice senderaddress broadcast)
//  Send funds and data to a contract.
//
//Arguments:
//  1. "contractaddress" (string, required) The contract address that will receive the funds and data.
//  2. "datahex"  (string, required) data to send.
//  3. "amount"      (numeric or string, optional) The amount in QTUM to send. eg 0.1, default: 0
//  4. gasLimit  (numeric or string, optional) gasLimit, default: 250000, max: 40000000
//  5. gasPrice  (numeric or string, optional) gasPrice Qtum price per gas unit, default: 0.0000004, min:0.0000004
//  6. "senderaddress" (string, optional) The quantum address that will be used as sender.
//  7. "broadcast" (bool, optional, default=true) Whether to broadcast the transaction or not.
//  8. "changeToSender" (bool, optional, default=true) Return the change to the sender.
//
//Result:
//  [
//  {
//  "txid" : (string) The transaction id.
//  "sender" : (string) QTUM address of the sender.
//  "hash160" : (string) ripemd-160 hash of the sender.
//  }
//  ]
//
//Examples:
//  > qtum-cli sendtocontract "c6ca2697719d00446d4ea51f6fac8fd1e9310214" "54f6127f"
//  > qtum-cli sendtocontract "c6ca2697719d00446d4ea51f6fac8fd1e9310214" "54f6127f" 12.0015 6000000 0.0000004 "QM72Sfpbz1BPpXFHz9m3CdqATR44Jvaydd"