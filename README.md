<h1>BlockReader</h1>
<p>     <img alt="GitHub last commit" src="https://img.shields.io/github/last-commit/Deeptiman/blockreader">  <img alt="GitHub language count" src="https://img.shields.io/github/languages/count/Deeptiman/blockreader"> <img alt="GitHub top language" src="https://img.shields.io/github/languages/top/Deeptiman/blockreader"></p>
<p><a href="https://www.hyperledger.org/projects/fabric"><img src="https://www.hyperledger.org/wp-content/uploads/2016/09/logo_hl_new.png" alt="N|Solid"></a></p>
<p><b>BlockReader</b> application extract and showcase the complete data structure of a Block that contains several details for a transaction. The application will require a transaction id to query the ledger to retrieve the associated block. Then the application will follow the Block data structure to read the content of the Block.</p>
<p><b>Medium writeup : </b><a href="https://medium.com/@deeptiman/whats-inside-the-block-in-hyperledger-fabric-69a0934fef08">https://medium.com/@deeptiman/whats-inside-the-block-in-hyperledger-fabric-69a0934fef08</a>


<h2> Run the application </h2>
<ol>
  <li> 
    Compile the application using following commands 
                
    cd /go/src/github.com/blockreader
    go build            
   
 You can see <b>blockreader</b> executable is generated.
  </li>
 
  <li>
  
   Now type following commands to read the Block.
  
   
                
     ./blockreader -txnId="cd2b072c880cdefbea66c5f9d73a5a5eb3c3977e77772fba42cec59204ca2980" -channelId="employeeledger"            
     -txnId : The transaction id to be query to retrieve the Block for that transaction from the ledger.
     -channelId: The network ChannelId
                
  
  </li>
  
  <li>
   After that, you will see the Block details in a JSON format.
       
    ************* BLOCK READER JSON ************* 
    {
        "BlockHeader": {
            "number": 2,
            "previous_hash": "2e8ddbf2dfd6b90fdadd0f2653f932f9f8bb7922244f991b66fc9ebeff3c63c1",
            "data_hash": "6df36db3c31ac485754d7e840574bf53639539d9b810c433cec143fe234557e8"
        },
        "BlockData": {
            "envelope": {
                "header": {
                    "payload": {
                        "channel_header": {
                            "type": "ENDORSER_TRANSACTION",
                            "version": 0,
                            "channelid": "employeeledger",
                            "txid": "c3d95a4e5606cae3a1aabced9fe532f7f628b2930eb09452f2b0db62f6e6ee0b",
                            "epoch": 0,
                            "extension": {
                                "chaincode_id": {
                                    "path": "",
                                    "name": "employeeledger",
                                    "version": ""
                                }
                            }
                        },
                        "signature_header": {
                            "creator": {
                                "msp_id": "org1.employee.ledger.com",
                                "cert_hash": "-----BEGIN CERTIFICATE-----\nMIICyjCCAnGgAwIBAgIUGjzRUXjNaP8bpU/9B2xjpDRPWRYwCgYIKoZIzj0EAwIw\ngYMxCzAJBgNVBAYTAlVTMRMwEQYDVQQIEwpDYWxpZm9ybmlhMRYwFAYDVQQHEw1T\nYW4gRnJhbmNpc2NvMSEwHwYDVQQKExhvcmcxLmVtcGxveWVlLmxlZGdlci5jb20x\nJDAiBgNVBAMTG2NhLm9yZzEuZW1wbG95ZWUubGVkZ2VyLmNvbTAeFw0yMDAxMTky\nMDE3MDBaFw0yMTAxMTgyMDIyMDBaMDYxGjALBgNVBAsTBHVzZXIwCwYDVQQLEwRv\ncmcxMRgwFgYDVQQDDA90ZXN0MUBnbWFpbC5jb20wWTATBgcqhkjOPQIBBggqhkjO\nPQMBBwNCAARifA6/HyRpmvTk6oNRCeB8QRWgysaxIxaTd36YHgTqFWlSA0oUE2PI\n2HZ7EiY/AOixTtQGzrAoxYgwpz219LBLo4IBDTCCAQkwDgYDVR0PAQH/BAQDAgeA\nMAwGA1UdEwEB/wQCMAAwHQYDVR0OBBYEFEvpcBUx2KVrXK3HJ/d7wLbYOsWRMCsG\nA1UdIwQkMCKAIGDnIjkscY/5fSt5a+QeZtR7sLnSYbI6t10GiS1huLBJMCQGA1Ud\nEQQdMBuCGWRlZXB0aW1hbnBjLUxlbm92by1HNTAtNDUwdwYIKgMEBQYHCAEEa3si\nYXR0cnMiOnsiaGYuQWZmaWxpYXRpb24iOiJvcmcxIiwiaGYuRW5yb2xsbWVudElE\nIjoidGVzdDFAZ21haWwuY29tIiwiaGYuVHlwZSI6InVzZXIiLCJ1c2VybW9kZSI6\nIkFkbWluIn19MAoGCCqGSM49BAMCA0cAMEQCIDUqcBCkezNYkoOpXzJb5F7ZMsaF\nkamc0bRNCLZy4nlSAiBYodLyCs7iivVonEn49L5GaH7c6zUmQw1hcL9AoGttVw==\n-----END CERTIFICATE-----\n",
                                "certificate": {
                                    "country": [
                                        "US"
                                    ],
                                    "organization": [
                                        "org1.employee.ledger.com"
                                    ],
                                    "organization_unit": null,
                                    "locality": [
                                        "San Francisco"
                                    ],
                                    "province": [
                                        "California"
                                    ],
                                    "serial_number": "",
                                    "not_before": "2020-01-19T20:17:00Z",
                                    "not_after": "2021-01-18T20:22:00Z"
                                }
                            }
                        }
                    }
                },
                "data": {
                    "transaction": {
                        "chaincode_proposal_payload": {
                            "chaincode_invocation_spec": {
                                "chaincode_spec": {
                                    "chaincode_id": "employeeledger",
                                    "chaincode_type": "GOLANG",
                                    "chaincode_args": [
                                        "invoke",
                                        "createUser",
                                        "Deeptiman Pattnaik",
                                        "test1@gmail.com",
                                        "Personal",
                                        "Software",
                                        "90,000",
                                        "Admin"
                                    ]
                                }
                            }
                        },
                        "chaincode_endorsed_action": {
                            "proposal_response_payload": {
                                "proposal_hash": "0597a5c323fa68d26e3063ccf640f8658c19f80454b573866cf0614fc2d57f11",
                                "chaincode_kv_rw_set": {
                                    "reads": {
                                        "Key": "employeeledger",
                                        "version": {
                                            "block_num": 1,
                                            "txnum": 0
                                        }
                                    },
                                    "range_queries_info": {
                                        "startkey": "",
                                        "endkey": "",
                                        "itr_exhausted": false
                                    },
                                    "writes": {
                                        "key": "",
                                        "is_delete": false
                                    },
                                    "metadata_writes": {
                                        "key": "",
                                        "name": ""
                                    }
                                },
                                "chaincode_events": {
                                    "chaincode_id": "employeeledger",
                                    "txid": "c3d95a4e5606cae3a1aabced9fe532f7f628b2930eb09452f2b0db62f6e6ee0b",
                                    "event_name": "addUserInvoke",
                                    "payload": ""
                                }
                            }
                        }
                    }
                }
            }
        },
        "BlockMetaData": {
            "value": null,
            "signature": "MEQCICshHeJiB1d1SdBmOyVaQLBrgmAOsrxwEEUO3ZkbB/gDAiAuGcyGDC/i8ExrfZg0p5zzSP/HnBODqSOe5kucVA0ZQg==",
            "signature_header": {
                "creator": {
                    "msp_id": "employee.ledger.com",
                    "cert_hash": "-----BEGIN CERTIFICATE-----\nMIICJDCCAcqgAwIBAgIQCNM0iKIvFrsC0jQGmS98nzAKBggqhkjOPQQDAjB5MQsw\nCQYDVQQGEwJVUzETMBEGA1UECBMKQ2FsaWZvcm5pYTEWMBQGA1UEBxMNU2FuIEZy\nYW5jaXNjbzEcMBoGA1UEChMTZW1wbG95ZWUubGVkZ2VyLmNvbTEfMB0GA1UEAxMW\nY2EuZW1wbG95ZWUubGVkZ2VyLmNvbTAeFw0yMDAxMTgwODA0MDFaFw0zMDAxMTUw\nODA0MDFaMGAxCzAJBgNVBAYTAlVTMRMwEQYDVQQIEwpDYWxpZm9ybmlhMRYwFAYD\nVQQHEw1TYW4gRnJhbmNpc2NvMSQwIgYDVQQDExtvcmRlcmVyLmVtcGxveWVlLmxl\nZGdlci5jb20wWTATBgcqhkjOPQIBBggqhkjOPQMBBwNCAASmE81Sa568qPkJHZR5\nsHtTpF7cLqyhNiHg8Qroq1xlZrdUNxBhIKpT1KMEWA8N1e2sr8HgTLZCL+6s/AHp\nKUXRo00wSzAOBgNVHQ8BAf8EBAMCB4AwDAYDVR0TAQH/BAIwADArBgNVHSMEJDAi\ngCDdKLH7Qvi6j2qO/O4Qb9iZHSAxS1KcLt/XePcoVfHz4jAKBggqhkjOPQQDAgNI\nADBFAiEA2i5wlHfoVlpPKVTinqcO6a9mqmWoPVf67f6V9XV16kICIDhgn18/hRXD\n7PYTnIUfx6l7ruLU+LuJywzxkKD0AScR\n-----END CERTIFICATE-----\n",
                    "certificate": {
                        "country": [
                            "US"
                        ],
                        "organization": [
                            "employee.ledger.com"
                        ],
                        "organization_unit": null,
                        "locality": [
                            "San Francisco"
                        ],
                        "province": [
                            "California"
                        ],
                        "serial_number": "",
                        "not_before": "2020-01-18T08:04:01Z",
                        "not_after": "2030-01-15T08:04:01Z"
                    }
                }
            }
        }
    }
  </li>
<h2>License</h2>
<p>This project is licensed under the <a href="https://github.com/Deeptiman/blockreader/blob/master/LICENSE">MIT License</a></p>
