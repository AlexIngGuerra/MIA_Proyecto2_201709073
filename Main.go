package main

import (
	"MIA/analizador"
	"fmt"
)

func main() {
	fmt.Println("WALTER ALEXANDER GUERRA DUQUE 201709073")

	initPrueba()
}

func initPrueba() {
	analizador := analizador.NewAnalizador()
	var comando string
	comando = "exec -path=\"C:/Users/WALTER/Dropbox/Mi PC (DESKTOP-DUSC6PO)/Documents/1/Entradas/entrada2.script\""
	analizador.Analizar(comando)
}

func initProyecto() {
	analizador := analizador.NewAnalizador()
	var comando string
	for true {
		fmt.Scanln(&comando)
		analizador.Analizar(comando)
		if comando == "salir" {
			break
		}
	}
}
