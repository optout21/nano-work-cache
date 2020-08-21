# Nano Work Cache

Nano PoW computation helper service, relieve light clients from the burden of work precomputation.
Nano PoW generation proxy and cache, written in Go.

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

In a real-world usage, cache hit ratio was **92%**.  In a larger wallet setup, over a period of several days, 3412 work requests were served from the cache, out of 3703.  During the same period, an additional 2627 work proofs have been precomputed.

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
curl -d '{"action": "work_generate","hash": "718CC2121C3E641059BC1C2CFC45666C99E8AE922F7A807B7D07B62C995D79E2","difficulty": "ffffffd21c3933f3"}' localhost:7376
```

```shell
cd src/sample
go run sample.go
```

## Usage by light wallets

Light wallets, working together with NanoWorkCache, can access precomputed work without much effort.
There are several styles for integration:

1. Fully automatic: the wallet has to do nothing extra:
- at early moment, requests account balance -- this triggers work computation 
- at send moment, ask for work by a regular work requests (the result will be delivered from cache quickly).

2. They also have the option to control even more the work generation.  Work generation can be triggered:
- Fully transparent work_generate request (with difficulty).  Wallet has to disregard the result, just fire-and-forget.
  _Variant:_ without difficulty, difficulty is filled by NanoWorkCache.
- Special async `work_pregenerate_by_hash` call invoked at an early time; only the hash is needed.  
  _Variant:_  work_pregenerate_by_account, when not event the last block hash is needed, only the account.

See also the details of the [Integration Options (API.md)](API.md).

## Not (yet) done

- Periodically retrieve current difficulty from node
- Store cache in Redis
- Support nano work peers
- Cache cleanup: remove old entries for accounts, for which newer one exists
- Listen on new blocks from node; if a new block is created for a recently used account, trigger work computation for new block right away
  (start pre-work even without any user activity).
