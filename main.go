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

// Função de log com timestamp
func logf(format string, a ...interface{}) {
	fmt.Printf("[%s] %s", time.Now().Format("15:04:05.000"), fmt.Sprintf(format, a...))
}

// Função principal do sistema Produtor-Consumidor
func IniciarSistema(cfg Config) {
	rand.Seed(time.Now().UnixNano())
	buffer := make(chan int, cfg.BufferSize)

	var wgProd sync.WaitGroup
	var wgCons sync.WaitGroup

	start := time.Now() // ⏱️ Início da medição

	// Monitor de buffer
	go func() {
		ticker := time.NewTicker(500 * time.Millisecond)
		defer ticker.Stop()
		for range ticker.C {
			logf("[Monitor] Tamanho atual do buffer: %d / %d\n", len(buffer), cap(buffer))
		}
	}()

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

	// Esperar produtores
	wgProd.Wait()
	close(buffer) // Sinaliza fim da produção
	wgCons.Wait() // Esperar consumidores

	elapsed := time.Since(start)
	logf("✅ Tempo total de execução: %s\n", elapsed)
}

// Produtor
func produtor(id int, buffer chan int, total int, wg *sync.WaitGroup) {
	defer wg.Done()
	for i := 0; i < total; i++ {
		item := rand.Intn(1000)
		logf("[Produtor %d] Produzindo item: %d\n", id, item)
		buffer <- item
		time.Sleep(time.Millisecond * time.Duration(rand.Intn(200)))
	}
	logf("[Produtor %d] Finalizou produção.\n", id)
}

// Consumidor
func consumidor(id int, buffer chan int, wg *sync.WaitGroup) {
	defer wg.Done()
	for item := range buffer {
		logf("  [Consumidor %d] Consumiu item: %d\n", id, item)
		time.Sleep(time.Millisecond * time.Duration(rand.Intn(500)))
	}
	logf("  [Consumidor %d] Finalizou consumo.\n", id)
}

// Função principal para executar o sistema
func main() {

	logf("***Caso 1: Produtores rapidos, consumidores lentos***\n")
	cfg1 := Config{
		BufferSize:       3,
		NumProdutores:    2,
		NumConsumidores:  1,
		ItemsPorProdutor: 5,
	}
	IniciarSistema(cfg1)
	logf("\n\n")
	//💡 Esperado:
	//Buffer enche rápido.
	//Consumidor fica sobrecarregado.

	logf("***Caso 2: Consumidores rapidos, produtores lentos***\n")
	cfg2 := Config{
		BufferSize:       5,
		NumProdutores:    1,
		NumConsumidores:  3,
		ItemsPorProdutor: 10,
	}
	IniciarSistema(cfg2)
	logf("\n\n")
	//💡 Esperado:
	//Buffer quase sempre vazio.
	//Consumidores competem pelos poucos itens.
	//Alguns consumidores podem terminar bem antes.

	logf("***Caso 3: Buffer pequeno, muitos produtores***\n")
	cfg3 := Config{
		BufferSize:       2,
		NumProdutores:    4,
		NumConsumidores:  2,
		ItemsPorProdutor: 4,
	}
	IniciarSistema(cfg3)
	logf("\n\n")
	//💡 Esperado:
	//Muita contenção no buffer.
	//Muitos produtores vão aguardar espaço para escrever.
	//Buffer frequentemente no limite.

	logf("***Caso 4: Buffer grande, produtores e consumidores equilibrados***\n")
	cfg4 := Config{
		BufferSize:       20,
		NumProdutores:    3,
		NumConsumidores:  3,
		ItemsPorProdutor: 10,
	}
	IniciarSistema(cfg4)
	logf("\n\n")
	//💡 Esperado:
	//Sistema flui com pouco bloqueio.
	//Boa concorrência.
	//Ideal para testar eficiência máxima.

	logf("***Caso 5: Teste de escala — produção alta***\n")
	cfg5 := Config{
		BufferSize:       50,
		NumProdutores:    10,
		NumConsumidores:  5,
		ItemsPorProdutor: 100,
	}
	IniciarSistema(cfg5)
	logf("\n\n")
	//💡 Esperado:
	//Estresse do sistema com alta carga.
	//Útil para observar se há travamentos ou lentidão.
	//Ideal para medir tempo total com muitos dados.

	logf("***Caso 6: Apenas 1 produtor e 1 consumidor***\n")
	cfg6 := Config{
		BufferSize:       1,
		NumProdutores:    1,
		NumConsumidores:  1,
		ItemsPorProdutor: 5,
	}
	IniciarSistema(cfg6)
	logf("\n\n")
	//💡 Esperado:
	//Tudo sequencial (quase sem concorrência).
	//Útil para ver o sistema mais simples possível.
	//Comportamento previsível.
}
