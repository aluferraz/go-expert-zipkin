# Busca CEP  - Zipkin

O Serviço A roda no endpoint "/"

O Serviço B roda no endpoint "/servicoB"

Otel grpc exporter criado no go_expert_zipkin.go

Otel middleware injetado no servidor no arquivo [webserver.go](https://github.com/aluferraz/go-expert-zipkin/blob/45fff0294478ba61b5a7ca481aee9f9b93e0c1f3/internal/infra/web/webserver/webserver.go#L48C2-L49C1)

Otel client usado em todas as requisiçoes [ZipkinOtelClient.go](https://github.com/aluferraz/go-expert-zipkin/blob/45fff0294478ba61b5a7ca481aee9f9b93e0c1f3/internal/infra/http_clients/ZipkinOtelClient.go#L1)

Os logs ficam disponiveis em http://127.0.0.1:9411

![traces](./screenshot.png)


Variáveis de ambiente:

| Variável        | Descrição                                                         |
|-----------------|-------------------------------------------------------------------|
| WEBSERVER_PORT  | A porta em que o servidor web estará disponível.                  |
| WEATHER_API_KEY | Chave de API para acessar a API de clima.                         |
| WEATHER_API_URL | URL da API de clima para obter dados de temperatura.              |
| CEP_API_URL     | URL da API de CEP para obter dados de localidade a partir do CEP. |


O projeto usa Viper para gerenciar as váriaveis de ambiente, que podem ser configuradas no OS ou em um arquivo .env

O projeto possui testes integrados ao github actions e também pode ser testado com ``go test -v ./...``


# Como utilizar
Clonar o repositório

Preencha sua API_KEY no docker-compose:
[docker-compose.prod.yml](https://github.com/aluferraz/go-expert-zipkin/blob/3ba456c240eaf155cb748d7a21df7ef5133873c9/docker-compose.prod.yml#L10-L23)

Em desenvolvimento, voce também pode criar um arquivo .env na raíz do projeto:
```
WEATHER_API_KEY=mysupersecret
```

Para executar a versão em dev:
```
docker compose build --no-cache
docker compose up
```

Para executar a versão em prd (lembre-se de configurar a variavel de ambiente no docker-compose):

```
docker compose -f docker-compose.prod.yml build --no-cache

docker compose -f docker-compose.prod.yml up 
```
