module github.com/GetWagz/hubspot-sdk

go 1.12

require (
	github.com/fatih/structs v1.1.0
	github.com/go-resty/resty v1.11.0
	github.com/sirupsen/logrus v1.3.0
	github.com/stretchr/testify v1.3.0
	golang.org/x/crypto v0.0.0-20190123085648-057139ce5d2b // indirect
	golang.org/x/net v0.0.0-20190119204137-ed066c81e75e // indirect
	golang.org/x/sys v0.0.0-20190123074212-c6b37f3e9285 // indirect
	gopkg.in/resty.v1 v1.12.0 // indirect
)

replace github.com/go-resty/resty => gopkg.in/resty.v1 v1.11.0
