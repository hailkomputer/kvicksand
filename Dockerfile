FROM golang:1.17-alpine as builder

ADD . /go/src/github.com/hailkomputer/kvicksand
WORKDIR /go/src/github.com/hailkomputer/kvicksand
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o kvicksand ./cmd/kvicksand

FROM gcr.io/distroless/static
COPY --from=builder /go/src/github.com/hailkomputer/kvicksand /bin
ENV GODEBUG madvdontneed=1
ENTRYPOINT ["/bin/kvicksand"]