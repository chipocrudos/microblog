# MicroBlog

Configurar variables de entorno:
```bash
MONGO_INITDB_ROOT_USERNAME=root
MONGO_INITDB_ROOT_PASSWORD=example

ME_CONFIG_MONGODB_ADMINUSERNAME=root
ME_CONFIG_MONGODB_ADMINPASSWORD=example
ME_CONFIG_MONGODB_URL=mongodb://root:example@mongo:27017/

DEBUG=true
PORT=8080
MONGO_URI=mongodb://root:example@deb:27017/?conect=direct
MONGO_DB=microblog
JWT_SALT=secret
# Minutes of expiration time
EXP_TIME=30
```


Ejecutar mongo
```bash
docker-composer up -d mongo
```

Ejecutar mongo-express
```bash
docker-composer up -d mongo-express
```

Ejecutar sin compilar dentro de la carpeta golang
```bash
cd golang
env $(cat ../.env)  go run cmd/**/*.go
```

Ejecutar dede docker-composer
```bash
docker-compose up microblog-go
```


Basado en el curso [Aprende lenguaje GO desde 0](https://www.udemy.com/share/102rbs3@sVC-jwM3qsG5aNG-Kyx4eJJoqJBTptkTxKT1O5tPcVI9W6OlLJiVVYSUrUfPJmXj/)

Github [Repositorio](https://github.com/ptilotta/twittor)

---
Apuntes:
**API Rest con Go (Golang)**
https://dev.to/orlmonteverde/api-rest-con-go-golang-y-postgresql-m0o