package objects

var FirmwareUpdateDescriptor = `{
      "Id": 5,
      "Name": "Firmware Update",
      "Multiple": false,
      "Mandatory": false,
      "Version": "1.1",
      "LwM2MVersion": "1.1",
      "URN": "urn:oma:lwm2m:oma:5:1.1",
      "Resources": [
        {
          "Id": 0,
          "Name": "Package",
          "Operations": "W",
          "Multiple": false,
          "Mandatory": true,
          "ResourceType": "opaque"
        },
        {
          "Id": 1,
          "Name": "Package URI",
          "Operations": "W",
          "Multiple": false,
          "Mandatory": true,
          "ResourceType": "string",
          "RangeOrEnums": "0-255 bytes",
          "ValueValidator": "NewRangeValidator(0 255)"
        },
        {
          "Id": 2,
          "Name": "Update",
          "Operations": "E",
          "Multiple": false,
          "Mandatory": true,
          "ResourceType": "string"
        },
        {
          "Id": 3,
          "Name": "State",
          "Operations": "R",
          "Multiple": false,
          "Mandatory": true,
          "ResourceType": "int",
          "RangeOrEnums": "1-3",
          "ValueValidator": "NewRangeValidator(1 3)"
        },
        {
          "Id": 4,
          "Name": "Update Supported Objects",
          "Operations": "RW",
          "Multiple": false,
          "Mandatory": false,
          "ResourceType": "bool"
        },
        {
          "Id": 5,
          "Name": "Update Result",
          "Operations": "R",
          "Multiple": false,
          "Mandatory": true,
          "ResourceType": "int",
          "RangeOrEnums": "0-6",
          "ValueValidator": "NewRangeValidator(0 6)"
        }
      ]
    }
`
