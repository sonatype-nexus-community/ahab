FROM centos:latest

WORKDIR /ahab-docker

RUN yum -y install epel-release which python3-pip

COPY ahab .

# Spit out these just for easier debugging
RUN yum list installed

RUN yum list installed | ./ahab chase

