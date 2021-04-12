FROM golang:1.15 as build
COPY . /project
WORKDIR /project
RUN make build

FROM ubuntu:latest
COPY --from=build /project/bin/ /
RUN apt update && apt install ca-certificates -y && rm -rf /var/cache/apt/*
CMD ["./server"]
