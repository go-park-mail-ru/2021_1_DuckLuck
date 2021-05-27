FROM golang:1.15 as build
WORKDIR /project
COPY go.mod .
RUN go mod download
COPY . /project
RUN make build

FROM golang:1.15 as api-server-build
WORKDIR /project
COPY go.mod .
RUN go mod download
COPY . /project
RUN make build api_server

FROM golang:1.15 as session-service-build
WORKDIR /project
COPY go.mod .
RUN go mod download
COPY . /project
RUN make build session_service

FROM golang:1.15 as cart-service-build
WORKDIR /project
COPY go.mod .
RUN go mod download
COPY . /project
RUN make build cart_service

FROM golang:1.15 as auth-service-build
WORKDIR /project
COPY go.mod .
RUN go mod download
COPY . /project
RUN make build auth_service

FROM ubuntu:latest as api-server
RUN apt update && apt install ca-certificates -y && rm -rf /var/cache/apt/*
COPY --from=api-server-build /project/bin/api_server /
CMD ["./api_server"]

FROM ubuntu:latest as session-service
RUN apt update && apt install ca-certificates -y && rm -rf /var/cache/apt/*
COPY --from=session-service-build /project/bin/session_service /
CMD ["./session_service"]

FROM ubuntu:latest as cart-service
RUN apt update && apt install ca-certificates -y && rm -rf /var/cache/apt/*
COPY --from=cart-service-build /project/bin/cart_service /
CMD ["./cart_service"]

FROM ubuntu:latest as auth-service
RUN apt update && apt install ca-certificates -y && rm -rf /var/cache/apt/*
COPY --from=auth-service-build /project/bin/auth_service /
CMD ["./auth_service"]

FROM postgres:13 as api-db
RUN apt update && \
    apt install myspell-ru -y
WORKDIR /usr/share/postgresql/13/tsearch_data
ENV DICT=/usr/share/hunspell/ru_RU
ENV POSTGRES_HOST_AUTH_METHOD=trust
RUN iconv -f koi8-r -t utf-8 -o russian.affix $DICT.aff && \
    iconv -f koi8-r -t utf-8 -o russian.dict $DICT.dic
