# Prueba


# Detener todos los contenedores
docker-compose down

# Eliminar el volumen
docker volume rm prueba_postgres_data

# Reconstruir y levantar los servicios
docker-compose up --build