FROM shortener/base:latest
ADD . /gopath/src/github.com/rlguarino/shortener/frontend
WORKDIR /gopath/src/github.com/rlguarino/shortener/frontend
RUN go get && go build
CMD ["./frontend"]