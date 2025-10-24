# Desafio: Cotação Dólar-Real com Go

Este repositório contém a solução para um desafio de programação em Go que envolve a criação de um sistema cliente-servidor para consultar a cotação atual do Dólar (USD) em relação ao Real (BRL).

O sistema é composto por dois executáveis:
1.  `server.go`: Um servidor HTTP que busca a cotação de uma API externa, persiste o resultado em um banco de dados SQLite e responde ao cliente.
2.  `client.go`: Um cliente que faz uma requisição ao servidor local, recebe a cotação e a salva em um arquivo de texto.

O principal objetivo do desafio é demonstrar o uso de `context` para gerenciar timeouts em operações de I/O (requisições HTTP e transações de banco de dados), além da manipulação de JSON e arquivos.

## ✨ Funcionalidades

* **Servidor HTTP**: Expõe o endpoint `/cotacao` na porta `8080`.
* **Consumo de API Externa**: O servidor busca dados da API `https://economia.awesomeapi.com.br`.
* **Persistência de Dados**: Cada cotação obtida é salva em um banco de dados **SQLite** (`cotacoes.db`).
* **Gerenciamento de Timeouts com `context`**:
    * **Servidor (API)**: Timeout de **200ms** para a chamada à API externa.
    * **Servidor (Banco de Dados)**: Timeout de **10ms** para a operação de escrita no banco de dados.
    * **Cliente**: Timeout de **300ms** para receber a resposta do servidor.
* **Manipulação de Arquivos**: O cliente salva a cotação recebida no arquivo `cotacao.txt`.
* **Logging**: O sistema registra em log os eventos importantes e os erros, especialmente os de timeout.

## 🔧 Pré-requisitos

Para executar este projeto, você precisa ter o **Go** (versão 1.18 ou superior) instalado em sua máquina.

* [Página oficial de download do Go](https://go.dev/doc/install)

## 🚀 Como Executar

Siga os passos abaixo para executar a aplicação.

### 1. Estrutura dos Arquivos

Certifique-se de que sua pasta de projeto contém os arquivos `server.go` e `client.go`.

### 2. Instale as Dependências

Abra um terminal na pasta do projeto e execute os seguintes comandos para inicializar o módulo e baixar as dependências.

```sh
# Inicializa o gerenciador de módulos
go mod init desafio-cotacao
```
# Baixa e instala as dependências (driver do SQLite)
go get [github.com/mattn/go-sqlite3](https://github.com/mattn/go-sqlite3)

### 3. Execute o Servidor
Em um terminal, inicie o servidor. Ele ficará em execução, aguardando requisições do cliente.
```sh
go run server.go
Bash

go run server.go
# Você deverá ver a seguinte mensagem, indicando que o servidor está pronto:
# 2025/10/24 18:55:00 Servidor iniciado na porta 8080
```
### 4. Execute o Cliente
Abra um novo terminal (mantendo o terminal do servidor em execução) e rode o cliente.
```sh
Bash

go run client.go
# Se tudo ocorrer bem, o cliente fará a requisição, receberá a resposta e salvará o resultado. Você verá a seguinte saída:

#2025/10/24 18:56:00 Cliente iniciado. Solicitando cotação...
#2025/10/24 18:56:00 Cotação salva com sucesso em cotacao.txt. Valor: 5.1234
```
**Verificação**
Após a execução, os seguintes arquivos serão criados na raiz do projeto:

**cotacao.txt:** Contém a cotação no formato Dólar: {valor}.

**cotacoes.db:** O banco de dados SQLite contendo o histórico de todas as cotações salvas.

**⚠️ Possíveis Erros e Soluções**
Este projeto foi desenhado para testar cenários de timeout. Abaixo estão os erros mais comuns que podem ocorrer e o que eles significam.

Erros no Log do Servidor
**1. Timeout na API externa**
Snippet de código

**Erro:** Timeout de 200ms excedido ao buscar cotação na API.
Causa: A API economia.awesomeapi.com.br demorou mais de 200 milissegundos para responder.

**Comportamento:** O servidor não conseguirá obter a cotação e retornará um erro HTTP 500 para o cliente.

**2. Timeout no Banco de Dados**
Snippet de código

**Erro:** Timeout de 10ms excedido ao salvar no banco de dados.
Causa: A operação de escrita no arquivo cotacoes.db demorou mais de 10 milissegundos.

**Comportamento:** O servidor obterá a cotação, mas falhará em salvá-la. A resposta ainda será enviada com sucesso para o cliente.

**Erros no Log do Cliente**
**1. Timeout para receber a resposta do servidor**
Snippet de código

**Erro:** Timeout de 300ms excedido para receber o resultado do servidor.
exit status 1
**Causa:** O tempo total da operação (requisição + processamento do servidor + resposta) ultrapassou 300 milissegundos.

**Solução:** Este é um comportamento esperado para demonstrar o context funcionando.

**2. Conexão Recusada**
Snippet de código

**Erro:** ... dial tcp 127.0.0.1:8080: connect: connection refused
exit status 1
**Causa:** O server.go não está em execução.

**Solução:** Certifique-se de que o servidor foi iniciado no outro terminal.