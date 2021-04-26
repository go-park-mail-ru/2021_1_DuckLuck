FROM golang:1.15 as build
COPY . /project
WORKDIR /project
RUN make build
RUN go test -coverprofile=coverage1.out -coverpkg=./... -cover ./...
RUN cat coverage1.out | grep -v mock > ./bin/cover.out
RUN go tool cover -func ./bin/cover.out

FROM ubuntu:latest as server
COPY --from=build /project/bin/ /
RUN apt update && apt install ca-certificates -y && rm -rf /var/cache/apt/*
CMD ["./server"]

FROM postgres:13 as postgres
RUN apt update && \
    apt install myspell-ru -y
WORKDIR /usr/share/postgresql/13/tsearch_data
ENV DICT=/usr/share/hunspell/ru_RU
ENV POSTGRES_HOST_AUTH_METHOD=trust
RUN iconv -f koi8-r -t utf-8 -o russian.affix $DICT.aff && \
    iconv -f koi8-r -t utf-8 -o russian.dict $DICT.dic

