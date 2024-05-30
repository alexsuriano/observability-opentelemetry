# Desafio Observabilidade & Open Telemetry

## Objetivo: 
Desenvolver um sistema em Go que receba um CEP, identifica a cidade e retorna o clima atual (temperatura em graus celsius, fahrenheit e kelvin) juntamente com a cidade. Esse sistema deverá implementar OTEL(Open Telemetry) e Zipkin.
Basedo no cenário conhecido "Sistema de temperatura por CEP" denominado Serviço B, será incluso um novo projeto, denominado Serviço A.

 ### Requisitos: Serviço A (responsável pelo input)

O sistema deve receber um input de 8 dígitos via POST, através do schema:  { "cep": "29902555" }
O sistema deve validar se o input é valido (contem 8 dígitos) e é uma STRING
Caso seja válido, será encaminhado para o Serviço B via HTTP
Caso não seja válido, deve retornar:
Código HTTP: 422
Mensagem: invalid zipcode

### Requisitos: Serviço B (responsável pela orquestração):

O sistema deve receber um CEP válido de 8 digitos
O sistema deve realizar a pesquisa do CEP e encontrar o nome da localização, a partir disso, deverá retornar as temperaturas e formata-lás em: Celsius, Fahrenheit, Kelvin juntamente com o nome da localização.
O sistema deve responder adequadamente nos seguintes cenários:
Em caso de sucesso:
Código HTTP: 200
Response Body: { "city: "São Paulo", "temp_C": 28.5, "temp_F": 28.5, "temp_K": 28.5 }
Em caso de falha, caso o CEP não seja válido (com formato correto):
Código HTTP: 422
Mensagem: invalid zipcode
​​​Em caso de falha, caso o CEP não seja encontrado:
Código HTTP: 404
Mensagem: can not find zipcode

Após a implementação dos serviços, adicione a implementação do OTEL + Zipkin:

Implementar tracing distribuído entre Serviço A - Serviço B
Utilizar span para medir o tempo de resposta do serviço de busca de CEP e busca de temperatura

### Como rodar a aplicação:
- Clone o projeto
- Vá até raiz do projeto
- Execute o seguinte comando para subir os 2 serviços e o zipkin:
```
docker compose up
```
- Você pode executar um `POST` para `http://localhost:8181` usando o seguinte schema json como body da requisiçãp:

  `{ "cep": "29902555" }`
  
- Para facilitar a chamada, a pasta `api` contém um arquivo `temperatureCity.http` que faz uma chamada de exemplo, sinta-se a vontade para trocar o CEP ao qual desejar.
- Você pode acessar o painel do Zipkin no seguinte endereço: `http://localhost:9411`
