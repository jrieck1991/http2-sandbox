# http2

sandbox for http/2 experiments

## TLS

generate server.key and server.crt:
`openssl req -newkey rsa:2048 -nodes -keyout server.key -x509 -days 365 -out server.crt`

## curls

```bash
curl https://localhost:8000
```
