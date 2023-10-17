package preset

var accessControlDescriptor = `{
      "Id": 2,
      "Name": "LwM2M Access Control",
      "Multiple": true,
      "Mandatory": false,
      "Version": "1.1",
      "LwM2MVersion": "1.0",
      "URN": "urn:oma:lwm2m:oma:2:1.1",
      "Resources": [
        {
          "Id": 0,
          "Name": "Object ID",
          "Operations": "R",
          "Multiple": false,
          "Mandatory": true,
          "ResourceType": "int",
          "RangeOrEnums": "1-65534",
          "ValueValidator": "NewRangeValidator(1 65534)"
        },
        {
          "Id": 1,
          "Name": "Object Instance ID",
          "Operations": "R",
          "Multiple": false,
          "Mandatory": true,
          "ResourceType": "int",
          "RangeOrEnums": "0-65535",
          "ValueValidator": "NewRangeValidator(0 65535)"
        },
        {
          "Id": 2,
          "Name": "ACL",
          "Operations": "RW",
          "Multiple": true,
          "Mandatory": false,
          "ResourceType": "int",
          "RangeOrEnums": "16-bit",
          "ValueValidator": "NewRangeValidator(-32768 32767)"
        },
        {
          "Id": 3,
          "Name": "Access Control Owner",
          "Operations": "RW",
          "Multiple": false,
          "Mandatory": true,
          "ResourceType": "int",
          "RangeOrEnums": "0-65535",
          "ValueValidator": "NewRangeValidator(0 65535)"
        }
      ]
    }
`
