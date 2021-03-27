# GoBase64

Base64 encoder/decoder based on https://tools.ietf.org/html/rfc4648
Reads a file, and prints the result to stdout.


### Usage:
Encode:
```
$ gobase64 <file>
```
Decode:
```
$ gobase64 -d <file>
```

If no file is provided, it will read from stdin.