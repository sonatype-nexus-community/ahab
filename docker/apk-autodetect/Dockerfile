FROM alpine:latest

WORKDIR /ahab-docker

#Requirement to run golang binary in apline
RUN apk add --no-cache libc6-compat

RUN apk add --no-cache python3

COPY ahab .

# Spit out these just for easier debugging
RUN apk info -vv | sort

RUN apk info -vv | sort | ./ahab chase