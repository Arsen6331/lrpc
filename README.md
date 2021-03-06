# lrpc

[![Go Reference](https://pkg.go.dev/badge/go.arsenm.dev/lrpc.svg)](https://pkg.go.dev/go.arsenm.dev/lrpc)
[![Go Report Card](https://goreportcard.com/badge/go.arsenm.dev/lrpc)](https://goreportcard.com/report/go.arsenm.dev/lrpc)

A lightweight RPC framework that aims to be as easy to use as possible, while also being as lightweight as possible. Most current RPC frameworks are bloated to the point of adding 7MB to my binary, like RPCX. That is what prompted me to create this.

---

### Channels

This RPC framework supports creating channels to transfer data from server to client. My use-case for this is to implement watch functions and transfer progress in [ITD](https://gitea.arsenm.dev/Arsen6331/itd), but it can be useful for many things.

---

### Codec

When creating a server or client, a `CodecFunc` can be provided. An `io.ReadWriter` is passed into the `CodecFunc` and it returns a `Codec`, which is an interface that contains encode and decode functions with the same signature as `json.Decoder.Decode()` and `json.Encoder.Encode()`.

This allows any codec to be used for the transfer of the data, making it easy to create clients in different languages.

---

### Web Client

Inside `client/web`, there is a web client for lrpc using WebSockets. It is written in ruby (I don't like JS) and translated to human-readable JS using Ruby2JS. With the `bundler` gem installed, cd into `client/web` and run `make`. This will create a new file called `lrpc.js`, which can be used within a browser. It uses `crypto.randomUUID()`, so it must be used on an https site, not http.