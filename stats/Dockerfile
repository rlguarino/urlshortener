FROM shortener/base:latest
ADD . /gopath/src/github.com/rlguarino/shortener/stats
WORKDIR /gopath/src/github.com/rlguarino/shortener/stats
RUN go get && go build
CMD ["./stats"]