FROM google/golang:latest
ENV GOPATH=/gopath
ADD . /gopath/src/github.com/rlguarino/shortener