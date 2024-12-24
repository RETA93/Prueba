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

# Tests unitarios
go test ./... -cover

# Tests de integración
go test ./tests/integration/... -tags=integration

# Tests de carga
k6 run tests/load/k6-test.js


# ---------------------------------------------------------------------------------------------------------------------

# Construye la imagen Docker: En la raíz de tu proyecto, ejecuta el siguiente comando:

docker build -t go-docker-project .


# Verifica que la imagen se construyó correctamente:

docker images

# ---------------------------------------------------------------------------------------------------------------------
# Autentica tu CLI de Docker con Google Cloud: Ejecuta este comando para autenticarte:

gcloud auth configure-docker

# Sube la imagen al Container Registry: Etiqueta la imagen de Docker con el nombre del repositorio de Google Container Registry y luego empújala:

docker tag go-docker-project gcr.io/pruebatecnicago/go-docker-project
docker push gcr.io/pruebatecnicago/go-docker-project

# Habilita Google Cloud Run: Si aún no lo has hecho, habilita Cloud Run en tu proyecto:

gcloud services enable run.googleapis.com

# Desplegar el contenedor a Cloud Run: Ejecuta el siguiente comando para desplegar tu imagen Docker en Cloud Run:
gcloud run deploy go-docker-project --image gcr.io/pruebatecnicago/go-docker-project --platform managed --region us-central1 --allow-unauthenticated

# ---------------------------------------------------------------------------------------------------------------------
# Para eliminar el servicio en Google Cloud Run, usa el siguiente comando:
gcloud run services delete go-docker-project --region us-central1

# Para verificar que el servicio se ha eliminado correctamente, puedes listar los servicios de Cloud Run:

gcloud run services list --region us-central1