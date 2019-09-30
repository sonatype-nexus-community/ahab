FROM centos:latest

RUN yum -y install epel-release && yum -y install python-pip

COPY ahab .

RUN yum list installed
RUN yum list installed | ./ahab chase -os centos