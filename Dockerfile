FROM golang:1.14 as build

COPY . /project

WORKDIR /project

RUN make build


FROM ubuntu:latest

COPY --from=build /project/bin/server /

CMD ["./server"]