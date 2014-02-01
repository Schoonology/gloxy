# Gloxy

A trivial reverse proxy implementation in Go that logs all requests made and
responses received.

## Why?

I wanted a really simple tool for two purposes:

 1. Debugging in-development REST API servers. I could pepper the server
 application with logging everywhere, but that takes a lot of time, and might
 not be trivial (depending on the environment). Gloxy gets the same result, and
 can be run on the server (to log everything) or on the client (to log that
 individual's requests).
 1. Reverse engineering existing REST APIs. This was Gloxy's original purpose,
 as I wanted an easy-to-read log of all the transactions a CouchDB server
 requested during replication. Since Gloxy logs only what it proxies, I could
 have the Admin console open and pointing directly to Couch (so the HTML
 requests weren't logged), with the replication done through Gloxy so that I
 could see every juicy detail.
