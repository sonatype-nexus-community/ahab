Contains a few test dockerfiles that can be built using the version of ahab of the branch you are currently on. 
Simply helps to validate that ahab behaves like it should.

##How to use these docker files : 
1. Simply running docker_test.go will build ahab, and build the docker files and validate the output is as expected
2. From the docker directory you can run the following

### Yum
```
GOOS=linux GOARCH=amd64 go build -o ahab ../main.go
docker build -f yum/Dockerfile .
```

### Apt
```
GOOS=linux GOARCH=amd64 go build -o ahab ../main.go
docker build -f apt/Dockerfile .
```