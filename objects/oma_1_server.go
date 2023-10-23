package objects

var ServerDescriptor = `{
      "Id": 1,
      "Name": "LwM2M Server",
      "Multiple": true,
      "Mandatory": true,
      "Version": "1.2",
      "LwM2MVersion": "1.2",
      "URN": "urn:oma:lwm2m:oma:1:1.2",
      "Resources": [
        {
          "Id": 0,
          "Name": "Short Server ID",
          "Operations": "R",
          "Multiple": false,
          "Mandatory": true,
          "ResourceType": "int",
          "RangeOrEnums": "1-65535",
          "ValueValidator": "NewRangeValidator(1 65535)"
        },
        {
          "Id": 1,
          "Name": "Lifetime",
          "Operations": "RW",
          "Multiple": false,
          "Mandatory": true,
          "ResourceType": "int",
          "Units": "s"
        },
        {
          "Id": 2,
          "Name": "Default Minimum Period",
          "Operations": "RW",
          "Multiple": false,
          "Mandatory": false,
          "ResourceType": "int",
          "Units": "s"
        },
        {
          "Id": 3,
          "Name": "Default Maximum Period",
          "Operations": "RW",
          "Multiple": false,
          "Mandatory": false,
          "ResourceType": "int",
          "Units": "s"
        },
        {
          "Id": 4,
          "Name": "Disable",
          "Operations": "E",
          "Multiple": false,
          "Mandatory": false,
          "ResourceType": "string"
        },
        {
          "Id": 5,
          "Name": "Disable Timeout",
          "Operations": "RW",
          "Multiple": false,
          "Mandatory": false,
          "ResourceType": "int",
          "Units": "s"
        },
        {
          "Id": 6,
          "Name": "Notification Storing When Disabled or Offline",
          "Operations": "RW",
          "Multiple": false,
          "Mandatory": true,
          "ResourceType": "bool"
        },
        {
          "Id": 7,
          "Name": "Binding",
          "Operations": "RW",
          "Multiple": false,
          "Mandatory": true,
          "ResourceType": "string",
          "RangeOrEnums": "The possible values of Resource are listed in 5.2.1.1"
        },
        {
          "Id": 8,
          "Name": "Registration Update Trigger",
          "Operations": "E",
          "Multiple": false,
          "Mandatory": true,
          "ResourceType": "string"
        },
        {
          "Id": 9,
          "Name": "BootstrapRequest Trigger",
          "Operations": "E",
          "Multiple": false,
          "Mandatory": false,
          "ResourceType": "string"
        },
        {
          "Id": 10,
          "Name": "APN Link",
          "Operations": "RW",
          "Multiple": false,
          "Mandatory": false,
          "ResourceType": "objectlink"
        },
        {
          "Id": 11,
          "Name": "TLS-DTLS Alert Code",
          "Operations": "R",
          "Multiple": false,
          "Mandatory": false,
          "ResourceType": "int",
          "RangeOrEnums": "0-255",
          "ValueValidator": "NewRangeValidator(0 255)"
        },
        {
          "Id": 12,
          "Name": "Last Bootstrapped",
          "Operations": "R",
          "Multiple": false,
          "Mandatory": false,
          "ResourceType": "time"
        },
        {
          "Id": 13,
          "Name": "Registration Priority Order",
          "Operations": "R",
          "Multiple": false,
          "Mandatory": false,
          "ResourceType": "int"
        },
        {
          "Id": 14,
          "Name": "Initial Registration Delay Timer",
          "Operations": "RW",
          "Multiple": false,
          "Mandatory": false,
          "ResourceType": "int",
          "Units": "s"
        },
        {
          "Id": 15,
          "Name": "Registration Failure Block",
          "Operations": "R",
          "Multiple": false,
          "Mandatory": false,
          "ResourceType": "bool"
        },
        {
          "Id": 16,
          "Name": "Bootstrap on Registration Failure",
          "Operations": "R",
          "Multiple": false,
          "Mandatory": false,
          "ResourceType": "bool"
        },
        {
          "Id": 17,
          "Name": "Communication Retry Count",
          "Operations": "RW",
          "Multiple": false,
          "Mandatory": false,
          "ResourceType": "int"
        },
        {
          "Id": 18,
          "Name": "Communication Retry Timer",
          "Operations": "RW",
          "Multiple": false,
          "Mandatory": false,
          "ResourceType": "int",
          "Units": "s"
        },
        {
          "Id": 19,
          "Name": "Communication Sequence Delay Timer",
          "Operations": "RW",
          "Multiple": false,
          "Mandatory": false,
          "ResourceType": "int",
          "Units": "s"
        },
        {
          "Id": 20,
          "Name": "Communication Sequence Retry Count",
          "Operations": "RW",
          "Multiple": false,
          "Mandatory": false,
          "ResourceType": "int"
        },
        {
          "Id": 21,
          "Name": "Trigger",
          "Operations": "RW",
          "Multiple": false,
          "Mandatory": false,
          "ResourceType": "bool"
        },
        {
          "Id": 22,
          "Name": "Preferred Transport",
          "Operations": "RW",
          "Multiple": false,
          "Mandatory": false,
          "ResourceType": "string"
        },
        {
          "Id": 23,
          "Name": "Mute Send",
          "Operations": "RW",
          "Multiple": false,
          "Mandatory": false,
          "ResourceType": "bool"
        },
        {
          "Id": 24,
          "Name": "Alternate APN Links",
          "Operations": "RW",
          "Multiple": true,
          "Mandatory": false,
          "ResourceType": "objectlink"
        },
        {
          "Id": 25,
          "Name": "Supported Server Versions",
          "Operations": "RW",
          "Multiple": true,
          "Mandatory": false,
          "ResourceType": "string"
        },
        {
          "Id": 26,
          "Name": "Default Notification Mode",
          "Operations": "RW",
          "Multiple": false,
          "Mandatory": false,
          "ResourceType": "int",
          "RangeOrEnums": "0-1",
          "ValueValidator": "NewRangeValidator(0 1)"
        },
        {
          "Id": 27,
          "Name": "Profile ID Hash Algorithm",
          "Operations": "RW",
          "Multiple": false,
          "Mandatory": false,
          "ResourceType": "int",
          "RangeOrEnums": "0-255",
          "ValueValidator": "NewRangeValidator(0 255)"
        }
      ]
    }
`
