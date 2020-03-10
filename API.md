## Integration options for light wallets

'Full' wallets work this way:

1. At an 'early' moment -- at startup, wallet open, etc. -- start work computation for current block, and store result
2. At the 'send' moment, when work is needed, they take stored result (and start work computation once the new block has been created).

Light wallets, working together with NanoWorkCache, have a much easier job:

1. Fully automatic: the wallet has to do nothing extra:
- at early moment, requests account balance -- this triggers work computation 
- at send moment, ask for work by a regular work requests (the result will be delivered from cache quickly).

2. They also have the option to control work generation more precisely.  Work generation can be triggered by:
- Fully transparent work_generate request (with difficulty).  Wallet has to disregard the result, just fire-and-forget.
  _Variant:_ without difficulty, difficulty is filled by NanoWorkCache.
- Special async `work_pregenerate_by_hash` call invoked at an early time; only the hash is needed.  
  _Variant:_  work_pregenerate_by_account, when not event the last block hash is needed, only the account.

### Transparent account_balance request

The client simply makes an `account_balance` call when the wallet is started, to display the current balance.
NanoWorkCache forwards the call to the node, and returns the balance.
However, work precomputation is triggered, so when later a work is requested for the last block of this account, likely it will be delivered quickly from cache.

The `accounts_balances` call is similarly supported.

```shell
curl -d '{"action":"account_balance","account":"nano_3rpb7ddcd6kux978gkwxh1i1s6cyn7pw3mzdb9aq7jbtsdfzceqdt3jureju"}' http://localhost:7176
```

### Fully transparent work_generate request with difficulty

Clients can do a `work_generate` request, exactly like they would call to a node.

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
    "source":"cache"
}
```

### Work_generate request without difficulty

This is similar to the full `work_generate` request, but without the difficulty.
The client does not have to worry about the current difficulty; NanoWorkCache will use the current network difficulty.

Example:

```shell
curl -d '{"action":"work_generate","hash":"DDDA8C4CB5825FF4F5D00C5F923BC6F632414F67D17039228325392671C50FA2"}' http://localhost:7176
```

### Simplified hash-based precompute call

This is a simplified call for work precompute: the client only has to specify the relevant block hash.  The action is `work_pregenerate_by_hash`.

- No need to wait for the response, as this call returns immediately
- No need to keep track of current difficulty

Example:

```shell
curl -d '{"action":"work_pregenerate_by_hash","hash":"DDDA8C4CB5825FF4F5D00C5F923BC6F632414F67D17039228325392671C50FA2"}' http://localhost:7176
```

### Simplified account-based precompute call

This is a simplified call for work precompute: the client only has to specify the relevant account, nothing else.
The action is `work_pregenerate_by_account`.  NanoWorkCache retrieves the current latest block of the account, and triggers work computation for it (unless it is in the cache already).

- No need to wait for the response, as this call returns immediately
- No need to keep track of current last block hash
- No need to keep track of current difficulty

Example:

```shell
curl -d '{"action":"work_pregenerate_by_account","account":"nano_3rpb7ddcd6kux978gkwxh1i1s6cyn7pw3mzdb9aq7jbtsdfzceqdt3jureju"}' http://localhost:7176
```
