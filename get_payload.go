package main

import (
	"encoding/base64"
	"crypto/x509"
	"encoding/pem"
	"github.com/hyperledger/fabric-protos-go/msp"
	"github.com/hyperledger/fabric-protos-go/peer"
	"github.com/hyperledger/fabric-protos-go/common"
	"github.com/golang/protobuf/proto"
	"github.com/pkg/errors"
)


func GetPayloadJson(envelope *common.Envelope) (Payload, error){

 	//Retrieve the Payload from the Envelope
	/*
		type Payload struct {
			Header *Header 
			Data   []byte 
		}
	*/

	payload := &common.Payload{}
	err := proto.Unmarshal(envelope.Payload, payload)
	if err != nil {
		return Payload{}, errors.WithMessage(err,"unmarshaling Payload error: ")
	}

	/*
		type Header struct {
			ChannelHeader   []byte 
			SignatureHeader []byte 
		}
	*/

	// 1. ChannelHeader
		channelHeaderJson, err := GetChannelHeader(payload)	
		if err != nil {
			return Payload{},  errors.WithMessage(err, "failed get Channel Header")
		}

	// 2. SignatureHeader
		signatureHeaderJson, err := GetSignatureHeader(payload)
		if err != nil {
			return Payload{}, errors.WithMessage(err, "failed get Signature Header")
		}

	payloadJson := Payload {
		ChannelHeader: 		channelHeaderJson,
		SignatureHeader:	signatureHeaderJson,
	}

	return payloadJson, nil
}


func GetChannelHeader(payload *common.Payload) (ChannelHeader, error) {

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
		err := proto.Unmarshal(payload.Header.ChannelHeader, channelHeader)
		if err != nil {
			return ChannelHeader{}, errors.WithMessage(err,"unmarshaling Channel Header error: ")
		}
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
		extension := &peer.ChaincodeHeaderExtension{}
		err = proto.Unmarshal(channelHeader.Extension, extension)
		if err != nil {
			return ChannelHeader{}, errors.WithMessage(err,"unmarshaling Channel Header error: ")
		}

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

	return channelHeaderJson, nil
}

func GetSignatureHeader(payload *common.Payload) (SignatureHeader, error){

	/*
		type SignatureHeader struct {
			Creator []byte
			Nonce   []byte
		}
	*/

		// Creator is the marshalled object of msp.SerializedIdentity
		signatureHeader := &common.SignatureHeader{}
		err := proto.Unmarshal(payload.Header.SignatureHeader, signatureHeader)
		if err != nil {
			return SignatureHeader{}, errors.WithMessage(err,"unmarshaling Signature Header error: ")
		}

		creator := &msp.SerializedIdentity{}
		err = proto.Unmarshal(signatureHeader.Creator, creator)
		if err != nil {
			return SignatureHeader{}, errors.WithMessage(err,"unmarshaling Creator error: ")
		}

		uEnc := base64.URLEncoding.EncodeToString([]byte(creator.IdBytes))

		// Base64 Url Decoding
		certText, err := base64.URLEncoding.DecodeString(uEnc)
		if err != nil {
			return SignatureHeader{}, errors.WithMessage(err,"Error decoding string: ")
		}
		
		end, _ := pem.Decode([]byte(string(certText)))
		if end == nil {
			return SignatureHeader{}, errors.New("Error Pem decoding: ")
		}
		cert, err := x509.ParseCertificate(end.Bytes)
		if err != nil {
			return SignatureHeader{}, errors.New("failed to parse certificate:: ")
		}		

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

	return signatureHeaderJson, nil
}
