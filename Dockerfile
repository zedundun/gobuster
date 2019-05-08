FROM centos:7
MAINTAINER The CentOS Project
ENTRYPOINT ["top", "-b"]
CMD echo "Hello World"
