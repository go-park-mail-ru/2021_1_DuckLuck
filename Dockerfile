FROM golang:1.14 as build

COPY . /project

WORKDIR /project

RUN make build


FROM alphine:latest

COPY --from=build /project/bin/server /bin/

CMD ["server"]