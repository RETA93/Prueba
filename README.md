# Prueba


# Detener todos los contenedores
docker-compose down

# Eliminar el volumen
docker volume rm prueba_postgres_data

# Reconstruir y levantar los servicios
docker-compose up --build


# Eliminar la carpeta docs
rd /s /q docs

# O también puedes usar
del /f /q docs

# Regenerar la documentación
swag init

# Reconstruir y reiniciar la aplicación
docker-compose down
docker-compose up --build