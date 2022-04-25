# Shorthills TreeService Signed URL Lambda Function

Lambda function to create signed URLs


## Issues

Ran into an issue when compiling:
`//go:linkname must refer to declared function or variable`

Found resolution here: https://stackoverflow.com/questions/71507321/go-1-18-build-error-on-mac-unix-syscall-darwin-1-13-go253-golinkname-mus

Essentially it involves installing this package: `go get -u golang.org/x/sys`
