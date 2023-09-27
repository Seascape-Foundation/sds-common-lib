# Common Lib

This module consists of the data types used in all go modules.

## Data Types

Stored in the `data_type` package.
This package consists of various data types.

### Queue
First in, First out queue. Available in `data_type.Queue`.
The queue has a cap.

### Database
Stored in the `data_type/database` package.
It Consists of the SQL database interfaces that all databases must implement.

### KeyValue
Stored in the `data_type/key_value.KeyValue`.
The `KeyValue` is the wrapper around `map[string]interface{}`.
The structure adds methods to extract the data with the validation.

The structure also can not contain the null parameters.

All the keys are string.

The values can be:
* `NestedValue` &ndash; nested `KeyValue`
* `NestedListValue` &ndash; list of `KeyValue`.
* `Uint64` &ndash; any natural numbers and zero are converted into go's `uint64` type. **KeyValue negative numbers are represented as Float64**.
* `Float64` &ndash; the number is represented as a float of 64 bits.
* `String` &ndash; a string.
* `Strings` &ndash; a slice of string.
* `BigNumber` &ndash; the number of `big.Int` format.
* `Bool` &ndash; a boolean parameter.

### KeyValueList
Stored in the `data_type/key_value.List`.
The `List` is the `KeyValue` with two conditions:
Keys are interface, not string. 
The key types must be identical.
The value types must be identical.

**The `List` has a cap**.

If the first value that you set is the number key, and struct A value.
Then all keys must be number, and all values must be struct A.

---

## Message
The messages are the data that SDS uses for intercommunication.
There are two types of the messages: Requests and Replies.

Internally, the SDS framework uses the interface: `RequestInterface` and `ReplyInterface`.
Any message type must implement the interfaces.

### Operations
All messages are grouped into the `message.Operations`.
It's a structure with the function references.
The SDS Framework accepts the `message.Operations` for a message manipulation.

* `NewReq([]messages) (RequestInterface, error)`
* `NewReply([]messages) (ReplyInterface, error)`
* `EmptyReq() RepquestInterface`
* `EmptyReply() ReplyInterface`

### Built in message types
The SDS comes with two types of messages as well as their operations.

The default message is the JSON compatible Structure.
The raw message is the represents the strings.

#### Default Message
The default message types are JSON compatible.
The operations are defined as `DefaultMessage` in the `message/request.go` file.

The request format is defined in `message.Request`.
The reply format is defined in `message.Reply`.

#### Raw Message
The messages are the wrappers around zeromq message envelopes.
The SDS framework as a framework to write distributed systems uses zeromq internally.

The `Raw` messages are the wrappers around zeromq's messages.
The wrapped messages are fully compatible with the SDS framework.
Because raw messages implement `RequestInterface` and `ReplyInterface`.

Both request and replies are defined in `message/raw.go` file.
The request format is defined in `message.RawRequest`.
The reply format is defined in `message.RawReply`.

The operations are defined as `message.RawMessage()`.

> **Todo** 
> Set an actual link to the documentation for usage of the messages in clients and handlers.

Refer to [client-lib](https://github.com/ahmetson/client-lib) and
[handler-lib](https://github.com/ahmetson/handler-lib) for changing default message format.

### Trace
To make debugging easy in the distributed systems, the messages come with the traces.
The message that was sent by a user to the service A, 
which then forwarded to service B must include the service information as a stack.

Therefore, there is an important rule:

**The next request must from service A to service B must be defined through `request.Next`**.
The package that handled the message must include its stack with the information about the service, handler and a time.