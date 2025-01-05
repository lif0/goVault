FROM golang:1.23.2 AS builder

LABEL org.opencontainers.image.title="goVault"
LABEL org.opencontainers.image.description="Blazing-fast in-memory database written in Go."
LABEL org.opencontainers.image.url="https://github.com/lif0/goVault"
LABEL org.opencontainers.image.logo="https://raw.githubusercontent.com/lif0/goVault/refs/heads/main/.github/assets/goVault_poster_round.png"

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o server ./cmd/server/main.go

FROM scratch

WORKDIR /app

COPY --from=builder /app/server .

CMD ["./server"]