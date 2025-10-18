// client.go
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

	// Cria um contexto com timeout de 300ms para a requisição.
	ctx, cancel := context.WithTimeout(context.Background(), requestTimeout)
	defer cancel()

	// Cria a requisição HTTP com o contexto.
	req, err := http.NewRequestWithContext(ctx, "GET", serverURL, nil)
	if err != nil {
		log.Fatalf("Erro ao criar requisição: %v", err)
	}

	// Executa a requisição.
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		// Verifica se o erro foi um timeout do contexto.
		if ctx.Err() == context.DeadlineExceeded {
			log.Fatalf("Erro: Timeout de 300ms excedido para receber o resultado do servidor.")
		}
		log.Fatalf("Erro ao fazer requisição ao servidor: %v", err)
	}
	defer res.Body.Close()

	// Lê o corpo da resposta.
	body, err := io.ReadAll(res.Body)
	if err != nil {
		log.Fatalf("Erro ao ler a resposta do servidor: %v", err)
	}

	// Faz o parse do JSON da resposta.
	var serverResponse ServerResponse
	err = json.Unmarshal(body, &serverResponse)
	if err != nil {
		log.Fatalf("Erro ao fazer parse da resposta JSON: %v", err)
	}

	// Salva a cotação no arquivo.
	err = saveCotacaoToFile(serverResponse.Bid)
	if err != nil {
		log.Fatalf("Erro ao salvar cotação no arquivo: %v", err)
	}

	log.Printf("Cotação salva com sucesso em %s. Valor: %s", outputFile, serverResponse.Bid)
}

// saveCotacaoToFile salva o valor da cotação em cotacao.txt.
func saveCotacaoToFile(bid string) error {
	// Cria a string no formato especificado.
	content := fmt.Sprintf("Dólar: %s", bid)

	// Escreve o conteúdo no arquivo. os.WriteFile cria o arquivo se não existir
	// ou o sobrescreve se já existir.
	return os.WriteFile(outputFile, []byte(content), 0644)
}
