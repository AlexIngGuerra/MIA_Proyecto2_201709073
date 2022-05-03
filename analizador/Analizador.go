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

			} else if strings.ToLower(comando[0]) == "mount" {
				self.cmdMount(comando)

			} else if strings.ToLower(comando[0]) == "rep" {
				self.cmdRep(comando)

			} else if strings.ToLower(comando[0]) == "mkfs" {
				self.cmdMkfs(comando)

			} else if strings.ToLower(comando[0]) == "login" {
				self.cmdLogin(comando)

			} else if strings.ToLower(comando[0]) == "logout" {
				self.cmdLogout(comando)

			} else if strings.ToLower(comando[0]) == "mkdir" {
				self.cmdMkdir(comando)

			} else if strings.ToLower(comando[0]) == "mkgrp" {
				self.cmdMkgrp(comando)

			} else if strings.ToLower(comando[0]) == "mkusr" {
				self.cmdMkuser(comando)

			} else if strings.ToLower(comando[0]) == "mkfile" {
				self.cmdMkfile(comando)

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
			comillas[i] = comillas[i] + " "
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

//COMANDO MOUNT
func (self Analizador) cmdMount(comando []string) {
	cmd := comandos.NewMount()

	for i := 1; i < len(comando); i++ {
		if i%2 != 0 && (i+1) < len(comando) {

			if strings.ToLower(comando[i]) == "-path" {
				cmd.Path = comando[i+1]
			} else if strings.ToLower(comando[i]) == "-name" {
				cmd.Name = comando[i+1]
			} else {
				fmt.Println("Error: El comando exec no contiene el comando \"" + comando[i] + "\"")
			}

		}
	}

	cmd.Ejecutar()
}

//COMANDO REP
func (self Analizador) cmdRep(comando []string) {
	cmd := comandos.NewRep()

	for i := 1; i < len(comando); i++ {
		if i%2 != 0 && (i+1) < len(comando) {

			if strings.ToLower(comando[i]) == "-path" {
				cmd.Path = comando[i+1]

			} else if strings.ToLower(comando[i]) == "-name" {
				valor := strings.ToLower(comando[i+1])
				if valor == "file" || valor == "disk" || valor == "tree" {
					cmd.Name = strings.ToLower(comando[i+1])
				} else {
					fmt.Println("Error: Valor del parametro -name incorrecto")
				}

			} else if strings.ToLower(comando[i]) == "-id" {
				cmd.Id = comando[i+1]

			} else if strings.ToLower(comando[i]) == "-ruta" {
				cmd.Ruta = comando[i+1]

			} else {
				fmt.Println("Error: El comando exec no contiene el comando \"" + comando[i] + "\"")
			}
		}
	}

	cmd.Ejecutar()
}

//COMANDO MKFS
func (self Analizador) cmdMkfs(comando []string) {
	cmd := comandos.NewMkfs()

	for i := 1; i < len(comando); i++ {
		if i%2 != 0 && (i+1) < len(comando) {

			if strings.ToLower(comando[i]) == "-id" {
				cmd.Id = comando[i+1]

			} else if strings.ToLower(comando[i]) == "-type" {

				valor := strings.ToLower(comando[i+1])
				if valor == "fast" || valor == "full" {
					cmd.Type = valor
				} else {
					fmt.Println("Error: El comando type solo puede recibir como valor fast y full")
				}

			} else {
				fmt.Println("Error: El comando exec no contiene el comando \"" + comando[i] + "\"")
			}

		}
	}

	cmd.Ejecutar()
}

//COMANDO LOGIN
func (self Analizador) cmdLogin(comando []string) {
	cmd := comandos.NewLogin()

	for i := 1; i < len(comando); i++ {
		if i%2 != 0 && (i+1) < len(comando) {

			if strings.ToLower(comando[i]) == "-id" {
				cmd.Id = comando[i+1]

			} else if strings.ToLower(comando[i]) == "-usuario" {
				cmd.Usuario = comando[i+1]

			} else if strings.ToLower(comando[i]) == "-password" {
				cmd.Password = comando[i+1]

			} else {
				fmt.Println("Error: El comando exec no contiene el comando \"" + comando[i] + "\"")
			}

		}
	}

	cmd.Ejecutar()
}

//COMANOD LOGOUT
func (self Analizador) cmdLogout(comando []string) {
	cmd := comandos.NewLogin()
	cmd.Logout()
}

//COMANDO MKDIR
func (self Analizador) cmdMkdir(comando []string) {
	cmd := comandos.NewMkdir()

	for i := 1; i < len(comando); i++ {
		if strings.ToLower(comando[i]) == "-path" {
			if (i + 1) < len(comando) {
				cmd.Path = comando[i+1]
			}
		} else if strings.ToLower(comando[i]) == "-p" {
			cmd.P = true
		}

	}

	cmd.Ejecutar()
}

//COMANDO MKGRP
func (self Analizador) cmdMkgrp(comando []string) {
	cmd := comandos.NewMkgrp()

	for i := 1; i < len(comando); i++ {
		if i%2 != 0 && (i+1) < len(comando) {
			if strings.ToLower(comando[i]) == "-name" {
				if (i + 1) < len(comando) {
					cmd.Name = comando[i+1]
				}
			} else {
				fmt.Println("Error: El comando exec no contiene el comando \"" + comando[i] + "\"")
			}
		}
	}

	cmd.Ejecutar()
}

//COMANDO MKuSR
func (self Analizador) cmdMkuser(comando []string) {

	cmd := comandos.NewMkuser()

	for i := 1; i < len(comando); i++ {
		if i%2 != 0 && (i+1) < len(comando) {
			if strings.ToLower(comando[i]) == "-usuario" {
				if (i + 1) < len(comando) {
					cmd.Usuario = comando[i+1]
				}
			} else if strings.ToLower(comando[i]) == "-pwd" {
				if (i + 1) < len(comando) {
					cmd.Password = comando[i+1]
				}
			} else if strings.ToLower(comando[i]) == "-grp" {
				if (i + 1) < len(comando) {
					cmd.Group = comando[i+1]
				}
			} else {
				fmt.Println("Error: El comando exec no contiene el comando \"" + comando[i] + "\"")
			}
		}
	}

	cmd.Ejecutar()
}

//COMANDO MKFILE OJALA FUNCIONE
func (self Analizador) cmdMkfile(comando []string) {
	cmd := comandos.NewMkfile()

	for i := 1; i < len(comando); i++ {
		if strings.ToLower(comando[i]) == "-path" {
			if (i + 1) < len(comando) {
				cmd.Path = comando[i+1]
			}
		} else if strings.ToLower(comando[i]) == "-r" {
			cmd.R = true

		} else if strings.ToLower(comando[i]) == "-size" {
			if (i + 1) < len(comando) {

				valor, err := strconv.Atoi(comando[i+1])
				if err != nil {
					cmd.Size = -1
				} else {
					cmd.Size = valor
				}
			}
		} else if strings.ToLower(comando[i]) == "-cont" {
			if (i + 1) < len(comando) {
				cmd.Cont = comando[i+1]
			}
		}

	}

	cmd.Ejecutar()
}
