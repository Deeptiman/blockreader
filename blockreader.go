package main

import (
	"fmt"	
	"encoding/json"
	"github.com/hyperledger/fabric-sdk-go/pkg/core/config"
	"github.com/hyperledger/fabric-sdk-go/pkg/fabsdk"
	contextAPI "github.com/hyperledger/fabric-sdk-go/pkg/common/providers/context"
	"github.com/hyperledger/fabric-protos-go/peer"
	"github.com/hyperledger/fabric-sdk-go/pkg/client/ledger"
	"github.com/hyperledger/fabric-protos-go/common"
	"github.com/golang/protobuf/proto"
	"github.com/pkg/errors"
)

var channelCtx	contextAPI.ChannelProvider

const (

	ChannelID = "employeeledger"
	OrgAdmin = "Admin"
	OrgName = "org1"
	ConfigFile = "config.yaml"
	TxnID = "cd2b072c880cdefbea66c5f9d73a5a5eb3c3977e77772fba42cec59204ca2980"
)

func Initialize() error {

	sdk, err := fabsdk.New(config.FromFile(ConfigFile))
	if err != nil {
		return errors.WithMessage(err, "failed to create SDK")
	}

	fmt.Println("SDK Initialized Successfully")

	channelCtx = sdk.ChannelContext(ChannelID, 
		fabsdk.WithUser(OrgAdmin), 
		fabsdk.WithOrg(OrgName))

	fmt.Println("Channel Context Initialized Successfully")
	return nil
}

func CreateLedgerClient()(*ledger.Client, error){

	lc , err := ledger.New(channelCtx)

	if err != nil {
		return nil, errors.WithMessage(err, "failed to create ledger client")
	}
	fmt.Println("Ledger Client Created Successfully\n")
	return lc, nil
}

func BlockReader() error {

	lc , err := CreateLedgerClient()
	if err != nil {
		return errors.WithMessage(err, "failed to create ledger client")
	}

	block, err := QueryBlock(lc)
	if err != nil {
		return errors.WithMessage(err, "failed to query block")
	}	

	/********************** READ THE BLOCK **************************/	

	//BlockData - 
	blockData := block.Data.Data

	//First Get the Envelope from the BlockData

	/*
		type Envelope struct {
			Payload   []byte 
			Signature []byte
		}
	*/
	envelope, err := GetEnvelopeFromBlock(blockData[0])
	if err != nil {
		return errors.WithMessage(err,"unmarshaling Envelope error: ")
	}	

	//Retrieve the Payload from the Envelope
	/*
		type Payload struct {
			Header *Header 
			Data   []byte 
		}
	*/
	payload := &common.Payload{}
	err = proto.Unmarshal(envelope.Payload, payload)
	if err != nil {
		return errors.WithMessage(err,"unmarshaling Payload error: ")
	}
	
	payloadJson, err := GetPayloadJson(envelope)
	if err != nil {
		return errors.WithMessage(err,"unmarshaling Payload error: ")	
	}

	//Read the Transaction from the Payload Data

	/*
		type TransactionAction struct {
			Header  []byte 
			Payload []byte 
		}
	*/
	transaction := &peer.Transaction{}
	err = proto.Unmarshal(payload.Data, transaction)
	if err != nil {
		return errors.WithMessage(err,"unmarshaling Payload Transaction error: ")
	}

	// The Header is simialr to the Payload Header retrieving the Creator Identity and Nonce //

	// Payload field is marshalled object of ChaincodeActionPayload

	/*
		type ChaincodeActionPayload struct {
			ChaincodeProposalPayload []byte 
			Action *ChaincodeEndorsedAction 
		}
	*/

	chaincodeActionPayload := &peer.ChaincodeActionPayload{}
	err = proto.Unmarshal(transaction.Actions[0].Payload, chaincodeActionPayload)
	if err != nil {
		return errors.WithMessage(err,"unmarshaling Chaincode Action Payload error: ")
	}

	transactionJson, err := GetTransactionJson(chaincodeActionPayload)
	if err != nil {
		return errors.WithMessage(err,"failed to get Transaction Json error: ")
	}	 
		

		headerJson := Header{
			Payload: payloadJson,
		}

		dataJson := Data {
			Transaction: transactionJson,
		}

		envelopeJson := Envelope {
			Header: headerJson,
			Data: 	dataJson,
		}

		blockReader := Block{
			Envelope: envelopeJson,
		}

		fmt.Println("************* JSON FORMAT ************* ")
		var jsonData []byte
		jsonData, err = json.MarshalIndent(blockReader, "","    ")
		if err != nil {
			errors.WithMessage(err, "failed to marshal Json")
		}
		fmt.Println(string(jsonData))

	return nil	
}

func QueryBlock(lc *ledger.Client) (*common.Block, error){

	block, err := lc.QueryBlockByTxID(TxnID)
	if err != nil {
		return nil, errors.WithMessage(err, "failed to query block by transaction ID")
	}

	if block == nil {
		return nil, errors.New("No Block exists for this TxnID - "+TxnID)
	}

	return block, nil
}

func GetEnvelopeFromBlock(data []byte) (*common.Envelope, error){

	var err error
	env := &common.Envelope{}
	if err = proto.Unmarshal(data, env); err != nil {
		return nil, errors.Wrap(err, "error unmarshaling Envelope")
	}

	return env, nil
}

func CToGoString(c []byte) string {
    n := -1
    for i, b := range c {
        if b == 0 {
            break
        }
        n = i
    }
    return string(c[:n+1])
}

func main() {

	err := Initialize()
	if err != nil {
		fmt.Println("failed to initialize")
	}

	err = BlockReader()
	if err != nil {
		fmt.Println(" failed to read the Block - ", err)
	}

}




