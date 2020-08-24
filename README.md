# TWHelp Golang SDK

Basic Golang SDK for making HTTP requests to [TWHelp API](https://api.tribalwarshelp.com).

## Getting started

This SDK requires Go version with Modules support. Make sure to initialize a Go module:

```
go mod init github.com/my/repo
```

and then install SDK:

```
go get -u github.com/tribalwarshelp/golang-sdk
```

### How to initialize SDK in your Go code?

```go
api := sdk.New("url")
```

Example url: https://api.tribalwarshelp.com/
