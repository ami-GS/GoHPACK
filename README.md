GoHPACK
=======

HPACK implementation in Golang based on [draft 09](http://tools.ietf.org/html/draft-ietf-httpbis-header-compression-09 "draft 09")

## Decode
Pass test cases (hpack-test-case/haskell-http2...)

## Encode
Pass test cases

# Test
```
$ go run test <param (-e -d -a)>
```
* -e: encode test (this must emit error)
* -d: decode test
* -a: encode to decode test