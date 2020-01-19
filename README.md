<h1>BlockReader</h1>
<p><a href="https://www.hyperledger.org/projects/fabric"><img src="https://www.hyperledger.org/wp-content/uploads/2016/09/logo_hl_new.png" alt="N|Solid"></a></p>
<p><b>BlockReader</b> application extract and showcase the complete data structure of a Block that contains several details for a transaction. The application will require a transaction id to query the ledger to retrieve the associated block. Then the application will follow the Block data structure to read the content of the Block.</p>

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
                        "Number": 2,
                        "PreviousHash": "ff4c6159dd0d4c3f755c5e9963588e86362d8852918490c200d636dd1b60c407",
                        "DataHash": "64d10332350b4d9f9baf99093656b364b6319f40b6efe7e7240965f909c2b6bf"
                    },
                    "BlockData": {
                        "Envelope": {
                            "Header": {
                                "payload": {
                                    "ChannelHeader": {
                                        "Type": "ENDORSER_TRANSACTION",
                                        "Version": 0,
                                        "ChannelId": "employeeledger",
                                        "TxId": "cd2b072c880cdefbea66c5f9d73a5a5eb3c3977e77772fba42cec59204ca2980",
                                        "Epoch": 0,
                                        "Extension": {
                                            "ChaincodeId": {
                                                "Path": "",
                                                "Name": "employeeledger",
                                                "Version": ""
                                            }
                                        }
                                    },
                                    "SignatureHeader": {
                                        "Creator": {
                                            "Mspid": "org1.employee.ledger.com",
                                            "CertText": "-----BEGIN CERTIFICATE-----\nMIICyzCCAnGgAwIBAgIUVUY43C+/a4LxPGaMCuYPtKHGYLAwCgYIKoZIzj0EAwIw\ngYMxCzAJBgNVBAYTAlVTMRMwEQYDVQQIEwpDYWxpZm9ybmlhMRYwFAYDVQQHEw1T\nYW4gRnJhbmNpc2NvMSEwHwYDVQQKExhvcmcxLmVtcGxveWVlLmxlZGdlci5jb20x\nJDAiBgNVBAMTG2NhLm9yZzEuZW1wbG95ZWUubGVkZ2VyLmNvbTAeFw0yMDAxMTgx\nNTA2MDBaFw0yMTAxMTcxNTExMDBaMDYxGjALBgNVBAsTBHVzZXIwCwYDVQQLEwRv\ncmcxMRgwFgYDVQQDDA90ZXN0MUBnbWFpbC5jb20wWTATBgcqhkjOPQIBBggqhkjO\nPQMBBwNCAAQNfOKAY+2J9OdI8eAoLUeBArn8VnSIA8ElkzJdkHirDp0IRBc+j/4L\nTH/bgQnwsgNwijoWRdlCK+ZfNWcHolgyo4IBDTCCAQkwDgYDVR0PAQH/BAQDAgeA\nMAwGA1UdEwEB/wQCMAAwHQYDVR0OBBYEFGoRvg3MuNSfoAYtf/USag2NssxzMCsG\nA1UdIwQkMCKAIGDnIjkscY/5fSt5a+QeZtR7sLnSYbI6t10GiS1huLBJMCQGA1Ud\nEQQdMBuCGWRlZXB0aW1hbnBjLUxlbm92by1HNTAtNDUwdwYIKgMEBQYHCAEEa3si\nYXR0cnMiOnsiaGYuQWZmaWxpYXRpb24iOiJvcmcxIiwiaGYuRW5yb2xsbWVudElE\nIjoidGVzdDFAZ21haWwuY29tIiwiaGYuVHlwZSI6InVzZXIiLCJ1c2VybW9kZSI6\nIkFkbWluIn19MAoGCCqGSM49BAMCA0gAMEUCIQDzM4CaWqaux+Mko/iovrqOHkQS\noQqgkg8t+xaA7kirZwIgTst30Yee+IqzHGSbl7f07M/d3yOX3mvZsa1DFk3HoVI=\n-----END CERTIFICATE-----\n",
                                            "Certificate": {
                                                "Country": [
                                                    "US"
                                                ],
                                                "Organization": [
                                                    "org1.employee.ledger.com"
                                                ],
                                                "OrganizationalUnit": null,
                                                "Locality": [
                                                    "San Francisco"
                                                ],
                                                "Province": [
                                                    "California"
                                                ],
                                                "SerialNumber": "",
                                                "NotBefore": "2020-01-18T15:06:00Z",
                                                "NotAfter": "2021-01-17T15:11:00Z"
                                            }
                                        }
                                    }
                                }
                            },
                            "Data": {
                                "Transaction": {
                                    "ChaincodeProposalPayload": {
                                        "ChaincodeInvocationSpec": {
                                            "ChaincodeSpec": {
                                                "ChaincodeId": "employeeledger",
                                                "ChaincodeType": "GOLANG",
                                                "ChaincodeArgs": [
                                                    "invoke",
                                                    "createUser",
                                                    "Putul",
                                                    "test1@gmail.com",
                                                    "Personal",
                                                    "Software",
                                                    "87",
                                                    "Admin"
                                                ]
                                            }
                                        }
                                    },
                                    "ChaincodeEndorsedAction": {
                                        "ProposalResponsePayload": {
                                            "ProposalHash": "92eb219beb95fe92b1164d23d4951377b1d4161b34e491055b0cc4b66bb57d3e",
                                            "ChaincodeKVRWSet": {
                                                "Reads": {
                                                    "Key": "employeeledger",
                                                    "Version": {
                                                        "BlockNum": 1,
                                                        "TxNum": 0
                                                    }
                                                },
                                                "RangeQueriesInfo": {
                                                    "StartKey": "",
                                                    "EndKey": "",
                                                    "ItrExhausted": false
                                                },
                                                "Writes": {
                                                    "Key": "",
                                                    "IsDelete": false
                                                },
                                                "MetadataWrites": {
                                                    "Key": "",
                                                    "Name": ""
                                                }
                                            },
                                            "ChaincodeEvents": {
                                                "ChaincodeId": "employeeledger",
                                                "TxId": "cd2b072c880cdefbea66c5f9d73a5a5eb3c3977e77772fba42cec59204ca2980",
                                                "EventName": "addUserInvoke",
                                                "Payload": ""
                                            }
                                        }
                                    }
                                }
                            }
                        }
                    }
                }
        
  
  </li>
  
