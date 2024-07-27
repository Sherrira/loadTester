# HTTP Client Load Tester
Este projeto é um cliente HTTP que permite realizar testes de carga em um serviço web. Ele suporta múltiplos métodos HTTP (GET, POST, PUT, DELETE) e permite a configuração de cabeçalhos HTTP e corpo da requisição via linha de comando.

## Funcionalidades
- Suporte para métodos HTTP: GET, POST, PUT, DELETE
- Configuração de cabeçalhos HTTP via linha de comando
- Configuração do corpo da requisição para métodos POST e PUT
- Controle de concorrência para limitar o número de requisições simultâneas
- Relatório detalhado com o tempo total de execução, quantidade de requisições realizadas e distribuição dos códigos de status HTTP

## Pré-requisitos
- Go
- Docker

## Como usar

### Execução local

Para executar o programa, utilize a seguinte sintaxe:
```sh
$ go run main.go \
--url=<URL> \
--requests=<NÚMERO_DE_REQUESTS> \
--concurrency=<NÍVEL_DE_CONCORRÊNCIA> \
--method=<MÉTODO_HTTP> \
--headers="<CHAVE_1>:<VALOR_1>,<CHAVE_2>:<VALOR_2>" \
--body="<CORPO_DA_REQUISIÇÃO>"
```

### Parâmetros 
- `--url`: UR do serviço a ser testado (obrigatório)
- `--requests`: Número total de requisições a serem realizadas (padrão: 100)
- `--concurrency`: Número de requisições simultâneas (padrão: 10)
- `--method`: Método HTTP a ser utilizado (GET, POST, PUT, DELETE) (padrão: GET)
- `--headers`: Cabeçalhos HTTP no formato Chave:Valor,Chave:Valor (opcional)
- `--body`: Corpo da requisição para métodos POST e PUT (opcional)

### Build da imagem Docker
Para compilar o programa, execute:
```sh
$ docker build -t loadtester .
```

### Exemplo executando o container
Para executar o programa:
```sh
$ docker run --rm loadtester \
--url=http://localhost:8080 \
--requests=500 \
--concurrency=50 \
--method=GET \
--headers="API_KEY:abc123"
```

### Exemplo apontando para um servidor em um container Docker na sua máquina
Caso você esteja executando a aplicação servidor com docker na sua máquina, para que o container do loadtester funcione chamando localhost, adicione o parâmetro de rede na chamada de execução:
```sh
$ docker run --rm --network=<NETWORK NAME OR ID> loadtester \
--url=http://<CONTAINER NAME OR ID>:8080 \
--requests=500 \
--concurrency=50 \
--method=GET \
--headers="API_KEY:abc123"
```

Para obter a rede de um container, caso esteja executando testes com docker, execute o comando a seguir:
```sh
$ docker inspect -f '{{json .NetworkSettings.Networks}}' <CONTAINER NAME OR ID>
```

### Exemplo de Uso

GET Request
```sh
$ docker run --rm loadtester --url=http://localhost:8080 --requests=50 --concurrency=5 --method=GET
```
POST Request com Cabeçalhos e Corpo
```sh
$ docker run --rm loadtester --url=http://localhost:8080 --requests=500 --concurrency=50 --headers="API_KEY:abc123" --method=POST --body='{"key":"value"}'
```

**Saída**

O programa gera um relatório com as seguintes informações:

```
Tempo total gasto na execução: <TEMPO>
Quantidade total de requisições realizadas: <QTD>
Quantidade de requisições com status HTTP 200: <QTD>
Distribuição de outros códigos de status HTTP:
<CODIGOS HTTP>: <QTD>
```