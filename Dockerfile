FROM golang:1.16.5-buster as builder
WORKDIR /app
COPY . /app
RUN make build

FROM scratch
WORKDIR /app
EXPOSE 8080
COPY --from=builder /app/_out/service /usr/bin/
ENTRYPOINT ["service"]