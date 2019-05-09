FROM centos:7
MAINTAINER cesec

RUN yum install -y git
RUN yum install -y epel-release
RUN yum install -y go
#RUN export GOPATH=/go \
RUN mkdir -p /root/go/src
WORKDIR /root/go/src
RUN git clone https://github.com/zedundun/gobuster.git \
  && mkdir -p /root/go/src/golang.org/x \
  && cd /root/go/src/golang.org/x \
  && git clone https://github.com/golang/crypto.git \
  && git clone https://github.com/golang/sys.git \
  && cd /root/go/src/gobuster \
  && go get -u -v github.com/OJ/gobuster/gobusterdir \
  && go get -u -v github.com/OJ/gobuster/gobusterdns \
  && go get -u -v github.com/OJ/gobuster/libgobuster \
  && go build -o /bin/gobuster
RUN /bin/gobuster -fw -m dir -u http://www.baidu.com -w /root/go/src/gobuster/wordlist.txt -o output.txt

ENTRYPOINT top -b
