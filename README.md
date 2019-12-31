# Nano Work Cache

Nano PoW computation helper service for light clients.  Nano PoW generation proxy and cache, written in Go.

## Overview

In Nano, wallets need to create a small proof-of-work (PoW) with each new block.
This is to prevent spamming the network with many junk blocks.

Normally, PoW is precomputed (it depends only on data from the previous block), so it is available when a new block is created (except if new blocks are created at a very fast rate).
Full wallets (with a full node) have all it takes to compute work (using device, GPU, etc.), trigger when to initiate pre computation, state where to persist pre-computed PoW, etc.

It is more difficult for light wallets to do pre-computation.
Light wallets -- mobile wallets, web/based wallets -- typcially rely on backend, as they do not have access to a local full node.  Work generation is also cumbersome for them -- GPU access or software library may not be accessible.

![diagram](https://github.com/catenocrypt/nano-work-cache/blob/master/doc/nano_work_cache_diag.png)

The Nano Work Cache is a backend component to assist light wallets, by caching work for them.
Nano Work Cache sits between the client and a PoW source, proxies PoW requests, and caches results.
Light wallets need only to fire a simple message to trigger work computation in the backend, and later ask for work 
when it is needed.  At that time PoW will be retured fast from the cache.
The Nano Work Cache does not compute work itself, it just proxies work reqests to a node, and caches results.

## Prerequisites

* Golang
* Access to a Nano node with work generation and RPC API.

## Setting up

Get sources

```shell
git clone https://github.com/catenocrypt/nano-work-cache.git
cd nano-work-cache
```

Configure: create local config file

```shell
cd src/main
cp config.SAMPLE.toml config.toml
```

Edit the file, and set the remote node RPC URL.

Run the service:

```shell
cd src/main
go run main.go
```

For testing, from another shell:

```shell
cd src/sample
go run sample.go
```

## Usage by light wallets

### Fully transparent work_generate request with difficulty

Clients can do a work_generate request, exactly like they would call to a node.

Example:

```shell
curl -d '{"action":"work_generate","hash":"DDDA8C4CB5825FF4F5D00C5F923BC6F632414F67D17039228325392671C50FA2","difficulty":"ffffffc000000000"}' http://localhost:7176
```

The first response is proxied from the remote node:

```json
{
    "hash":"DDDA8C4CB5825FF4F5D00C5F923BC6F632414F67D17039228325392671C50FA2",
    "work":"bbe869e32c992096",
    "difficulty":"fffffff8ad570225",
    "multiplier":"8.739717559668671",
    "source":"fresh"
}
```

Note that it is extended with a `source` field.
Subsequent calls are returned instantly from the cache:

```json
{
    "hash":"DDDA8C4CB5825FF4F5D00C5F923BC6F632414F67D17039228325392671C50FA2",
    "work":"bbe869e32c992096",
    "difficulty":"fffffff8ad570225",
    "multiplier":"8.739717559668671",
    "source":"fromcache"
}
```

### Work_generate request without difficulty

This is similar to the full work_generate request, but without the difficulty.
The client does not have to worry about the current difficulty; NanoWorkCache will use the current network difficulty.

Example:

```shell
curl -d '{"action":"work_generate","hash":"DDDA8C4CB5825FF4F5D00C5F923BC6F632414F67D17039228325392671C50FA2"}' http://localhost:7176
```

### Simplified account-based precompute call

/TODO/

This is a simplified call for work precompute: the client only has to specify the relevant account, nothing else.

- No need to wait for the response, as this call returns immediately
- No need to keep track of current last block hash
- No need to keep track of current difficulty

## Not (yet) done

- Simplified account-based precompute call
- Periodically retrieve current difficulty from node
- Listen on new blocks from node; if a new block is created for a recently used account, start work computation right away, without being requested
