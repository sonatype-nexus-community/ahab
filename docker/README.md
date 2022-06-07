<!--

    Copyright (c) 2019-present Sonatype, Inc.

    Licensed under the Apache License, Version 2.0 (the "License");
    you may not use this file except in compliance with the License.
    You may obtain a copy of the License at

        http://www.apache.org/licenses/LICENSE-2.0

    Unless required by applicable law or agreed to in writing, software
    distributed under the License is distributed on an "AS IS" BASIS,
    WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
    See the License for the specific language governing permissions and
    limitations under the License.

-->
Contains a few test dockerfiles that can be built using the version of ahab of the branch you are currently on. 
Simply helps to validate that ahab behaves like it should.

##How to use these docker files : 
1. Simply running docker_test.go will build ahab, and build the docker files and validate the output is as expected
2. From the docker directory you can run the following

### Yum
```
GOOS=linux GOARCH=amd64 go build -o ahab ../main.go
docker build --no-cache -f yum/Dockerfile .
```

### APK
```
GOOS=linux GOARCH=amd64 go build -o ahab ../main.go
docker build --no-cache -f apk-autodetect/Dockerfile .
```

### Dpkg-query
```
GOOS=linux GOARCH=amd64 go build -o ahab ../main.go
docker build --no-cache -f dpkg-query/Dockerfile .
```
