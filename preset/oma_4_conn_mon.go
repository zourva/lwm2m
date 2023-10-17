package preset

var connMonitorDescriptor = `{
      "Id": 4,
      "Name": "Connectivity Monitoring",
      "Multiple": false,
      "Mandatory": false,
      "Version": "1.3",
      "LwM2MVersion": "1.1",
      "URN": "urn:oma:lwm2m:oma:4:1.3",
      "Resources": [
        {
          "Id": 0,
          "Name": "Network Bearer",
          "Operations": "R",
          "Multiple": false,
          "Mandatory": true,
          "ResourceType": "int"
        },
        {
          "Id": 1,
          "Name": "Available Network Bearer",
          "Operations": "R",
          "Multiple": true,
          "Mandatory": true,
          "ResourceType": "int"
        },
        {
          "Id": 2,
          "Name": "Radio Signal Strength",
          "Operations": "R",
          "Multiple": false,
          "Mandatory": true,
          "ResourceType": "int",
          "Units": "dBm"
        },
        {
          "Id": 3,
          "Name": "Link Quality",
          "Operations": "R",
          "Multiple": false,
          "Mandatory": false,
          "ResourceType": "int"
        },
        {
          "Id": 4,
          "Name": "IP Addresses",
          "Operations": "R",
          "Multiple": true,
          "Mandatory": true,
          "ResourceType": "string"
        },
        {
          "Id": 5,
          "Name": "Router IP Addresses",
          "Operations": "R",
          "Multiple": true,
          "Mandatory": false,
          "ResourceType": "string"
        },
        {
          "Id": 6,
          "Name": "Link Utilization",
          "Operations": "R",
          "Multiple": false,
          "Mandatory": false,
          "ResourceType": "int",
          "RangeOrEnums": "0-100",
          "Units": "%",
          "ValueValidator": "NewRangeValidator(0 100)"
        },
        {
          "Id": 7,
          "Name": "APN",
          "Operations": "R",
          "Multiple": true,
          "Mandatory": false,
          "ResourceType": "string"
        },
        {
          "Id": 8,
          "Name": "Cell ID",
          "Operations": "R",
          "Multiple": false,
          "Mandatory": false,
          "ResourceType": "int"
        },
        {
          "Id": 9,
          "Name": "SMNC",
          "Operations": "R",
          "Multiple": false,
          "Mandatory": false,
          "ResourceType": "int",
          "Units": "%"
        },
        {
          "Id": 10,
          "Name": "SMCC",
          "Operations": "R",
          "Multiple": false,
          "Mandatory": false,
          "ResourceType": "int"
        }
      ]
    }
`
