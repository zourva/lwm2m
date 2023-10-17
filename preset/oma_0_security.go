package preset

var securityDescriptor = `{
      "Id": 0,
      "Name": "LwM2M Security",
      "Multiple": true,
      "Mandatory": true,
      "Version": "1.2",
      "LwM2MVersion": "1.1",
      "URN": "urn:oma:lwm2m:oma:0:1.2",
      "Resources": [
        {
          "Id": 0,
          "Name": "LwM2M Server URI",
          "Operations": "N",
          "Multiple": false,
          "Mandatory": true,
          "ResourceType": "string",
          "RangeOrEnums": "0-255 bytes",
          "ValueValidator": "NewRangeValidator(0 255)"
        },
        {
          "Id": 1,
          "Name": "Bootstrap Server",
          "Operations": "N",
          "Multiple": false,
          "Mandatory": true,
          "ResourceType": "bool"
        },
        {
          "Id": 2,
          "Name": "Security Mode",
          "Operations": "N",
          "Multiple": false,
          "Mandatory": true,
          "ResourceType": "int",
          "RangeOrEnums": "0-3",
          "ValueValidator": "NewRangeValidator(0 3)"
        },
        {
          "Id": 3,
          "Name": "Public Key or Identity",
          "Operations": "N",
          "Multiple": false,
          "Mandatory": true,
          "ResourceType": "opaque"
        },
        {
          "Id": 4,
          "Name": "Server Public Key or Identity",
          "Operations": "N",
          "Multiple": false,
          "Mandatory": true,
          "ResourceType": "opaque"
        },
        {
          "Id": 5,
          "Name": "Secret Key",
          "Operations": "N",
          "Multiple": false,
          "Mandatory": true,
          "ResourceType": "opaque"
        },
        {
          "Id": 6,
          "Name": "SMS Security Mode",
          "Operations": "N",
          "Multiple": false,
          "Mandatory": true,
          "ResourceType": "int",
          "RangeOrEnums": "0-255",
          "ValueValidator": "NewRangeValidator(0 255)"
        },
        {
          "Id": 7,
          "Name": "SMS Binding Key Parameters",
          "Operations": "N",
          "Multiple": false,
          "Mandatory": true,
          "ResourceType": "opaque",
          "RangeOrEnums": "6 bytes",
          "ValueValidator": "NewLengthValidator(6)"
        },
        {
          "Id": 8,
          "Name": "SMS Binding Secret Keys",
          "Operations": "N",
          "Multiple": false,
          "Mandatory": true,
          "ResourceType": "opaque",
          "RangeOrEnums": "32-48 bytes",
          "ValueValidator": "NewRangeValidator(32 48)"
        },
        {
          "Id": 9,
          "Name": "LwM2M Server SMS Number",
          "Operations": "N",
          "Multiple": false,
          "Mandatory": true,
          "ResourceType": "int"
        },
        {
          "Id": 10,
          "Name": "Short Server ID",
          "Operations": "N",
          "Multiple": false,
          "Mandatory": false,
          "ResourceType": "int",
          "RangeOrEnums": "1-65535",
          "ValueValidator": "NewRangeValidator(1 65535)"
        },
        {
          "Id": 11,
          "Name": "Client Hold Off Time",
          "Operations": "N",
          "Multiple": false,
          "Mandatory": true,
          "ResourceType": "int",
          "Units": "s"
        }
      ]
    }
`
