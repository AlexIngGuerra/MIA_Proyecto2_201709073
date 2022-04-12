package analizador

import (
	"MIA/comandos"
	"fmt"
	"strings"
)

type Analizador struct {
}

func NewAnalizador() Analizador {
	return Analizador{}
}

func (self Analizador) Analizar(entrada string) {
	linea := strings.Split(entrada, "\n")

	for i := 0; i < len(linea); i++ {

		if linea[i] != "" {

			fmt.Println(linea[i])
			comando := self.getComando(linea[i])

			if strings.ToLower(comando[0]) == "pause" {
				var tmp string
				fmt.Scanln(&tmp)

			} else if strings.ToLower(comando[0]) == "exec" {
				self.cmdExec(comando)

			} else if strings.ToLower(comando[0]) == "mkdisk" {
				self.cmdMkdisk(comando)
			}
		}

	}
}

func (self Analizador) getComando(linea string) []string {
	var comando []string
	comentario := strings.Split(linea, "#")
	comillas := strings.Split(comentario[0], "\"")

	for i := 0; i < len(comillas); i++ {

		if i%2 == 0 {
			//Quito espacios
			espacios := strings.Split(comillas[i], " ")
			for j := 0; j < len(espacios); j++ {

				iguales := strings.Split(espacios[j], "=")

				for k := 0; k < len(iguales); k++ {
					if iguales[k] != "" {
						comando = append(comando, iguales[k])
					}
				}
			}

		} else {
			//Es lo que estaba entre comillas y no se quitan espacios
			comando = append(comando, comillas[i])
		}

	}
	return comando
}

func (self Analizador) cmdExec(comando []string) {
	cmd := comandos.NewExec()

	for i := 1; i < len(comando); i++ {
		if i%2 != 0 && (i+1) < len(comando) {
			if strings.ToLower(comando[i]) == "-path" {
				cmd.Path = comando[i+1]
			} else {
				fmt.Println("Error: El comando exec no contiene el comando \"" + comando[i] + "\"")
			}
		}
	}

	valor := cmd.Ejecutar()
	self.Analizar(valor)
}

func (self Analizador) cmdMkdisk(comando []string) {
	cmd := comandos.NewMkdisk()

	for i := 1; i < len(comando); i++ {
		if i%2 != 0 && (i+1) < len(comando) {

		}
	}

	cmd.Ejecutar()
}
