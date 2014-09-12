install:
	go get code.google.com/p/goauth2/oauth github.com/google/go-github github.com/google/go-querystring

test:
	go test bitbucket.org/alexisjanvier/deployedpr/deptools
