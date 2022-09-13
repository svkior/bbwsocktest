## Simple workaround connection to blockbook public nodes

Problem: 

- websocket error handshake

Because:

- 403 because of cloudflare wants an User-Agent in request header

```mermaid
sequenceDiagram
    autonumber
    participant C as Client
    participant CF as CloudFlare
    participant BB as BlockBook Service

    C->>CF: HTTP Request With Upgrade
    CF->>CF: Check for income connection
    alt Request With User-Agent
        CF->>BB: HTTP Request With Upgrade
        activate BB
        BB->>C: Handshake ....
    else Request without User-Agent
        CF->>C: 403
    end
```

Solution:

- Put User-Agent tag in the http header

How to run:

1. make build-base
2. make watch-wsclient