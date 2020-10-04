FROM alpine:latest

WORKDIR /ahab-docker

RUN apk add --no-cache python3

COPY ahab .

# Spit out these just for easier debugging
RUN apk info -vv | sort

# Deprecated
RUN apk info -vv | sort | ./ahab chase --os alpine
# New way
RUN apk info -vv | sort | ./ahab chase --package-manager apk
