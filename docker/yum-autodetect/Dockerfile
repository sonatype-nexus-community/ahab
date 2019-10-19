FROM centos:latest

RUN yum -y install epel-release && yum -y install python3-pip

COPY ahab .

# Spit out these just for easier debugging
RUN yum list installed

RUN yum list installed | ./ahab chase