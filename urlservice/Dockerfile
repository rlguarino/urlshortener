FROM shortener/base:latest
ADD . /gopath/src/github.com/rlguarino/shortener/urlservice
WORKDIR /gopath/src/github.com/rlguarino/shortener/urlservice
RUN go get && go build
CMD ["./urlservice"]