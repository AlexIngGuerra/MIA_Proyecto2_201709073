package analizador

import (
	"MIA/comandos"
	"fmt"
	"strconv"
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
			if len(comando) <= 0 {
				continue
			}

			if strings.ToLower(comando[0]) == "pause" || strings.ToLower(linea[i]) == "linea" {
				var comando string
				fmt.Scanln(&comando)

			} else if strings.ToLower(comando[0]) == "exec" {
				self.cmdExec(comando)

			} else if strings.ToLower(comando[0]) == "mkdisk" {
				self.cmdMkdisk(comando)

			} else if strings.ToLower(comando[0]) == "rmdisk" {
				self.cmdRmdisk(comando)

			} else if strings.ToLower(comando[0]) == "fdisk" {
				self.cmdFdisk(comando)
			}

		}

	} //Fin For de lineas
}

//Obtener el comando
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

//Comando Exec
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
	fmt.Print("\n")
	self.Analizar(valor)
}

//Comando Mkdisk
func (self Analizador) cmdMkdisk(comando []string) {
	cmd := comandos.NewMkdisk()

	for i := 1; i < len(comando); i++ {
		if i%2 != 0 && (i+1) < len(comando) {

			if strings.ToLower(comando[i]) == "-size" {

				valor, err := strconv.Atoi(comando[i+1])
				if err != nil || valor < 1 {
					fmt.Println("Error: El parametro size solo puede contener numeros enteros positivos")
				} else {
					cmd.Size = valor
				}

			} else if strings.ToLower(comando[i]) == "-fit" {

				valor := strings.ToUpper(comando[i+1])

				if valor == "FF" || valor == "WF" || valor == "BF" {
					cmd.Fit = valor
				} else {
					fmt.Println("Error: El valor del parametro fit es incorrecto")
					continue
				}

			} else if strings.ToLower(comando[i]) == "-unit" {
				valor := strings.ToUpper(comando[i+1])

				if valor == "M" || valor == "K" {
					cmd.Unit = valor
				} else {
					fmt.Println("Error: El valor del parametro unit es incorrecto")
					continue
				}

			} else if strings.ToLower(comando[i]) == "-path" {
				cmd.Path = comando[i+1]
			}

		}
	}

	cmd.Ejecutar()
}

//Comando Rmdisk
func (self Analizador) cmdRmdisk(comando []string) {
	cmd := comandos.NewMkdisk()

	for i := 1; i < len(comando); i++ {
		if i%2 != 0 && (i+1) < len(comando) {

			if strings.ToLower(comando[i]) == "-path" {
				cmd.Path = comando[i+1]
			} else {
				fmt.Println("Error: El comando exec no contiene el comando \"" + comando[i] + "\"")
			}

		}
	}

	cmd.EjecutarRmdisk()
}

//COMANDO FDISK
func (self Analizador) cmdFdisk(comando []string) {
	cmd := comandos.NewFdisk()

	for i := 1; i < len(comando); i++ {
		if i%2 != 0 && (i+1) < len(comando) {

			if strings.ToLower(comando[i]) == "-size" {

				valor, err := strconv.Atoi(comando[i+1])
				if err != nil || valor < 1 {
					fmt.Println("Error: El parametro size solo puede contener numeros enteros positivos")
				} else {
					cmd.Size = valor
				}

			} else if strings.ToLower(comando[i]) == "-unit" {

				valor := strings.ToUpper(comando[i+1])

				if valor == "M" || valor == "K" || valor == "B" {
					cmd.Unit = valor
				} else {
					fmt.Println("Error: El valor del parametro unit es incorrecto")
					continue
				}

			} else if strings.ToLower(comando[i]) == "-path" {

				cmd.Path = comando[i+1]

			} else if strings.ToLower(comando[i]) == "-type" {

				valor := strings.ToUpper(comando[i+1])

				if valor == "P" || valor == "E" || valor == "L" {
					cmd.Type = valor
				} else {
					fmt.Println("Error: El valor del parametro type es incorrecto")
					continue
				}

			} else if strings.ToLower(comando[i]) == "-fit" {

				valor := strings.ToUpper(comando[i+1])

				if valor == "FF" || valor == "WF" || valor == "BF" {
					cmd.Fit = valor
				} else {
					fmt.Println("Error: El valor del parametro fit es incorrecto")
					continue
				}

			} else if strings.ToLower(comando[i]) == "-name" {
				cmd.Name = comando[i+1]

			} else {
				fmt.Println("Error: El comando exec no contiene el comando \"" + comando[i] + "\"")
			}

		}
	}

	cmd.Ejecutar()
}
