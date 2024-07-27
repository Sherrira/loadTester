package main

import (
	"bytes"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"sync"
	"time"
)

type Result struct {
	statusCode int
	duration   time.Duration
}

func main() {
	// Configurar o logger para saída padrão e arquivo
	logFile, err := os.OpenFile("client.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("Erro ao abrir arquivo de log: %v", err)
	}
	defer logFile.Close()
	log.SetOutput(logFile)

	// Parsing CLI parameters
	urlFlag := flag.String("url", "", "URL do serviço a ser testado")
	requests := flag.Int("requests", 100, "Número total de requests")
	concurrency := flag.Int("concurrency", 10, "Número de chamadas simultâneas")
	method := flag.String("method", "GET", "Método HTTP a ser utilizado (GET, POST, PUT, DELETE)")
	headers := flag.String("headers", "", "Cabeçalhos HTTP no formato 'Chave:Valor,Chave:Valor'")
	body := flag.String("body", "", "Corpo da requisição para métodos POST e PUT")
	flag.Parse()

	if *urlFlag == "" {
		log.Println("A URL do serviço é obrigatória.")
		return
	}

	// Parse headers
	headerMap := make(map[string]string)
	if *headers != "" {
		headerPairs := strings.Split(*headers, ",")
		for _, pair := range headerPairs {
			parts := strings.SplitN(pair, ":", 2)
			if len(parts) == 2 {
				headerMap[strings.TrimSpace(parts[0])] = strings.TrimSpace(parts[1])
			}
		}
	}

	// Channel to collect results
	results := make(chan Result, *requests)
	var wg sync.WaitGroup

	// Channel to limit concurrency
	sem := make(chan struct{}, *concurrency)

	// Custom HTTP client
	client := &http.Client{}

	// Start time
	start := time.Now()

	// Launch goroutines
	for i := 0; i < *requests; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			sem <- struct{}{}        // Acquire a slot
			defer func() { <-sem }() // Release the slot

			req, err := http.NewRequest(*method, *urlFlag, bytes.NewBufferString(*body))
			if err != nil {
				log.Printf("Erro ao criar requisição: %v", err)
				results <- Result{statusCode: 0, duration: 0}
				return
			}

			// Add headers to the request
			for key, value := range headerMap {
				req.Header.Add(key, value)
			}

			startTime := time.Now()
			resp, err := client.Do(req)
			duration := time.Since(startTime)
			if err != nil {
				log.Printf("Erro ao fazer requisição: %v", err)
				results <- Result{statusCode: 0, duration: duration}
				return
			}
			results <- Result{statusCode: resp.StatusCode, duration: duration}
			resp.Body.Close()
		}()
	}

	// Wait for all goroutines to finish
	go func() {
		wg.Wait()
		close(results)
	}()

	// Initialize counters for the report
	totalRequests := 0
	successfulRequests := 0
	statusDistribution := make(map[int]int)

	// Process results
	for result := range results {
		totalRequests++
		if result.statusCode == 200 {
			successfulRequests++
		} else {
			statusDistribution[result.statusCode]++
		}
	}

	// Print total duration
	totalDuration := time.Since(start)
	fmt.Printf("Tempo total gasto na execução: %v\n", totalDuration)

	// Print report
	fmt.Printf("Quantidade total de requests realizados: %d\n", totalRequests)
	fmt.Printf("Quantidade de requests com status HTTP 200: %d\n", successfulRequests)
	fmt.Println("Distribuição de outros códigos de status HTTP:")
	for statusCode, count := range statusDistribution {
		fmt.Printf("  HTTP %d: %d\n", statusCode, count)
	}
}
