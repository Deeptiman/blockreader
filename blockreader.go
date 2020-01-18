package main

import (
	"fmt"
	"encoding/base64"
	"crypto/sha256"
	"crypto/x509"
	"encoding/pem"
	"encoding/json"
	"encoding/hex"
	"github.com/hyperledger/fabric-protos-go/ledger/rwset/kvrwset"
	"github.com/hyperledger/fabric-protos-go/ledger/rwset"
	"github.com/hyperledger/fabric-protos-go/msp"
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

	fmt.Println("*********READ BLOCK HEADER FROM THE PAYLOAD*******")

	/*
		type Header struct {
			ChannelHeader   []byte 
			SignatureHeader []byte 
		}
	*/

	// 1. ChannelHeader
		/*
			type ChannelHeader struct {
				Type      int32 
				Version   int32 
				Timestamp *google_protobuf.Timestamp 
				ChannelId string 
				TxId      string 
				Epoch     uint64 
				Extension []byte 
			}

		*/
	
		channelHeader := &common.ChannelHeader{}
		err = proto.Unmarshal(payload.Header.ChannelHeader, channelHeader)
		if err != nil {
			return errors.WithMessage(err,"unmarshaling Channel Header error: ")
		}

		fmt.Println(" ******* Channel Headers ********")
		fmt.Println(" ChannelId = ", channelHeader.ChannelId)
		fmt.Println(" Type = ", channelHeader.Type)
		fmt.Println(" Version = ", channelHeader.Version)
		fmt.Println(" Timestamp = ", channelHeader.Timestamp)
		fmt.Println(" TxId = ", channelHeader.TxId)

		// The Exension field marshalled object from ChaincodeHeaderExtension

		/*
			type ChaincodeHeaderExtension struct {
				PayloadVisibility []byte 
				ChaincodeId  *ChaincodeID
			}

			type ChaincodeID struct {
				Path    string 
				Name    string 
				Version string 
			}
		*/
		fmt.Println(" --> Extension")

		extension := &peer.ChaincodeHeaderExtension{}
		err = proto.Unmarshal(channelHeader.Extension, extension)
		if err != nil {
			return errors.WithMessage(err,"unmarshaling Channel Header error: ")
		}
		fmt.Println(" 1. Chaincode Path = "+extension.ChaincodeId.Path)
		fmt.Println(" 2. Chaincode Name = "+extension.ChaincodeId.Name)
		fmt.Println(" 3. Chaincode Version = "+extension.ChaincodeId.Version)
		fmt.Println(" *******************************")


		chaincodeIdJson := ChaincodeID {
			Path:		extension.ChaincodeId.Path,
			Name:		extension.ChaincodeId.Name,
			Version:	extension.ChaincodeId.Version,
		}

		chaincodeHeaderExtensionJson := ChaincodeHeaderExtension{
			ChaincodeId: chaincodeIdJson,	
		}

		channelHeaderJson := ChannelHeader{
			Type: 		channelHeader.Type,
			Version:	channelHeader.Version,
			ChannelId:	channelHeader.ChannelId,
			TxId:		channelHeader.TxId,
			Epoch:		channelHeader.Epoch,
			Extension:	chaincodeHeaderExtensionJson,
		}		

	// 2. SignatureHeader

		/*
			type SignatureHeader struct {
				Creator []byte
				Nonce   []byte
			}
		*/

		// Creator is the marshalled object of msp.SerializedIdentity
		signatureHeader := &common.SignatureHeader{}
		err = proto.Unmarshal(payload.Header.SignatureHeader, signatureHeader)
		if err != nil {
			return errors.WithMessage(err,"unmarshaling Signature Header error: ")
		}

		creator := &msp.SerializedIdentity{}
		err = proto.Unmarshal(signatureHeader.Creator, creator)
		if err != nil {
			return errors.WithMessage(err,"unmarshaling Creator error: ")
		}

		fmt.Println(" \n******* Signature Headers ********")
		fmt.Println(" Creator = "+creator.Mspid)
		uEnc := base64.URLEncoding.EncodeToString([]byte(creator.IdBytes))

		// Base64 Url Decoding
		certText, err := base64.URLEncoding.DecodeString(uEnc)
		if err != nil {
			return errors.WithMessage(err,"Error decoding string: ")
		}
		
		fmt.Println("IdBytes = "+string(certText))

		end, _ := pem.Decode([]byte(string(certText)))
		if end == nil {
			panic("failed to parse certificate PEM")
		}
		cert, err := x509.ParseCertificate(end.Bytes)
		if err != nil {
			panic("failed to parse certificate: " + err.Error())
		}
		fmt.Println("Parse Certificate", cert.Issuer, cert.Subject, cert.SerialNumber, cert.NotBefore, cert.NotAfter, cert.PermittedEmailAddresses)
		

		certificateJson :=	Certificate{
			Country:			cert.Issuer.Country,
			Organization:		cert.Issuer.Organization,
			OrganizationalUnit:	cert.Issuer.OrganizationalUnit,
			Locality:			cert.Issuer.Locality,
			Province:			cert.Issuer.Province,
			SerialNumber:		cert.Issuer.SerialNumber,
			NotBefore:			cert.NotBefore,
			NotAfter:			cert.NotAfter,								
		}

		creatorJson := Creator {
			Mspid:  		creator.Mspid,
			CertText:		string(certText),
			Certificate:	certificateJson,
		}

		signatureHeaderJson := SignatureHeader {				
			Creator: creatorJson,
		}

		payloadJson := Payload {
			ChannelHeader: 		channelHeaderJson,
			SignatureHeader:	signatureHeaderJson,
		}

		// IdBytes SigningIdentityInfo
		fmt.Println(" *******************************\n")


	fmt.Println("--------------------------------------------")
	
	
	fmt.Println("\n*********READ BLOCK DATA FROM THE PAYLOAD*******")

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

	// 1. ChaincodeProposalPayload

		/*
			type ChaincodeProposalPayload struct {
				Input        []byte 
				TransientMap map[string][]byte 
			}
		*/

		chaincodeProposalPayload := &peer.ChaincodeProposalPayload{}
		err = proto.Unmarshal(chaincodeActionPayload.ChaincodeProposalPayload, chaincodeProposalPayload)
		if err != nil {
			return errors.WithMessage(err,"unmarshaling Chaincode Proposal Payload error: ")
		}

		// The Input field is marshalled object of ChaincodeInvocationSpec
		input := &peer.ChaincodeInvocationSpec{}
		err = proto.Unmarshal(chaincodeProposalPayload.Input, input)
		if err != nil {
			return errors.WithMessage(err,"unmarshaling Chaincode Proposal Payload Input error: ")
		}	

		/*
			type ChaincodeInvocationSpec struct {
				ChaincodeSpec  *ChaincodeSpec
			}
			type ChaincodeSpec struct {
				Type 		ChaincodeSpec_Type 
				ChaincodeId *ChaincodeID
				Input       *ChaincodeInput
				Timeout     int32
			}
			type ChaincodeInput struct {
				Args        [][]byte
				Decorations map[string][]byte
			}
		*/

		fmt.Println(" ******** Chanicode Input Spect ****** ")
		fmt.Println("Type = ", input.ChaincodeSpec.Type)
		fmt.Println("ChaincodeId = "+input.ChaincodeSpec.ChaincodeId.Name)
		fmt.Println(" --> Read the Chaincode Input Args \n")
		chaincodeArgs := make([]string, len(input.ChaincodeSpec.Input.Args))

		for i, c := range input.ChaincodeSpec.Input.Args {
			args := CToGoString(c[:])
			chaincodeArgs[i] = args
			fmt.Println(" --[CC Args == ",i," --> "+args, len(c), "]")
		}
		fmt.Println(" *********************************** \n\n")

		chaincodeSpecJson := ChaincodeSpec{
			ChaincodeId: 	input.ChaincodeSpec.ChaincodeId.Name,
			ChaincodeType:  string(input.ChaincodeSpec.Type),
			ChaincodeArgs:	chaincodeArgs,
		}

		chaincodeInvocationSpecJson := ChaincodeInvocationSpec{
			ChaincodeSpec:  chaincodeSpecJson,
		}

		chaincodeProposalPayloadJson := ChaincodeProposalPayload {
			ChaincodeInvocationSpec: chaincodeInvocationSpecJson,
		}

	//2. Action

		/*
			type ChaincodeEndorsedAction struct {

				ProposalResponsePayload []byte
				Endorsements  []*Endorsement
			}
		*/

		// ProposalResponsePayload field is the marshalled object of ProposalResponsePayload

		/*
			type ProposalResponsePayload struct {

				ProposalHash []byte
				Extension    []byte   
			}
		*/
		proposalResponsePayload	:= &peer.ProposalResponsePayload{}
		err = proto.Unmarshal(chaincodeActionPayload.Action.ProposalResponsePayload, proposalResponsePayload)
		if err != nil {
			return errors.WithMessage(err,"unmarshaling Proposal Response Payload error: ")
		}

		fmt.Println(" ******** Proposal Response Payload ******* ")		
		proposalHash := sha256.Sum256(proposalResponsePayload.ProposalHash)
		fmt.Printf(" ProposalHash = %x", proposalHash)
		fmt.Println(" \n****************************************** ")

		// Extension is the marshalled object of ChaincodeAction
		/*
			type ChaincodeAction struct {
				Results  []byte 
				Events   []byte 
				Response *Response 
			}
		*/
		chaincodeAction := &peer.ChaincodeAction{}
		err = proto.Unmarshal(proposalResponsePayload.Extension, chaincodeAction)
		if err != nil {
			return errors.WithMessage(err,"unmarshaling Extension error: ")
		}	

		// Results is the marshalled object of TxReadWriteSet
		/*
			type TxReadWriteSet struct {
				DataModel TxReadWriteSet_DataModel 
				NsRwset   []*NsReadWriteSet        
			}

			type NsReadWriteSet struct {
				Namespace string 
				Rwset     []byte 
			}
		*/
		txReadWriteSet := &rwset.TxReadWriteSet{}
		err = proto.Unmarshal(chaincodeAction.Results, txReadWriteSet)
		if err != nil {
			return errors.WithMessage(err,"unmarshaling txReadWriteSet error: ")
		}
		
		nameSpace := txReadWriteSet.NsRwset[0].Namespace
		RwSet := txReadWriteSet.NsRwset[1].Rwset

		//RwSet is the marshalled object of KVRWSet
		/*
			type KVRWSet struct {
				Reads            []*KVRead         
				RangeQueriesInfo []*RangeQueryInfo 
				Writes           []*KVWrite  
				MetadataWrites   []*KVMetadataWrite      
			}

			type KVRead struct {
				Key     string   
				Version *Version 
			}

			type Version struct {
				BlockNum uint64 
				TxNum    uint64 
			}

			type KVWrite struct {
				Key      string 
				IsDelete bool   
				Value    []byte 
			}

			type RangeQueryInfo struct {
				StartKey     string 
				EndKey       string 
				ItrExhausted bool   
				ReadsInfo isRangeQueryInfo_ReadsInfo 
			}

			type KVMetadataWrite struct {
				Key      string
				Entries  []*KVMetadataEntry
			}

			type KVMetadataEntry struct {
				Name     string
				Value    []byte
			}
		*/
		kvrwset := &kvrwset.KVRWSet{}
		err = proto.Unmarshal(RwSet, kvrwset)
		if err != nil {
			return errors.WithMessage(err,"unmarshaling kvrwset error: ")
		}
		fmt.Println(" Namespace = "+nameSpace)
		fmt.Println(" ****** KV Read Write Set ****** ")

		var versionJson			Version
		var kvReadJson 			KVRead
		var kvWriteJson			KVWrite
		var rangeQueryInfoJson	RangeQueryInfo
		var kvMetadataWriteJson KVMetadataWrite

		if len(kvrwset.Reads) != 0 {
			fmt.Println(" ## Read")
			fmt.Println(" Key = "+kvrwset.Reads[0].Key)
			fmt.Println(" BlockNum = ", kvrwset.Reads[0].Version.BlockNum)
			fmt.Println(" TxNum = ", kvrwset.Reads[0].Version.TxNum)

			versionJson = Version{
				BlockNum: kvrwset.Reads[0].Version.BlockNum,
				TxNum:  kvrwset.Reads[0].Version.TxNum,
			}	

			kvReadJson = KVRead{
				Key: 	 kvrwset.Reads[0].Key,
				Version: versionJson,
			}
		}
		
		if len(kvrwset.Writes) != 0 {
			fmt.Println(" \n## Write")
			fmt.Println(" Key = "+kvrwset.Writes[0].Key)
			fmt.Println(" IsDelete = ",kvrwset.Writes[0].IsDelete)		
			value := CToGoString(kvrwset.Writes[0].Value[:])
			fmt.Println(" Value = "+value)

			kvWriteJson = KVWrite {
				Key: 		kvrwset.Writes[0].Key,
				IsDelete: 	kvrwset.Writes[0].IsDelete,
			}
		}

		if len(kvrwset.RangeQueriesInfo) != 0 {
			fmt.Println(" \n## RangeQueriesInfo")
			fmt.Println(" Start Key = "+kvrwset.RangeQueriesInfo[0].StartKey)
			fmt.Println(" End Key = "+kvrwset.RangeQueriesInfo[0].EndKey)
			fmt.Println(" ItrExhausted = ",kvrwset.RangeQueriesInfo[0].ItrExhausted)
			fmt.Println(" ReadsInfo = ",kvrwset.RangeQueriesInfo[0].ReadsInfo)

			rangeQueryInfoJson = RangeQueryInfo{
				StartKey:  kvrwset.RangeQueriesInfo[0].StartKey,
				EndKey:	   kvrwset.RangeQueriesInfo[0].EndKey,
				ItrExhausted: kvrwset.RangeQueriesInfo[0].ItrExhausted, 	
			}
		}

		if len(kvrwset.MetadataWrites) != 0 {
			fmt.Println(" \n## KVMetadataWrite")
			fmt.Println(" Key = "+kvrwset.MetadataWrites[0].Key)
			fmt.Println(" Name = "+kvrwset.MetadataWrites[0].Entries[0].Name)
			metadataValue := CToGoString(kvrwset.MetadataWrites[0].Entries[1].Value[:])
			fmt.Println(" Value = "+metadataValue)

			kvMetadataWriteJson = KVMetadataWrite{
				Key: 	kvrwset.MetadataWrites[0].Key,
				Name:	kvrwset.MetadataWrites[0].Entries[0].Name,
			}
		}

		chaincodeKVRWSetJson := ChaincodeKVRWSet{
			Reads: 				kvReadJson,
			Writes: 			kvWriteJson,
			RangeQueriesInfo:	rangeQueryInfoJson,
			MetadataWrites:		kvMetadataWriteJson,
		}

		proposalResponsePayloadJson := ProposalResponsePayload{
			ProposalHash: 		hex.EncodeToString(proposalHash[:]),
			ChaincodeKVRWSet:   chaincodeKVRWSetJson,
		}

		chaincodeEndorsedActionJson := ChaincodeEndorsedAction {
			ProposalResponsePayload: proposalResponsePayloadJson,
		}

		transactionJson := Transaction {
			ChaincodeProposalPayload: 	chaincodeProposalPayloadJson,
			ChaincodeEndorsedAction:	chaincodeEndorsedActionJson,
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
