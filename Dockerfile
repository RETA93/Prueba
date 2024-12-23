FROM golang:1.21

WORKDIR /app

# Instalar swag
RUN go install github.com/swaggo/swag/cmd/swag@latest

# Copiar los archivos de dependencias primero
COPY go.mod go.sum ./
RUN go mod download

# Copiar el resto del código
COPY . .

# Generar documentación Swagger
RUN swag init

# Compilar la aplicación
RUN go build -o main .

EXPOSE 3000

CMD ["./main"]