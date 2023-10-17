package preset

var connStatsDescriptor = `{
      "Id": 7,
      "Name": "Connectivity Statistics",
      "Multiple": false,
      "Mandatory": false,
      "Version": "1.0",
      "LwM2MVersion": "1.0",
      "URN": "urn:oma:lwm2m:oma:7",
      "Resources": [
        {
          "Id": 0,
          "Name": "SMS Tx Counter",
          "Operations": "R",
          "Multiple": false,
          "Mandatory": false,
          "ResourceType": "int"
        },
        {
          "Id": 1,
          "Name": "SMS Rx Counter",
          "Operations": "R",
          "Multiple": false,
          "Mandatory": false,
          "ResourceType": "int"
        },
        {
          "Id": 2,
          "Name": "Tx Data",
          "Operations": "R",
          "Multiple": false,
          "Mandatory": false,
          "ResourceType": "int",
          "Units": "Kilo-Bytes"
        },
        {
          "Id": 3,
          "Name": "Rx Data",
          "Operations": "R",
          "Multiple": false,
          "Mandatory": false,
          "ResourceType": "int",
          "Units": "Kilo-Bytes"
        },
        {
          "Id": 4,
          "Name": "Max Message Size",
          "Operations": "R",
          "Multiple": false,
          "Mandatory": false,
          "ResourceType": "int",
          "Units": "Byte"
        },
        {
          "Id": 5,
          "Name": "Average Message Size",
          "Operations": "R",
          "Multiple": false,
          "Mandatory": false,
          "ResourceType": "int",
          "Units": "Byte"
        },
        {
          "Id": 6,
          "Name": "StartOrReset",
          "Operations": "E",
          "Multiple": false,
          "Mandatory": true,
          "ResourceType": "string"
        }
      ]
    }
`
