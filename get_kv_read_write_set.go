package main

import (	
	"github.com/hyperledger/fabric-protos-go/ledger/rwset/kvrwset"
	"github.com/hyperledger/fabric-protos-go/ledger/rwset"
	"github.com/hyperledger/fabric-protos-go/peer"
	"github.com/golang/protobuf/proto"
	"github.com/pkg/errors"
)

func GetKVRWSetJson(chaincodeAction *peer.ChaincodeAction)(ChaincodeKVRWSet, error){

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
		err := proto.Unmarshal(chaincodeAction.Results, txReadWriteSet)
		if err != nil {
			return ChaincodeKVRWSet{}, errors.WithMessage(err,"unmarshaling txReadWriteSet error: ")
		}

		RwSet := txReadWriteSet.NsRwset[0].Rwset

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
			return ChaincodeKVRWSet{}, errors.WithMessage(err,"unmarshaling kvrwset error: ")
		}

		var versionJson			Version
		var kvReadJson 			KVRead
		var kvWriteJson			KVWrite
		var rangeQueryInfoJson	RangeQueryInfo
		var kvMetadataWriteJson KVMetadataWrite

		if len(kvrwset.Reads) != 0 {
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
			kvWriteJson = KVWrite {
				Key: 		kvrwset.Writes[0].Key,
				IsDelete: 	kvrwset.Writes[0].IsDelete,
			}
		}

		if len(kvrwset.RangeQueriesInfo) != 0 {
			rangeQueryInfoJson = RangeQueryInfo{
				StartKey:  kvrwset.RangeQueriesInfo[0].StartKey,
				EndKey:	   kvrwset.RangeQueriesInfo[0].EndKey,
				ItrExhausted: kvrwset.RangeQueriesInfo[0].ItrExhausted, 	
			}
		}

		if len(kvrwset.MetadataWrites) != 0 {
			
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
	
	return chaincodeKVRWSetJson, nil
}
