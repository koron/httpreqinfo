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

## Commands for test

```
# 生JSON(中)を投げて、生版のハッシュが返ってくることを確認する
curl http://${TARGET_HOST}/test0/ -H 'Content-Type: application/json' --data-binary @testdata/bm20200819b.json

# 圧縮JSON(中)を投げて、圧縮版のハッシュが返ってくることを確認する
curl http://${TARGET_HOST}/test0/ -H 'Content-Type: application/json' --data-binary @testdata/bm20200819b.json.gz

# 生JSON(中)を投げて、生版のハッシュが返ってくることを確認する
curl http://${TARGET_HOST}/test1/ -H 'Content-Type: application/json' --data-binary @testdata/bm20200819b.json

# リクエスト展開プロキシを挟んで
# 圧縮JSON(中)を投げて、圧縮版のハッシュが返ってくることを確認する
curl http://${TARGET_HOST}/test1/ -H 'Content-Type: application/json' --data-binary @testdata/bm20200819b.json.gz

# リクエスト展開プロキシを挟んで
# 圧縮JSON(中)を投げて、生版のハッシュが返ってくることを確認する
curl http://${TARGET_HOST}/test1/ -H 'Content-Type: application/json' -H 'Content-Encoding: gzip' --data-binary @testdata/bm20200819b.json.gz
```

参考ハッシュ値:

```console
$ md5sum testdata/bm*
5d1f6efb229631c95d9e6d79ebff5319 *testdata/bm20200819a.json
f08318ebf2b28f0996b74a7a74e55373 *testdata/bm20200819a.json.gz
fc594bf63ae918046807dba16ec5aef1 *testdata/bm20200819b.json
fd7b83eaedd14cca6f3fe099211ab1a6 *testdata/bm20200819b.json.gz
```
