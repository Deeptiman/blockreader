package main

import (
	"fmt"	
	"flag"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"github.com/hyperledger/fabric-sdk-go/pkg/common/providers/fab"
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

var ChannelID string
var TxnID string

const (
	OrgAdmin = "Admin"
	OrgName = "org1"
	ConfigFile = "config.yaml"
)

type BlockMetadataIndex int32
const (
    BlockMetadataIndex_SIGNATURES          BlockMetadataIndex = 0
    BlockMetadataIndex_LAST_CONFIG         BlockMetadataIndex = 1 // Deprecated: Do not use.
    BlockMetadataIndex_TRANSACTIONS_FILTER BlockMetadataIndex = 2
    BlockMetadataIndex_ORDERER             BlockMetadataIndex = 3 // Deprecated: Do not use.
    BlockMetadataIndex_COMMIT_HASH         BlockMetadataIndex = 4
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

	/*		
		type Block struct {
			Header   *BlockHeader   
			Data     *BlockData     
			Metadata *BlockMetadata 
		}
	*/



	/////////////////////////// -> BlockHeader <- //////////////////////////////////

	blockHeader := block.Header

	/*
		type BlockHeader struct {
			Number   		uint64
			PreviousHash    []byte
			DataHash        []byte
		}
	*/

	previousHash := sha256.Sum256(blockHeader.PreviousHash)
	dataHash := sha256.Sum256(blockHeader.DataHash)

	blockHeaderJson := BlockHeader {
		Number: 		blockHeader.Number,
		PreviousHash:	hex.EncodeToString(previousHash[:]),
		DataHash:		hex.EncodeToString(dataHash[:]),
	}
	//////////////////////////////////////////////////////////////////////////////////


	/////////////////////////// -> BlockData <- ////////////////////////////////// 

	/*
		type BlockData struct {
			Data   [][]byte
		}
	*/
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
	
	payloadJson, err := GetPayloadJson(payload)
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

	blockDataJson := BlockData {
		Envelope: envelopeJson,
	}
	//////////////////////////////////////////////////////////////////////////////////

	
	/////////////////////////// -> BlockMetaData <- //////////////////////////////////
	
	blockMetaData := block.Metadata
	metadata := &common.Metadata{}
	err = proto.Unmarshal(blockMetaData.Metadata[BlockMetadataIndex_SIGNATURES], metadata)
	if err != nil {
		return errors.Wrapf(err, "error unmarshaling metadata")
	}

	/*
		type Metadata struct {
			Value   	[]byte
			Signatures  []*MetadataSignature
		}

		type MetadataSignature struct {
			SignatureHeader  []byte
			Signature        []byte
		}
	*/
		
		signatureHeader := &common.SignatureHeader{}
		err = proto.Unmarshal(metadata.Signatures[0].SignatureHeader, signatureHeader)
		if err != nil {
			return errors.WithMessage(err,"unmarshaling Signature Header error: ")
		}

		signatureHeaderJson, err := GetSignatureHeaderJson(signatureHeader)
		if err != nil {
			return errors.WithMessage(err, "failed get Signature Header")
		}

		blockMetaDataJson := BlockMetaData{
			Value: metadata.Value,
			Signature: metadata.Signatures[0].Signature,
			SignatureHeader: signatureHeaderJson,
		}
		//////////////////////////////////////////////////////////////////////////////////



		blockReader := Block{
			BlockHeader: 	blockHeaderJson,
			BlockData: 		blockDataJson,
			BlockMetaData: 	blockMetaDataJson,
		}
		
		fmt.Println("************* BLOCK READER JSON ************* ")
		var jsonData []byte
		jsonData, err = json.MarshalIndent(blockReader, "","    ")
		if err != nil {
			errors.WithMessage(err, "failed to marshal Json")
		}
		fmt.Println(string(jsonData))

	return nil	
}

func QueryBlock(lc *ledger.Client) (*common.Block, error){

	block, err := lc.QueryBlockByTxID(fab.TransactionID(TxnID))
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

	flag.StringVar(&ChannelID, "channelId", "", "add channel name")
	flag.StringVar(&TxnID, "txnId",  "", "add txnId")
	flag.Parse()	

	if len(TxnID)==0 || len(ChannelID) == 0 {
		fmt.Println("Please add the 'txnId' and 'channelId' to continue...")
	} else {

		err := Initialize()
		if err != nil {
			fmt.Println("failed to initialize")
		}

		err = BlockReader()
		if err != nil {
			fmt.Println(" failed to read the Block - ", err)
		}
	}

}
