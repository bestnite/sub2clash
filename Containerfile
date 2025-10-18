FROM node:lts-alpine as frontend_builder
WORKDIR /app/server/frontend
COPY server/frontend/package*.json ./
RUN npm install
COPY server/frontend .
ARG version
RUN VITE_APP_VERSION=${version} npm run build

FROM golang:1.25 as builder
LABEL authors="nite07"
WORKDIR /app
COPY . .
COPY --from=frontend_builder /app/server/frontend/dist /app/server/frontend/dist
RUN go mod download
ARG version
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w -X github.com/bestnite/sub2clash/constant.Version=${version}" -o sub2clash .

FROM alpine:latest
WORKDIR /app
COPY --from=builder /app/sub2clash /app/sub2clash
ENTRYPOINT ["/app/sub2clash"]
