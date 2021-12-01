package main

const ProtoExample string = `syntax = "proto3";
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
    }`
