GoHPACK
=======

HPACK implementation in Golang based on [RFC 7541](http://tools.ietf.org/html/rfc7541 "RFC 7541")

## Decode
Pass test cases (hpack-test-case/haskell-http2...)

## Encode
Pass test cases

# Install & Test
```
$ go get github.com/ami-GS/GoHPACK
$ cd $GOPATH/src/github.com/ami-GS/
$ mv GoHPACK/test.go .
$ git clone https://github.com/http2jp/hpack-test-case
$ go test
```

#### License
The MIT License (MIT) Copyright (c) 2014 ami-GS