FROM golang:1.15 as build
COPY . /project
WORKDIR /project
RUN make build
RUN go test -coverprofile=coverage1.out -coverpkg=./... -cover ./...
RUN cat coverage1.out | grep -v mock > ./bin/cover.out
RUN go tool cover -func ./bin/cover.out

FROM ubuntu:latest
COPY --from=build /project/bin/ /
RUN apt update && apt install ca-certificates -y && rm -rf /var/cache/apt/*
CMD ["./server"]
