

# Alviss

## Introducción
Proyecto sencillo de acortador de URLs, escrito en Golang.

## Configuración y ejecución

Puedes usar fácilmente el siguiente comando para ejecutar el proyecto:
```
docker-compose -f deployments/docker-compose.yml up
```

De lo contrario, puedes instalar el proyecto con `go install` en el directorio `cmd/alviss/` y luego usar:
```
alviss runserver
```
Para ejecutar el proyecto en tu máquina; y también existe una bandera `-p` o `--port` para especificar el puerto del servidor.

**Nota: si vas a ejecutar el proyecto en tu máquina local, debes tener un `redis-server` ejecutándose en segundo plano.

## Rutas
### `localhost:8080/`:
Envía una solicitud GET y recibe una cálida bienvenida :)
```JSON
{
  "message": "Welcome to Alviss! Your mythical URL shortener."
}
```


### `localhost:8080/shorten`:
Realiza una petición POST con un objeto JSON como el siguiente y, a cambio, recibe el enlace corto generado:
```JSON
{
  "LongURL": "https://gist.github.com/joshbuchea/6f47e86d2510bce28f8e7f42ae84c716",
  "ExpTime": "2d"
}
```
Formato válido para la fecha de expiración: `2d` para 2 días, `2h` para 2 horas, `2m` para 2 minutos y `2s` para 2 segundos.

#### Resultado:
```JSON
{
  "message": "Short url created successfully",
  "ShortURL": "http://localhost:8080/ZLgJHJB2"
}
```

### `localhost:8080/url/{YOUR-SHORTENED-URL}`
Envía una solicitud GET y obtén los detalles de tu URL, como `UsedCount` o `ExpDate`:
```JSON
{
  "ExpDate": "2021-12-20T15:38:26.48860767Z",
  "OriginalURL": "https://gist.github.com/joshbuchea/6f47e86d2510bce28f8e7f42ae84c716",
  "ShortURL": "http://localhost:8080/ZLgJHJB2",
  "UsedCount": 3
}
```
