package structs

import (
	"fmt"
	"strconv"
	"strings"
)

func LeerBloquesArchivo(bloques []BloqueArchivo) string {
	contenido := ""
	for i := 0; i < len(bloques); i++ {
		bloque := bloques[i]

		for j := 0; j < 64; j++ {
			if bloque.Contenido[j] == 0 {
				break
			}
			contenido = contenido + string(bloque.Contenido[j])
		}

	}

	return contenido
}

func EscribirBloqueArchivo(texto string) []BloqueArchivo {
	var bloques []BloqueArchivo
	caracteres := []uint8(texto)
	numCaracteres := len(caracteres)
	numBloques := 1
	for numCaracteres > 64 {
		numBloques++
		numCaracteres -= 64
	}
	var char uint8
	fmt.Println(caracteres)

	for i := 0; i < numBloques; i++ {
		var bloque BloqueArchivo
		nulo := uint8(0)
		for j := 0; j < 64; j++ {
			if len(caracteres) == 0 {
				nulo = nulo % 10
				bloque.Contenido[j] = nulo
				nulo++
			} else {
				char, caracteres = caracteres[0], caracteres[1:]
				bloque.Contenido[j] = char
			}

		}
		bloques = append(bloques, bloque)
		if len(caracteres) == 0 {
			break
		}

	}

	return bloques
}

func GetGruposYUsuarios(texto string) []Grupo {
	var grupos []Grupo
	linea := strings.Split(texto, "\n")
	for i := 0; i < len(linea); i++ {
		if len(linea[i]) < 1 {
			continue
		}

		parametros := strings.Split(linea[i], ",")

		tipo := strings.Trim(parametros[1], " ")
		if tipo == "G" {
			var grupo Grupo
			grupo.Id, _ = strconv.Atoi(strings.Trim(parametros[0], " "))
			grupo.Name = strings.Trim(parametros[2], " ")
			grupos = append(grupos, grupo)

		} else if tipo == "U" {
			var usuario Usuario
			usuario.Id, _ = strconv.Atoi(strings.Trim(parametros[0], " "))
			usuario.User = strings.Trim(parametros[3], " ")
			usuario.Password = strings.Trim(parametros[4], " ")

			for gr := 0; gr < len(grupos); gr++ {
				if grupos[gr].Name == strings.Trim(parametros[2], " ") {
					grupos[gr].Users = append(grupos[gr].Users, usuario)
				}
			}

		}

	}

	return grupos
}

func AddGrupo(grupos []Grupo, grupo Grupo) []Grupo {
	for i := 0; i < len(grupos); i++ {
		if grupos[i].Id == grupo.Id || grupos[i].Name == grupo.Name {
			fmt.Println("Error: No se puede agregar un grupo repetido")
			return grupos
		}
	}
	grupos = append(grupos, grupo)
	return grupos
}

func AddUsuario(grupos []Grupo, usuario Usuario, grupo string) []Grupo {

	for i := 0; i < len(grupos); i++ {
		for j := 0; j < len(grupos[i].Users); j++ {

			if grupos[i].Users[j].User == usuario.User {
				fmt.Println("Error: No puede exisitr un usuario repetido")
				return grupos
			}

		}
	}

	for i := 0; i < len(grupos); i++ {
		if grupos[i].Name == grupo {
			grupos[i].Users = append(grupos[i].Users, usuario)
			break
		}
	}

	return grupos
}

func GroupsToString(grupos []Grupo) string {
	cadena := ""

	for i := 0; i < len(grupos); i++ {
		cadena = cadena + strconv.Itoa(grupos[i].Id) + ", G, " + grupos[i].Name + "\n"

		for j := 0; j < len(grupos[i].Users); j++ {
			cadena = cadena + strconv.Itoa(grupos[i].Users[j].Id) + ", U, " + grupos[i].Name + ", " + grupos[i].Users[j].User + ", " + grupos[i].Users[j].Password + "\n"
		}

	}

	return cadena
}
