pro:
	GOOS=linux GOARCH=amd64 go build src/*go
	mv article api_server
