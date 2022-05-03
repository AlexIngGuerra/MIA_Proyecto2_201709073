package main

import (
	"MIA/analizador"
	"bufio"
	"fmt"
	"os"
	"strings"
)

func main() {
	fmt.Println("WALTER ALEXANDER GUERRA DUQUE 201709073")
	initProyecto()
}

func initProyecto() {
	analizador := analizador.NewAnalizador()
	reader := bufio.NewReader(os.Stdin)
	for {
		cmd, _ := reader.ReadString('\n')
		cmd = strings.Replace(cmd, "\n", "", -1)
		cmd = strings.Trim(cmd, string(uint8(13)))
		if cmd == "salir" {
			break
		}

		analizador.Analizar(cmd)

	}
}
