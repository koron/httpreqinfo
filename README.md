# httpreqinfo - HTTP server to dump request information

httpreqinfo is a HTTP server, returns a JSON with request information.

```console
$ go run

# or

$ go build
$ ./httpreqinfo
```

It starts to listen `:8000`.

Example of response:

```console
$ curl http://127.0.0.1:8000/
{
  "method": "GET",
  "url": "/",
  "proto": "HTTP/1.1",
  "header": {
    "Accept": "*/*",
    "User-Agent": "curl/7.71.1"
  },
  "content-length": 0,
  "host": "127.0.0.1:8000"
}
```
