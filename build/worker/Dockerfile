FROM golang:1.26-alpine AS build
WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 go build -o /bin/worker ./cmd/worker


FROM scratch
COPY --from=build /bin/worker /bin/worker

ENTRYPOINT ["/bin/worker"]
