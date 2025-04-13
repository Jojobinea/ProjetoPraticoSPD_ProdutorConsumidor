package main

import (
	"testing"
)

func TestSistemaProdutorConsumidor(t *testing.T) {
	t.Log("Iniciando teste com 2 produtores, 2 consumidores e buffer de tamanho 3")
	cfg := Config{
		BufferSize:       3,
		NumProdutores:    2,
		NumConsumidores:  2,
		ItemsPorProdutor: 5,
	}

	IniciarSistema(cfg)
	t.Log("Sistema finalizado com sucesso.")
}

func TestBufferMaiorComMaisThreads(t *testing.T) {
	t.Log("Iniciando teste com 4 produtores, 4 consumidores e buffer de tamanho 10")
	cfg := Config{
		BufferSize:       10,
		NumProdutores:    4,
		NumConsumidores:  4,
		ItemsPorProdutor: 3,
	}

	IniciarSistema(cfg)
	t.Log("Sistema finalizado com sucesso.")
}
