package objects

var LocationDescriptor = `{
      "Id": 6,
      "Name": "Location",
      "Multiple": false,
      "Mandatory": false,
      "Version": "1.0",
      "LwM2MVersion": "1.0",
      "URN": "urn:oma:lwm2m:oma:6",
      "Resources": [
        {
          "Id": 0,
          "Name": "Latitude",
          "Operations": "R",
          "Multiple": false,
          "Mandatory": true,
          "ResourceType": "string",
          "Units": "Deg"
        },
        {
          "Id": 1,
          "Name": "Longitude",
          "Operations": "R",
          "Multiple": false,
          "Mandatory": true,
          "ResourceType": "string",
          "Units": "Deg"
        },
        {
          "Id": 2,
          "Name": "Altitude",
          "Operations": "R",
          "Multiple": false,
          "Mandatory": false,
          "ResourceType": "string",
          "Units": "m"
        },
        {
          "Id": 3,
          "Name": "Uncertainty",
          "Operations": "R",
          "Multiple": false,
          "Mandatory": false,
          "ResourceType": "string",
          "Units": "m"
        },
        {
          "Id": 4,
          "Name": "Velocity",
          "Operations": "R",
          "Multiple": false,
          "Mandatory": false,
          "ResourceType": "opaque",
          "Units": "Refers to 3GPP GAD specs"
        },
        {
          "Id": 5,
          "Name": "Timestamp",
          "Operations": "R",
          "Multiple": false,
          "Mandatory": true,
          "ResourceType": "time",
          "RangeOrEnums": "0-6",
          "ValueValidator": "NewRangeValidator(0 6)"
        }
      ]
    }`
