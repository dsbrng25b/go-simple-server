# go-simple-server
go-simple-server is a HTTP server which prints all requests it receives to stdout.
I often use this code as a staring point to get a HTTP server to debug HTTP clients.

* Binds to `0.0.0.0:8080` (can be changed with `-bind`)
* Print HTTP method, path and headers to stdout (unless `-header=false`).
* Print request body to stdout (only printable characters).
  * If `Content-Type` is `application/json` or `-json` is set the body is treated as json and pretty-printed.
* If `-tls` is set it serves TLS with `tls.key` and `tls.crt` from the current working directory.

## Usage
* Start server
```
go-simple-server
```

* Send requests
```
curl http://localhost:8080
```
