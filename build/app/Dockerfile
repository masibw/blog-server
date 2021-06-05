FROM golang:1.15.6 AS build

ENV GOOS=linux
ENV GOARCH=amd64
ENV CGO_ENABLED=0
ENV GO111MODULE=on

WORKDIR /go/src/github.com/masibw/blog-server

COPY go.mod  ./
RUN go mod download

COPY . .
RUN go build .

RUN go build tools/admin/admin.go

FROM gcr.io/distroless/static-debian10:latest
ENV DOCKERIZE_VERSION v0.6.1

ENV TZ=Asia/Tokyo

COPY ./migrations /migrations

COPY --from=build /go/src/github.com/masibw/blog-server/blog-server /
COPY --from=build /go/src/github.com/masibw/blog-server/admin /
ENTRYPOINT ["/blog-server"]