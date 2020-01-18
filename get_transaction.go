package main

import (
	"crypto/sha256"
	"encoding/hex"
	"github.com/hyperledger/fabric-protos-go/peer"
	"github.com/golang/protobuf/proto"
	"github.com/pkg/errors"
)

func GetTransactionJson(chaincodeActionPayload *peer.ChaincodeActionPayload) (Transaction, error){

	/*
		type TransactionAction struct {
			Header  []byte 
			Payload []byte 
		}
	*/

	// Payload field is marshalled object of ChaincodeActionPayload

	/*
		type ChaincodeActionPayload struct {
			ChaincodeProposalPayload []byte 
			Action *ChaincodeEndorsedAction 
		}
	*/

	// 1. ChaincodeProposalPayload

	chaincodeProposalPayloadJson, err := GetChainCodeProposalPayload(chaincodeActionPayload)
	if err != nil {
		return Transaction{}, errors.WithMessage(err,"chaincode proposal payload error: ")
	}

	// 2. ChaincodeEndorsedAction
	chaincodeEndorsedActionJson, err := GetChainCodeEndorsedAction(chaincodeActionPayload)
	if err != nil {
		return Transaction{}, errors.WithMessage(err,"chaincode endorse action error: ")
	}

	transactionJson := Transaction {
		ChaincodeProposalPayload: 	chaincodeProposalPayloadJson,
		ChaincodeEndorsedAction:	chaincodeEndorsedActionJson,
	}

	return transactionJson, nil
}

func GetChainCodeProposalPayload(chaincodeActionPayload *peer.ChaincodeActionPayload) (ChaincodeProposalPayload, error){

	/*
		type ChaincodeProposalPayload struct {
			Input        []byte 
			TransientMap map[string][]byte 
		}
	*/

	chaincodeProposalPayload := &peer.ChaincodeProposalPayload{}
	err := proto.Unmarshal(chaincodeActionPayload.ChaincodeProposalPayload, chaincodeProposalPayload)
		if err != nil {
			return ChaincodeProposalPayload{}, errors.WithMessage(err,"unmarshaling Chaincode Proposal Payload error: ")
		}

		// The Input field is marshalled object of ChaincodeInvocationSpec
		input := &peer.ChaincodeInvocationSpec{}
		err = proto.Unmarshal(chaincodeProposalPayload.Input, input)
		if err != nil {
			return ChaincodeProposalPayload{}, errors.WithMessage(err,"unmarshaling Chaincode Proposal Payload Input error: ")
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

		chaincodeArgs := make([]string, len(input.ChaincodeSpec.Input.Args))

		for i, c := range input.ChaincodeSpec.Input.Args {
			args := CToGoString(c[:])
			chaincodeArgs[i] = args
		}
	
		ChaincodeType := [5]string{"UNDEFINED", "GOLANG", "NODE", "CAR", "JAVA"}
	
		chaincodeSpecJson := ChaincodeSpec{
			ChaincodeId: 	input.ChaincodeSpec.ChaincodeId.Name,
			ChaincodeType:  ChaincodeType[input.ChaincodeSpec.Type],
			ChaincodeArgs:	chaincodeArgs,
		}

		chaincodeInvocationSpecJson := ChaincodeInvocationSpec{
			ChaincodeSpec:  chaincodeSpecJson,
		}

		chaincodeProposalPayloadJson := ChaincodeProposalPayload {
			ChaincodeInvocationSpec: chaincodeInvocationSpecJson,
		}

	return chaincodeProposalPayloadJson, nil
}


func GetChainCodeEndorsedAction(chaincodeActionPayload *peer.ChaincodeActionPayload) (ChaincodeEndorsedAction, error){

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
		err := proto.Unmarshal(chaincodeActionPayload.Action.ProposalResponsePayload, proposalResponsePayload)
		if err != nil {
			return ChaincodeEndorsedAction{}, errors.WithMessage(err,"unmarshaling Proposal Response Payload error: ")
		}

		proposalHash := sha256.Sum256(proposalResponsePayload.ProposalHash)
		
		chaincodeKVRWSetJson, err := GetKVRWSetJson(proposalResponsePayload)
		if err != nil {
			return ChaincodeEndorsedAction{}, errors.WithMessage(err,"failed to get KVRWSet Json error: ")
		} 

		proposalResponsePayloadJson := ProposalResponsePayload{
			ProposalHash: 		hex.EncodeToString(proposalHash[:]),
			ChaincodeKVRWSet:   chaincodeKVRWSetJson,
		}

		chaincodeEndorsedActionJson := ChaincodeEndorsedAction {
			ProposalResponsePayload: proposalResponsePayloadJson,
		}

	return chaincodeEndorsedActionJson, nil
}
