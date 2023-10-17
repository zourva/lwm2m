package preset

var deviceDescriptor = `{
      "Id": 3,
      "Name": "Device",
      "Multiple": false,
      "Mandatory": true,
      "Version": "1.2",
      "LwM2MVersion": "1.1",
      "URN": "urn:oma:lwm2m:oma:3:1.2",
      "Resources": [
        {
          "Id": 0,
          "Name": "Manufacturer",
          "Operations": "R",
          "Multiple": false,
          "Mandatory": false,
          "ResourceType": "string"
        },
        {
          "Id": 1,
          "Name": "Model Number",
          "Operations": "R",
          "Multiple": false,
          "Mandatory": false,
          "ResourceType": "string"
        },
        {
          "Id": 2,
          "Name": "Serial Number",
          "Operations": "R",
          "Multiple": false,
          "Mandatory": false,
          "ResourceType": "string"
        },
        {
          "Id": 3,
          "Name": "Firmware Version",
          "Operations": "R",
          "Multiple": false,
          "Mandatory": false,
          "ResourceType": "string"
        },
        {
          "Id": 4,
          "Name": "Reboot",
          "Operations": "E",
          "Multiple": false,
          "Mandatory": true,
          "ResourceType": "string"
        },
        {
          "Id": 5,
          "Name": "Factory Reset",
          "Operations": "E",
          "Multiple": false,
          "Mandatory": false,
          "ResourceType": "string"
        },
        {
          "Id": 6,
          "Name": "Available Power Sources",
          "Operations": "R",
          "Multiple": true,
          "Mandatory": false,
          "ResourceType": "int",
          "RangeOrEnums": "0-7",
          "ValueValidator": "NewRangeValidator(0 7)"
        },
        {
          "Id": 7,
          "Name": "Power Source Voltage",
          "Operations": "R",
          "Multiple": true,
          "Mandatory": false,
          "ResourceType": "int",
          "Units": "mV"
        },
        {
          "Id": 8,
          "Name": "Power Source Current",
          "Operations": "R",
          "Multiple": true,
          "Mandatory": false,
          "ResourceType": "int",
          "Units": "mA"
        },
        {
          "Id": 9,
          "Name": "Battery Level",
          "Operations": "R",
          "Multiple": false,
          "Mandatory": false,
          "ResourceType": "int",
          "RangeOrEnums": "0-100",
          "Units": "%",
          "ValueValidator": "NewRangeValidator(0 100)"
        },
        {
          "Id": 10,
          "Name": "Memory Free",
          "Operations": "R",
          "Multiple": false,
          "Mandatory": false,
          "ResourceType": "int",
          "Units": "KB"
        },
        {
          "Id": 11,
          "Name": "Error Code",
          "Operations": "R",
          "Multiple": true,
          "Mandatory": true,
          "ResourceType": "int"
        },
        {
          "Id": 12,
          "Name": "Reset Error Code",
          "Operations": "E",
          "Multiple": false,
          "Mandatory": false,
          "ResourceType": "string"
        },
        {
          "Id": 13,
          "Name": "Current Time",
          "Operations": "RW",
          "Multiple": false,
          "Mandatory": false,
          "ResourceType": "time"
        },
        {
          "Id": 14,
          "Name": "UTC Offset",
          "Operations": "RW",
          "Multiple": false,
          "Mandatory": false,
          "ResourceType": "string"
        },
        {
          "Id": 15,
          "Name": "Timezone",
          "Operations": "RW",
          "Multiple": false,
          "Mandatory": false,
          "ResourceType": "string"
        },
        {
          "Id": 16,
          "Name": "Supported Binding and Modes",
          "Operations": "R",
          "Multiple": false,
          "Mandatory": true,
          "ResourceType": "string"
        },
        {
          "Id": 17,
          "Name": "Device Type",
          "Operations": "R",
          "Multiple": false,
          "Mandatory": false,
          "ResourceType": "string"
        },
        {
          "Id": 18,
          "Name": "Hardware Version",
          "Operations": "R",
          "Multiple": false,
          "Mandatory": false,
          "ResourceType": "string"
        },
        {
          "Id": 19,
          "Name": "Software Version",
          "Operations": "R",
          "Multiple": false,
          "Mandatory": false,
          "ResourceType": "string"
        },
        {
          "Id": 20,
          "Name": "Battery Status",
          "Operations": "R",
          "Multiple": false,
          "Mandatory": false,
          "ResourceType": "int",
          "RangeOrEnums": "0-6",
          "ValueValidator": "NewRangeValidator(0 6)"
        },
        {
          "Id": 21,
          "Name": "Memory Total",
          "Operations": "R",
          "Multiple": false,
          "Mandatory": false,
          "ResourceType": "int"
        },
        {
          "Id": 22,
          "Name": "ExtDevInfo",
          "Operations": "R",
          "Multiple": true,
          "Mandatory": false,
          "ResourceType": "objectlink"
        }
      ]
    }
`
