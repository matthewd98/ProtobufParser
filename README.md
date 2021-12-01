# ProtobufParser
Hackathon project at work: parse protobuf (.proto) files to JSON. The JSON can then be used to build documentation.

Example JSON:

```json
{
  "repo": {
    "name": "",
    "url": ""
  },
  "schema": [
    {
      "filePath": "",
      "fileName": "Example file",
      "url": "",
      "packageName": "custom",
      "enums": [
        {
          "name": "Color",
          "comment": "Color enumeration",
          "values": [
            {
              "name": "White",
              "value": 0
            },
            {
              "name": "Black",
              "value": 1
            }
          ]
        },
        {
          "name": "Foo.Letter",
          "comment": "Nested enumeration",
          "values": [
            {
              "name": "A",
              "value": 0
            },
            {
              "name": "B",
              "value": 1
            },
            {
              "name": "C",
              "value": 2
            }
          ]
        }
      ],
      "messages": [
        {
          "name": "Foo",
          "comment": "Foo message definition",
          "extensions": {
            "minTag": 0,
            "maxTag": 0
          },
          "fields": [
            {
              "name": "field1",
              "comment": "Field 1 is optional",
              "type": "string",
              "tag": 1,
              "isRequired": false,
              "isRepeated": false,
              "isExtension": false,
              "annotation": ""
            },
            {
              "name": "field2",
              "comment": "Field 2 is required",
              "type": "string",
              "tag": 2,
              "isRequired": true,
              "isRepeated": false,
              "isExtension": false,
              "annotation": ""
            },
            {
              "name": "field3",
              "comment": "Field 3 is implicitly optional",
              "type": "int32",
              "tag": 3,
              "isRequired": false,
              "isRepeated": false,
              "isExtension": false,
              "annotation": ""
            },
            {
              "name": "letter",
              "comment": "",
              "type": "Letter",
              "tag": 3,
              "isRequired": false,
              "isRepeated": false,
              "isExtension": false,
              "annotation": ""
            },
            {
              "name": "color",
              "comment": "",
              "type": "Color",
              "tag": 4,
              "isRequired": false,
              "isRepeated": false,
              "isExtension": false,
              "annotation": ""
            },
            {
              "name": "bar",
              "comment": "",
              "type": "Bar",
              "tag": 7,
              "isRequired": false,
              "isRepeated": true,
              "isExtension": false,
              "annotation": ""
            }
          ],
          "oneofs": [
            {
              "name": "Options",
              "comment": "One of the options can be set",
              "fields": [
                {
                  "name": "option1",
                  "comment": "",
                  "type": "bool",
                  "tag": 5,
                  "isRepeated": false,
                  "annotation": ""
                },
                {
                  "name": "option2",
                  "comment": "",
                  "type": "bool",
                  "tag": 6,
                  "isRepeated": false,
                  "annotation": ""
                }
              ]
            }
          ],
          "nestedMessages": [
            {
              "name": "Foo.Bar",
              "comment": "Nested message",
              "extensions": {
                "minTag": 0,
                "maxTag": 0
              },
              "fields": [
                {
                  "name": "nestedField",
                  "comment": "",
                  "type": "field",
                  "tag": 1,
                  "isRequired": false,
                  "isRepeated": false,
                  "isExtension": false,
                  "annotation": ""
                }
              ],
              "oneofs": null,
              "nestedMessages": null
            }
          ]
        }
      ],
      "services": [
        {
          "name": "FooBar",
          "comment": "",
          "rpcs": [
            {
              "name": "RPCFooBar",
              "comment": "",
              "rpcInput": {
                "type": "Foo",
                "isStream": false
              },
              "rpcOutput": {
                "type": "Foo",
                "isStream": false
              }
            }
          ]
        }
      ]
    }
  ]
}
```

Based on this example proto file:

``` proto
    syntax = "proto3";
    package custom;

    // Color enumeration
    enum Color {
        White = 0;
        Black = 1;
    }

    /* Foo message definition */
    message Foo {
        // Field 1 is optional
        optional string field1 = 1;
        // Field 2 is required
        required string field2 = 2;
        // Field 3 is implicitly optional
        int32 field3 = 3;

        // Nested enumeration
        enum Letter {
            A = 0;
            B = 1;
            C = 2;
        }
        optional Letter letter = 3;

        Color color = 4;

        // One of the options can be set
        oneof Options {
            bool option1 = 5;
            bool option2 = 6;
        }

        // Nested message
        message Bar {
            string field nestedField = 1;
        }
        repeated Bar bar = 7;

        extensions 100 to 200;
    }

    service FooBar {
        /* FooBar remote procedure call */
        rpc RPCFooBar (Foo) returns (Foo) {}
    }
```
