FROM golang:1.21.5 as build
WORKDIR /app
COPY . . 
RUN GOOS=linux CGO_ENABLED=0 go build -ldflags="-w -s" -o service_a .

FROM scratch
COPY --from=build /app/service_a .
ENTRYPOINT ["./service_a"]