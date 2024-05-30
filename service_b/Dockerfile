FROM golang:1.22.3 as build
WORKDIR /app
COPY . . 
RUN GOOS=linux CGO_ENABLED=0 go build -ldflags="-w -s" -o service_b .

FROM scratch
COPY --from=build /app/service_b .
ENTRYPOINT ["./service_b"]