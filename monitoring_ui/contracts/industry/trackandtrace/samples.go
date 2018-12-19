package main


	import (
		"github.com/hyperledger/fabric/core/chaincode/shim"
		iot "github.com/ibm-watson-iot/blockchain-samples/contracts/platform/iotcontractplatform"
)

var samples = `

{
    "API": {
        "deletePropertiesFromAssetSurgicalKit": {
            "args": [
                {
                    "qprops": [
                        "carpe noctem"
                    ],
                    "surgicalkit": {
                        "skitID": "A surgicalkit's ID"
                    }
                }
            ],
            "function": "deletePropertiesFromAssetSurgicalKit"
        },
        "readAssetStateHistorySurgicalKit": {
            "args": [
                {
                    "daterange": {
                        "begin": "timestamp formatted yyyy-mm-dd hh:mm:ss",
                        "end": "timestamp formatted yyyy-mm-dd hh:mm:ss"
                    },
                    "filter": {
                        "match": "all",
                        "select": [
                            {
                                "qprop": "Qualified property to compare, for example 'asset.assetID'",
                                "value": "Value to be compared"
                            }
                        ]
                    },
                    "surgicalkit": {
                        "skitID": "A surgicalkit's ID"
                    }
                }
            ],
            "function": "readAssetStateHistorySurgicalKit",
            "result": [
                {
                    "^CON": "INVALID OBJECT - MISSING PROPERTIES"
                }
            ]
        },
        "readAssetSurgicalKit": {
            "args": [
                {
                    "surgicalkit": {
                        "skitID": "A surgicalkit's ID"
                    }
                }
            ],
            "function": "readAssetSurgicalKit",
            "result": {
                "AssetKey": "This surgicalkit's world state surgicalkit ID",
                "alerts": [
                    "An alert name"
                ],
                "class": {},
                "compliant": true,
                "eventin": {
                    "surgicalkit": {
                        "common": {
                            "appdata": [
                                {
                                    "K": "carpe noctem",
                                    "V": "carpe noctem"
                                }
                            ],
                            "deviceID": "A unique identifier for the device that sent the current event",
                            "devicetimestamp": "A timestamp recoded by the device that sent the current event",
                            "location": {
                                "latitude": 123.456,
                                "longitude": 123.456
                            }
                        },
                        "hospital": {
                            "address": {
                                "city": "carpe noctem",
                                "country": "carpe noctem",
                                "postcode": "carpe noctem",
                                "streetandnumber": "carpe noctem"
                            },
                            "fence": {
                                "center": {
                                    "latitude": 123.456,
                                    "longitude": 123.456
                                },
                                "radius": 123.456
                            },
                            "name": "carpe noctem"
                        },
                        "sensors": {
                            "begin": "timestamp formatted yyyy-mm-dd hh:mm:ss",
                            "currtilt": 123.456,
                            "end": "timestamp formatted yyyy-mm-dd hh:mm:ss",
                            "endlocation": {
                                "latitude": 123.456,
                                "longitude": 123.456
                            },
                            "maxgforce": 123.456,
                            "maxtilt": 123.456,
                            "startlocation": {
                                "latitude": 123.456,
                                "longitude": 123.456
                            }
                        },
                        "skitID": "A surgicalkit's ID",
                        "status": "oem",
                        "transit": {
                            "begintransit": "timestamp formatted yyyy-mm-dd hh:mm:ss",
                            "carrier": "carpe noctem",
                            "endtransit": "timestamp formatted yyyy-mm-dd hh:mm:ss",
                            "receiver": "oem",
                            "shipper": "oem"
                        }
                    }
                },
                "eventout": {
                    "surgicalkit": {
                        "name": "The chaincode event's name",
                        "payload": {}
                    }
                },
                "state": {
                    "distanceFromCenter": 123.456,
                    "surgicalkit": {
                        "common": {
                            "appdata": [
                                {
                                    "K": "carpe noctem",
                                    "V": "carpe noctem"
                                }
                            ],
                            "deviceID": "A unique identifier for the device that sent the current event",
                            "devicetimestamp": "A timestamp recoded by the device that sent the current event",
                            "location": {
                                "latitude": 123.456,
                                "longitude": 123.456
                            }
                        },
                        "hospital": {
                            "address": {
                                "city": "carpe noctem",
                                "country": "carpe noctem",
                                "postcode": "carpe noctem",
                                "streetandnumber": "carpe noctem"
                            },
                            "fence": {
                                "center": {
                                    "latitude": 123.456,
                                    "longitude": 123.456
                                },
                                "radius": 123.456
                            },
                            "name": "carpe noctem"
                        },
                        "sensors": {
                            "begin": "timestamp formatted yyyy-mm-dd hh:mm:ss",
                            "currtilt": 123.456,
                            "end": "timestamp formatted yyyy-mm-dd hh:mm:ss",
                            "endlocation": {
                                "latitude": 123.456,
                                "longitude": 123.456
                            },
                            "maxgforce": 123.456,
                            "maxtilt": 123.456,
                            "startlocation": {
                                "latitude": 123.456,
                                "longitude": 123.456
                            }
                        },
                        "skitID": "A surgicalkit's ID",
                        "status": "oem",
                        "transit": {
                            "begintransit": "timestamp formatted yyyy-mm-dd hh:mm:ss",
                            "carrier": "carpe noctem",
                            "endtransit": "timestamp formatted yyyy-mm-dd hh:mm:ss",
                            "receiver": "oem",
                            "shipper": "oem"
                        }
                    }
                },
                "txnid": "Transaction UUID matching the blockchain",
                "txnts": "Transaction timestamp matching the blockchain"
            }
        }
    },
    "Model": {
        "invokeevent": {
            "name": "The chaincode event's name",
            "payload": {}
        },
        "ioteventcommon": {
            "appdata": [
                {
                    "K": "carpe noctem",
                    "V": "carpe noctem"
                }
            ],
            "deviceID": "A unique identifier for the device that sent the current event",
            "devicetimestamp": "A timestamp recoded by the device that sent the current event",
            "location": {
                "latitude": 123.456,
                "longitude": 123.456
            }
        },
        "stateFilter": {
            "match": "all",
            "select": [
                {
                    "qprop": "Qualified property to compare, for example 'asset.assetID'",
                    "value": "Value to be compared"
                }
            ]
        },
        "surgicalkit": {
            "common": {
                "appdata": [
                    {
                        "K": "carpe noctem",
                        "V": "carpe noctem"
                    }
                ],
                "deviceID": "A unique identifier for the device that sent the current event",
                "devicetimestamp": "A timestamp recoded by the device that sent the current event",
                "location": {
                    "latitude": 123.456,
                    "longitude": 123.456
                }
            },
            "hospital": {
                "address": {
                    "city": "carpe noctem",
                    "country": "carpe noctem",
                    "postcode": "carpe noctem",
                    "streetandnumber": "carpe noctem"
                },
                "fence": {
                    "center": {
                        "latitude": 123.456,
                        "longitude": 123.456
                    },
                    "radius": 123.456
                },
                "name": "carpe noctem"
            },
            "sensors": {
                "begin": "timestamp formatted yyyy-mm-dd hh:mm:ss",
                "currtilt": 123.456,
                "end": "timestamp formatted yyyy-mm-dd hh:mm:ss",
                "endlocation": {
                    "latitude": 123.456,
                    "longitude": 123.456
                },
                "maxgforce": 123.456,
                "maxtilt": 123.456,
                "startlocation": {
                    "latitude": 123.456,
                    "longitude": 123.456
                }
            },
            "skitID": "A surgicalkit's ID",
            "status": "oem",
            "transit": {
                "begintransit": "timestamp formatted yyyy-mm-dd hh:mm:ss",
                "carrier": "carpe noctem",
                "endtransit": "timestamp formatted yyyy-mm-dd hh:mm:ss",
                "receiver": "oem",
                "shipper": "oem"
            }
        },
        "surgicalkitstate": {
            "AssetKey": "This surgicalkit's world state surgicalkit ID",
            "alerts": [
                "An alert name"
            ],
            "class": {},
            "compliant": true,
            "eventin": {
                "surgicalkit": {
                    "common": {
                        "appdata": [
                            {
                                "K": "carpe noctem",
                                "V": "carpe noctem"
                            }
                        ],
                        "deviceID": "A unique identifier for the device that sent the current event",
                        "devicetimestamp": "A timestamp recoded by the device that sent the current event",
                        "location": {
                            "latitude": 123.456,
                            "longitude": 123.456
                        }
                    },
                    "hospital": {
                        "address": {
                            "city": "carpe noctem",
                            "country": "carpe noctem",
                            "postcode": "carpe noctem",
                            "streetandnumber": "carpe noctem"
                        },
                        "fence": {
                            "center": {
                                "latitude": 123.456,
                                "longitude": 123.456
                            },
                            "radius": 123.456
                        },
                        "name": "carpe noctem"
                    },
                    "sensors": {
                        "begin": "timestamp formatted yyyy-mm-dd hh:mm:ss",
                        "currtilt": 123.456,
                        "end": "timestamp formatted yyyy-mm-dd hh:mm:ss",
                        "endlocation": {
                            "latitude": 123.456,
                            "longitude": 123.456
                        },
                        "maxgforce": 123.456,
                        "maxtilt": 123.456,
                        "startlocation": {
                            "latitude": 123.456,
                            "longitude": 123.456
                        }
                    },
                    "skitID": "A surgicalkit's ID",
                    "status": "oem",
                    "transit": {
                        "begintransit": "timestamp formatted yyyy-mm-dd hh:mm:ss",
                        "carrier": "carpe noctem",
                        "endtransit": "timestamp formatted yyyy-mm-dd hh:mm:ss",
                        "receiver": "oem",
                        "shipper": "oem"
                    }
                }
            },
            "eventout": {
                "surgicalkit": {
                    "name": "The chaincode event's name",
                    "payload": {}
                }
            },
            "state": {
                "distanceFromCenter": 123.456,
                "surgicalkit": {
                    "common": {
                        "appdata": [
                            {
                                "K": "carpe noctem",
                                "V": "carpe noctem"
                            }
                        ],
                        "deviceID": "A unique identifier for the device that sent the current event",
                        "devicetimestamp": "A timestamp recoded by the device that sent the current event",
                        "location": {
                            "latitude": 123.456,
                            "longitude": 123.456
                        }
                    },
                    "hospital": {
                        "address": {
                            "city": "carpe noctem",
                            "country": "carpe noctem",
                            "postcode": "carpe noctem",
                            "streetandnumber": "carpe noctem"
                        },
                        "fence": {
                            "center": {
                                "latitude": 123.456,
                                "longitude": 123.456
                            },
                            "radius": 123.456
                        },
                        "name": "carpe noctem"
                    },
                    "sensors": {
                        "begin": "timestamp formatted yyyy-mm-dd hh:mm:ss",
                        "currtilt": 123.456,
                        "end": "timestamp formatted yyyy-mm-dd hh:mm:ss",
                        "endlocation": {
                            "latitude": 123.456,
                            "longitude": 123.456
                        },
                        "maxgforce": 123.456,
                        "maxtilt": 123.456,
                        "startlocation": {
                            "latitude": 123.456,
                            "longitude": 123.456
                        }
                    },
                    "skitID": "A surgicalkit's ID",
                    "status": "oem",
                    "transit": {
                        "begintransit": "timestamp formatted yyyy-mm-dd hh:mm:ss",
                        "carrier": "carpe noctem",
                        "endtransit": "timestamp formatted yyyy-mm-dd hh:mm:ss",
                        "receiver": "oem",
                        "shipper": "oem"
                    }
                }
            },
            "txnid": "Transaction UUID matching the blockchain",
            "txnts": "Transaction timestamp matching the blockchain"
        },
        "surgicalkitstatearray": [
            {
                "^CON": "INVALID OBJECT - MISSING PROPERTIES"
            }
        ],
        "surgicalkitstateexternal": {
            "^CON": "INVALID OBJECT - MISSING PROPERTIES"
        }
    }
}`


	var readAssetSamples iot.ChaincodeFunc = func(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
		return []byte(samples), nil
	}

	func init() {
		iot.AddRoute("readAssetSamples", "query", iot.SystemClass, readAssetSamples)
	}
	