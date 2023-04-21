# Chap 04

## 4.4 Transaction

```sh
MULTI/EXEC

WATCH/UNWATCH/DISCARD
```

- Transactions in Redis are different from transactions that exist in more traditional relational databases. In a relational database, we can tell the database server `BEGIN`, at which point we can perform a variety of read and write operations that will be consistent with respect to each other, after which we can run either `COMMIT` to make our changes permanent or `ROLLBACK` to discard our changes.

- Within Redis, there’s a simple method for handling a sequence of reads and writes that will be consistent with each other. The transaction is begun by calling the special command `MULTI`, passing our series of commands, followed by `EXEC`. `The problem is that this simple transaction doesn’t actually do anything until EXEC is called`

- `WATCH` combine with `MULTI and EXEC`. When we’ve watched keys with `WATCH`, if at any time some other client replaces, updates, or deletes any keys that we’ve WATCHed before we have performed the EXEC operation, our operations against Redis will fail with an error message when we try to `EXEC` (at which point we can retry or abort the operation).

## 4.5 Non-transaction

In the case where we don’t need transactions, but where we still want to do a lot of work, we could still use `MULTI/EXEC` for their ability to send all of the commands at the same time to minimize round trips and latency.

```python
def update_token_pipeline(conn, token, user, item=None):
    timestamp = time.time()
    pipe = conn.pipeline(False)
    pipe.hset('login:', token, user)
    pipe.zadd('recent:', token, timestamp)
    if item:
        pipe.zadd('viewed:' + token, item, timestamp)
        pipe.zremrangebyrank('viewed:' + token, 0, -26)
        pipe.zincrby('viewed:', item, -1)
    pipe.execute()
```
