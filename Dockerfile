# Usar una imagen base de Go para construir el binario
FROM golang:1.22 as builder

# Establecer el directorio de trabajo dentro del contenedor
WORKDIR /app

# Copiar los archivos del proyecto al directorio de trabajo
COPY . .

# Descargar las dependencias y construir el binario
RUN go mod download
RUN go build -o main .

FROM ubuntu:22.04

# Instalar las dependencias necesarias
RUN apt-get update && apt-get install -y ca-certificates && rm -rf /var/lib/apt/lists/*

# Establecer el directorio de trabajo dentro del contenedor
WORKDIR /app

# Copiar el binario desde la etapa de construcción
COPY --from=builder /app/main .
COPY --from=builder /app/services /app/services

# Exponer el puerto en el que corre tu aplicación, si es necesario
EXPOSE 8080

# Comando para ejecutar el binario
CMD ["./main"]
