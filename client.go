package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"time"
)

// Estrutura para receber a resposta do servidor local.
type ServerResponse struct {
	Bid string `json:"bid"`
}

const (
	serverURL      = "http://localhost:8080/cotacao"
	requestTimeout = 300 * time.Millisecond
	outputFile     = "cotacao.txt"
)

func main() {
	log.Println("Cliente iniciado. Solicitando cotação...")

	// 1. Cria um contexto com timeout de 300ms.
	ctx, cancel := context.WithTimeout(context.Background(), requestTimeout)
	defer cancel()

	// 2. Cria a requisição HTTP com o contexto.
	req, err := http.NewRequestWithContext(ctx, "GET", serverURL, nil)
	if err != nil {
		log.Fatalf("Erro ao criar requisição: %v", err)
	}

	// 3. Executa a requisição.
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		// Verifica se o erro foi causado pelo cancelamento do contexto (timeout).
		if ctx.Err() == context.DeadlineExceeded {
			log.Fatalf("Erro: Timeout de 300ms excedido para receber o resultado do servidor.")
		}
		log.Fatalf("Erro ao fazer requisição ao servidor: %v", err)
	}
	defer res.Body.Close()

	// 4. Lê o corpo da resposta.
	body, err := io.ReadAll(res.Body)
	if err != nil {
		log.Fatalf("Erro ao ler a resposta do servidor: %v", err)
	}

	// 5. Faz o parse do JSON da resposta.
	var serverResponse ServerResponse
	if err := json.Unmarshal(body, &serverResponse); err != nil {
		log.Fatalf("Erro ao fazer parse da resposta JSON: %v", err)
	}

	// 6. Salva a cotação no arquivo.
	if err := saveCotacaoToFile(serverResponse.Bid); err != nil {
		log.Fatalf("Erro ao salvar cotação no arquivo: %v", err)
	}

	log.Printf("Cotação salva com sucesso em %s. Valor: %s", outputFile, serverResponse.Bid)
}

// saveCotacaoToFile cria ou sobrescreve o arquivo cotacao.txt com o valor da cotação.
func saveCotacaoToFile(bid string) error {
	content := fmt.Sprintf("Dólar: %s", bid)

	// A função WriteFile lida com a criação, abertura, escrita e fechamento do arquivo.
	// O terceiro parâmetro (0644) define as permissões do arquivo.
	return os.WriteFile(outputFile, []byte(content), 0644)
}
