# Usa la imagen base de Go
FROM golang:1.21

# Establece el directorio de trabajo
WORKDIR /app

# Instalar swag
RUN go install github.com/swaggo/swag/cmd/swag@latest

# Copia los archivos de dependencias primero
COPY go.mod go.sum ./
# Descarga las dependencias
RUN go mod download

# Copia el resto del código de la aplicación
COPY . .

# Limpia las dependencias innecesarias
RUN go mod tidy

# Genera la documentación Swagger
RUN swag init

# Compila la aplicación
RUN go build -o main .

# Expone el puerto 8080 para la aplicación
EXPOSE 8080

# Usa variables de entorno para configuración
ENV DB_HOST=postgres
ENV DB_PORT=5432
ENV DB_USER=root
ENV DB_PASSWORD=root
ENV DB_NAME=root

# Comando de inicio
CMD ["./main"]
