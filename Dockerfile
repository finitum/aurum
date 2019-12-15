FROM golang as build
ENV GO111MODULE=on

WORKDIR /go/src/aurum
COPY . .
RUN go get -d
RUN CGO_ENABLED=0 go install -ldflags '-extldflags "-static"'

# Run stage
FROM scratch
COPY --from=build /go/bin/aurum /aurum
ENTRYPOINT ["/aurum"]
