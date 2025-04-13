package main

import (
	"fmt"
	"math/rand"
	"sync"
	"time"
)

// Configuração do sistema
type Config struct {
	BufferSize       int
	NumProdutores    int
	NumConsumidores  int
	ItemsPorProdutor int
}

// Função para iniciar o sistema Produtor-Consumidor com goroutines
func IniciarSistema(cfg Config) {
	rand.Seed(time.Now().UnixNano())
	buffer := make(chan int, cfg.BufferSize)

	var wgProd sync.WaitGroup
	var wgCons sync.WaitGroup

	// Iniciar produtores
	for i := 0; i < cfg.NumProdutores; i++ {
		wgProd.Add(1)
		go produtor(i+1, buffer, cfg.ItemsPorProdutor, &wgProd)
	}

	// Iniciar consumidores
	for i := 0; i < cfg.NumConsumidores; i++ {
		wgCons.Add(1)
		go consumidor(i+1, buffer, &wgCons)
	}

	// Esperar todos os produtores terminarem
	wgProd.Wait()
	// Fechar o buffer para indicar fim da produção
	close(buffer)
	// Esperar todos os consumidores terminarem
	wgCons.Wait()
}

func produtor(id int, buffer chan int, total int, wg *sync.WaitGroup) {
	defer wg.Done()
	for i := 0; i < total; i++ {
		item := rand.Intn(1000)
		fmt.Printf("[Produtor %d] Produzindo item: %d\n", id, item)
		buffer <- item
		time.Sleep(time.Millisecond * time.Duration(rand.Intn(200)))
	}
	fmt.Printf("[Produtor %d] Finalizou produção.\n", id)
}

func consumidor(id int, buffer chan int, wg *sync.WaitGroup) {
	defer wg.Done()
	for item := range buffer {
		fmt.Printf("  [Consumidor %d] Consumiu item: %d\n", id, item)
		time.Sleep(time.Millisecond * time.Duration(rand.Intn(500)))
	}
	fmt.Printf("  [Consumidor %d] Finalizou consumo.\n", id)
}
