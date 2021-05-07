FROM golang:1.15 as build
COPY . /project
WORKDIR /project
RUN make build
RUN go test -coverprofile=coverage1.out -coverpkg=./... -cover ./...
RUN cat coverage1.out | grep -v mock | grep -v proto | grep -v cmd > cover.out
RUN go tool cover -func cover.out

FROM golang:1.15 as api-server-build
COPY . /project
WORKDIR /project
RUN make build api_server

FROM golang:1.15 as session-service-build
COPY . /project
WORKDIR /project
RUN make build session_service

FROM golang:1.15 as cart-service-build
COPY . /project
WORKDIR /project
RUN make build cart_service

FROM golang:1.15 as auth-service-build
COPY . /project
WORKDIR /project
RUN make build auth_service

FROM ubuntu:latest as api-server
COPY --from=api-server-build /project/bin/api_server /
RUN apt update && apt install ca-certificates -y && rm -rf /var/cache/apt/*
CMD ["./api_server"]

FROM ubuntu:latest as session-service
COPY --from=session-service-build /project/bin/session_service /
RUN apt update && apt install ca-certificates -y && rm -rf /var/cache/apt/*
CMD ["./session_service"]

FROM ubuntu:latest as cart-service
COPY --from=cart-service-build /project/bin/cart_service /
RUN apt update && apt install ca-certificates -y && rm -rf /var/cache/apt/*
CMD ["./cart_service"]

FROM ubuntu:latest as auth-service
COPY --from=auth-service-build /project/bin/auth_service /
RUN apt update && apt install ca-certificates -y && rm -rf /var/cache/apt/*
CMD ["./auth_service"]

FROM postgres:13 as api-db
RUN apt update && \
    apt install myspell-ru -y
WORKDIR /usr/share/postgresql/13/tsearch_data
ENV DICT=/usr/share/hunspell/ru_RU
ENV POSTGRES_HOST_AUTH_METHOD=trust
RUN iconv -f koi8-r -t utf-8 -o russian.affix $DICT.aff && \
    iconv -f koi8-r -t utf-8 -o russian.dict $DICT.dic

