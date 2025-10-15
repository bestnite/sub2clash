FROM golang:1.25 as builder
LABEL authors="nite07"
WORKDIR /app
COPY . .
RUN go mod download
ARG version
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w -X github.com/bestnite/sub2clash/constant.Version=${version}" -o sub2clash .

FROM alpine:latest
WORKDIR /app
COPY --from=builder /app/sub2clash /app/sub2clash
ENTRYPOINT ["/app/sub2clash"]