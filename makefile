deploy:
	GOOS=linux GOARCH=amd64 go build src/*go
	mv article api_server
	scp ./api_server ubuntu@cdr.today:~
	rm api_server
