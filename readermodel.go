package main

import (
	"time"
)

type Block struct {

	BlockHeader 	BlockHeader 	`json:block_header`
	BlockData	BlockData	`json:block_data`
	BlockMetaData 	BlockMetaData 	`json:block_metadata`
}

type BlockHeader struct {
	Number   		uint64	`json:"number"`
	PreviousHash    	string	`json:"previous_hash"`	
	DataHash        	string	`json:"data_hash"`
}

type BlockData struct {

	Envelope Envelope	`json:"envelope"`
}

type BlockMetaData struct {
	Value			[]byte			 `json:"value"`
	Signature		[]byte			 `json:"signature"`
	SignatureHeader 	SignatureHeader  	 `json:"signature_header"`
}

type Envelope struct {

	Header 	Header 	`json:"header"`
	Data	Data	`json:"data"`
}

type Header struct {
	Payload	 Payload `json:"payload"`
}

type Payload struct {

	ChannelHeader 	ChannelHeader 	 `json:"channel_header"`
	SignatureHeader SignatureHeader  `json:"signature_header"`
}

type ChannelHeader struct {

	Type      string 		`json:"type"`
	Version   int32 		`json:"version"`
	ChannelId string 		`json:"channelid"`
	TxId      string 		`json:"txid"`
	Epoch     uint64 		`json:"epoch"`
	Extension ChaincodeHeaderExtension 	`json:"extension"`
}

type ChaincodeHeaderExtension struct {
	ChaincodeId  ChaincodeID  `json:"chaincode_id"`
}

type ChaincodeID struct {
	Path    	string 		`json:"path"`
	Name    	string 		`json:"name"`
	Version 	string		`json:"version"`
}

type SignatureHeader struct {
	Creator 	Creator  	`json:"creator"`
}

type Creator struct {
	Mspid 		string 			`json:"msp_id"`
	CertHash	string			`json:"cert_hash"`
	Certificate 	Certificate		`json:"certificate"`
}

type Certificate struct {

	Country  			[]string 		`json:"country"`
	Organization			[]string		`json:"organization"`
	OrganizationalUnit		[]string		`json:"organization_unit"`
	Locality			[]string		`json:"locality"`
	Province			[]string		`json:"province"`
	SerialNumber			string			`json:"serial_number"`
	NotBefore			time.Time		`json:"not_before"`	 
	NotAfter 			time.Time		`json:"not_after"`
}

type Data struct {

	Transaction Transaction `json:"transaction"`
}

type Transaction struct {

	ChaincodeProposalPayload ChaincodeProposalPayload `json:"chaincode_proposal_payload"`
	ChaincodeEndorsedAction	 ChaincodeEndorsedAction  `json:"chaincode_endorsed_action"`
}

type ChaincodeProposalPayload struct {

	ChaincodeInvocationSpec ChaincodeInvocationSpec `json:"chaincode_invocation_spec"`
}

type ChaincodeInvocationSpec struct {
	ChaincodeSpec ChaincodeSpec `json:"chaincode_spec"`
}

type ChaincodeSpec struct {

	ChaincodeId    string    `json:"chaincode_id"`
	ChaincodeType  string	 `json:"chaincode_type"`	
	ChaincodeArgs  []string  `json:"chaincode_args"`
}


type ChaincodeEndorsedAction struct {
	ProposalResponsePayload ProposalResponsePayload `json:"proposal_response_payload"`
}

type ProposalResponsePayload struct {
	ProposalHash   		string 			`json:"proposal_hash"`
	ChaincodeKVRWSet	ChaincodeKVRWSet	`json:"chaincode_kv_rw_set"`
	ChaincodeEvents		ChaincodeEvents		`json:"chaincode_events"`
}

type ChaincodeEvents struct {
	ChaincodeId     	string  `json:"chaincode_id"`
	TxId			string	`json:"txid"`
	EventName		string	`json:"event_name"`
	Payload			string 	`json:"payload"`	
}

type ChaincodeKVRWSet struct {

	Reads            KVRead         	`json:"reads"`
	RangeQueriesInfo RangeQueryInfo 	`json:"range_queries_info"`
	Writes           KVWrite  		`json:"writes"`
	MetadataWrites   KVMetadataWrite	`json:"metadata_writes"`
}

type KVRead struct {
	Key     	string   	`json:"key`
	Version 	Version 	`json:"version"`
}

type Version struct {
	BlockNum 	uint64 		`json:"block_num"`
	TxNum    	uint64 		`json:"txnum"`
}

type KVWrite struct {
	Key      	string 		`json:"key"`
	IsDelete 	bool  		`json:"is_delete"`
}

type RangeQueryInfo struct {
	StartKey     	string 	`json:"startkey"`
	EndKey       	string 	`json:"endkey"`
	ItrExhausted 	bool 	`json:"itr_exhausted"`
}

type KVMetadataWrite struct {
	Key     string		`json:"key"`
	Name	string		`json:"name"`
}



