// server.go
package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"time"

	_ "github.com/mattn/go-sqlite3" // Import implícito do driver SQLite
)

// Estrutura para parsear a resposta da API de cotação.
// A API retorna um JSON com a chave "USDBRL", que contém os detalhes da cotação.
type ApiResponse struct {
	USDBRL ExchangeRate `json:"USDBRL"`
}

// Estrutura com os detalhes da cotação do dólar.
type ExchangeRate struct {
	Code       string `json:"code"`
	Codein     string `json:"codein"`
	Name       string `json:"name"`
	High       string `json:"high"`
	Low        string `json:"low"`
	VarBid     string `json:"varBid"`
	PctChange  string `json:"pctChange"`
	Bid        string `json:"bid"` // O valor que nos interessa
	Ask        string `json:"ask"`
	Timestamp  string `json:"timestamp"`
	CreateDate string `json:"create_date"`
}

const (
	apiURL     = "https://economia.awesomeapi.com.br/json/last/USD-BRL"
	apiTimeout = 200 * time.Millisecond
	dbTimeout  = 10 * time.Millisecond
	dbDriver   = "sqlite3"
	dbName     = "./cotacoes.db"
)

func main() {
	// Inicializa o banco de dados
	db, err := setupDatabase()
	if err != nil {
		log.Fatalf("Falha ao configurar o banco de dados: %v", err)
	}
	defer db.Close()

	// Registra o handler para o endpoint /cotacao
	http.HandleFunc("/cotacao", cotacaoHandler(db))

	log.Println("Servidor iniciado na porta 8080")
	// Inicia o servidor HTTP
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatalf("Falha ao iniciar o servidor: %v", err)
	}
}

// setupDatabase prepara o banco de dados SQLite, criando a tabela se não existir.
func setupDatabase() (*sql.DB, error) {
	db, err := sql.Open(dbDriver, dbName)
	if err != nil {
		return nil, err
	}

	createTableSQL := `
	CREATE TABLE IF NOT EXISTS cotacoes (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		bid TEXT,
		timestamp TEXT
	);`
	_, err = db.Exec(createTableSQL)
	if err != nil {
		return nil, err
	}
	return db, nil
}

// cotacaoHandler é uma closure que tem acesso à conexão do banco de dados.
func cotacaoHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Println("Requisição recebida em /cotacao")

		// Busca a cotação da API externa
		exchangeRate, err := fetchExchangeRate()
		if err != nil {
			log.Printf("Erro ao buscar cotação da API: %v", err)
			http.Error(w, "Erro ao buscar cotação externa", http.StatusInternalServerError)
			return
		}

		// Salva a cotação no banco de dados
		err = saveExchangeRate(db, exchangeRate)
		if err != nil {
			log.Printf("Erro ao salvar cotação no banco de dados: %v", err)
			// Não retornamos erro ao cliente aqui, pois a cotação foi obtida com sucesso.
			// O erro é apenas logado no servidor.
		}

		// Prepara a resposta para o cliente (apenas o campo "bid")
		response := struct {
			Bid string `json:"bid"`
		}{
			Bid: exchangeRate.Bid,
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(response)
		log.Println("Resposta enviada com sucesso para o cliente.")
	}
}

// fetchExchangeRate busca a cotação na API externa com um timeout de 200ms.
func fetchExchangeRate() (*ExchangeRate, error) {
	ctx, cancel := context.WithTimeout(context.Background(), apiTimeout)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, "GET", apiURL, nil)
	if err != nil {
		return nil, err
	}

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		// Verifica se o erro foi causado pelo timeout do contexto
		if ctx.Err() == context.DeadlineExceeded {
			log.Println("Erro: Timeout de 200ms excedido ao buscar cotação na API.")
		}
		return nil, err
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	var apiResponse ApiResponse
	err = json.Unmarshal(body, &apiResponse)
	if err != nil {
		return nil, err
	}

	return &apiResponse.USDBRL, nil
}

// saveExchangeRate salva a cotação no banco de dados com um timeout de 10ms.
func saveExchangeRate(db *sql.DB, rate *ExchangeRate) error {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	stmt, err := db.Prepare("INSERT INTO cotacoes(bid, timestamp) VALUES(?, ?)")
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.ExecContext(ctx, rate.Bid, time.Now().Format(time.RFC3339))
	if err != nil {
		// Verifica se o erro foi causado pelo timeout do contexto
		if ctx.Err() == context.DeadlineExceeded {
			log.Println("Erro: Timeout de 10ms excedido ao salvar no banco de dados.")
		}
		return err
	}

	log.Println("Cotação salva no banco de dados com sucesso.")
	return nil
}
