FROM golang:1.13

WORKDIR /root

RUN apt-get update

RUN curl https://raw.githubusercontent.com/golang/dep/master/install.sh | sh

ENV PATH $PATH:/root/go/bin

RUN echo 'export PS1="\$ "' >> /root/.bashrc \
&& echo 'export TEST_ENV="local"' >> /root/.bashrc

WORKDIR /go/src/github.com/paulhindemith

CMD ["/bin/bash"]


