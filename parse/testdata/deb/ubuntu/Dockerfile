# Usage: From project root directory:
# $ docker build --tag ahab-deb -f parse/testdata/deb/ubuntu/Dockerfile  .
# $ docker run -it ahab-deb
#
# From inside running docker container shell:
# # go build
# # dpkg-query --show --showformat='${Package} ${Version}\n' | ./ahab chase
# When done chasing, end via:
# # exit

FROM ubuntu:latest

RUN apt-get update
RUN apt-get install -y software-properties-common
RUN add-apt-repository -y ppa:longsleep/golang-backports
RUN apt-get update
RUN apt-get install -y golang-go vim

COPY ./ /data

VOLUME /data

WORKDIR /data
