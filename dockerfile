# Använd en officiell Go-bild som basbild
FROM golang:1.24-alpine AS builder

# Ställ in arbetskatalogen inuti containern
WORKDIR /app

# Kopiera go.mod och go.sum till arbetskatalogen
COPY go.mod go.sum ./

# Ladda ner beroenden
RUN go mod download

# Kopiera resten av applikationens källkod
COPY . .

# Bygg Go-applikationen till en binär fil med namnet "app"
ENV CGO_ENABLED=0
RUN go build -o /app/app ./main.go

# Använd en minimal basbild för den slutliga containern
FROM alpine:latest

# Ställ in arbetskatalogen inuti containern
WORKDIR /root/

# Kopiera den kompilerade binärfilen från byggsteget
COPY --from=builder /app/app .

# Kopiera styles-katalogen (om nödvändigt)
COPY --from=builder /app/styles ./styles

# Exponera port 8080 (eller den port som din app använder)
EXPOSE 8080

# Kommando för att köra applikationen
CMD ["./app"]
