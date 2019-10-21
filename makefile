
deploy:
	GOOS=linux GOARCH=amd64 go build src/*go
	mv _author api
	scp ./api ubuntu@49.234.50.44:~
	rm api
