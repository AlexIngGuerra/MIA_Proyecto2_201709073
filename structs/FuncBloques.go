package structs

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"os"
	"strconv"
	"strings"
	"unsafe"
)

//FUNCION PARA LEER BLOQUES DE TIPO ARCHIVO
func LeerBloqueA(archivo *os.File, inicio int64) BloqueArchivo {
	bloque := BloqueArchivo{}

	archivo.Seek(inicio, 0)

	size := int(unsafe.Sizeof(bloque))
	data := LeerArchivo(archivo, size)
	buffer := bytes.NewBuffer(data)
	binary.Read(buffer, binary.BigEndian, &bloque)

	return bloque
}

//FUNCION PARA LEER BLOQUES DE TIPO CARPETA
func LeerBloqueC(archivo *os.File, inicio int64) BloqueCarpeta {
	bloque := BloqueCarpeta{}

	archivo.Seek(inicio, 0)

	size := int(unsafe.Sizeof(bloque))
	data := LeerArchivo(archivo, size)
	buffer := bytes.NewBuffer(data)
	binary.Read(buffer, binary.BigEndian, &bloque)

	return bloque
}

//FUNCION PARA ESCRIBIR BLOQUES DE ARCHIVO
func EscribirBloqueA(archivo *os.File, superBloque SuperBloque, bloque BloqueArchivo) SuperBloque {

	if superBloque.Free_inodes_count < 1 {
		fmt.Println("Error: No se pueden crear más bloques")
		return superBloque
	}

	archivo.Seek(superBloque.First_bloc, 0)
	buffer := bytes.Buffer{}
	binary.Write(&buffer, binary.BigEndian, &bloque)
	EscribirArchivo(archivo, buffer.Bytes())

	superBloque.Free_blocks_count -= 1                                                  // -1 Inodos libres
	superBloque.First_bloc += int64(superBloque.Block_Size)                             //+1 Inodo al inicio
	MarcarPrimerBitLibre(archivo, superBloque.Bm_block_start, superBloque.Blocks_count) //marcamos el primer bit como libre

	return superBloque
}

//FUNCION PARA ESCRIBIR BLOQUES DE CARPETA
func EscribirBloqueC(archivo *os.File, superBloque SuperBloque, bloque BloqueCarpeta) SuperBloque {

	if superBloque.Free_inodes_count < 1 {
		fmt.Println("Error: No se pueden crear más bloques")
		return superBloque
	}

	archivo.Seek(superBloque.First_bloc, 0)
	buffer := bytes.Buffer{}
	binary.Write(&buffer, binary.BigEndian, &bloque)
	EscribirArchivo(archivo, buffer.Bytes())

	superBloque.Free_blocks_count -= 1                                                  // -1 Inodos libres
	superBloque.First_bloc += int64(superBloque.Block_Size)                             //+1 Inodo al inicio
	MarcarPrimerBitLibre(archivo, superBloque.Bm_block_start, superBloque.Blocks_count) //marcamos el primer bit como libres

	return superBloque
}

//OBTENEMOS UN STRING DE UN ARREGLO DE BLOQUES
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

//ESCRIBIMOS UN STRING EN UN ARREGLO DE BLOQUES
func EscribirTextoEnBloques(texto string) []BloqueArchivo {
	var bloques []BloqueArchivo
	caracteres := []uint8(texto)
	numCaracteres := len(caracteres)
	numBloques := 1
	for numCaracteres > 64 {
		numBloques++
		numCaracteres -= 64
	}
	var char uint8

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

//TRANSFORMAMOS UN STRING UN ARREGLO DE GRUPOS
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

//AGREGAMOS UN GRUPO A UN ARREGLO
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

//ARREGLAMOS UN USUARIO A UN ARREGLO
func AddUsuario(grupos []Grupo, usuario Usuario, grupo string) []Grupo {

	for i := 0; i < len(grupos); i++ {
		for j := 0; j < len(grupos[i].Users); j++ {

			if grupos[i].Users[j].User == usuario.User {
				fmt.Println("Error: No puede existur un usuario repetido")
				return grupos
			}

		}
	}

	for i := 0; i < len(grupos); i++ {
		if grupos[i].Name == grupo {

			usuario.Id = len(grupos[i].Users) + 1
			fmt.Println(usuario.Id)
			grupos[i].Users = append(grupos[i].Users, usuario)

		}
	}

	return grupos
}

//PASAMOS UN ARREGLO DE GRUPOS A UN STRING
func GroupsToString(grupos []Grupo) string {
	cadena := ""

	for i := 0; i < len(grupos); i++ {
		cadena = cadena + strconv.Itoa(grupos[i].Id) + ", G, " + grupos[i].Name + "\n"

		for j := 0; j < len(grupos[i].Users); j++ {
			convirtiendo := strconv.Itoa(grupos[i].Users[j].Id)
			cadena += convirtiendo + ", U, " + grupos[i].Name + ", " + grupos[i].Users[j].User + ", " + grupos[i].Users[j].Password + "\n"
		}

	}

	return cadena
}

//OBTENEMOS EL NOMBRE DE UN BLOQUE
func GetNameBloque(Name string) [12]uint8 {
	var retorno [12]uint8
	arreglo := []uint8(Name)
	for i := 0; i < 12; i++ {
		if i < len(arreglo) {
			retorno[i] = arreglo[i]
		} else {
			retorno[i] = 0
		}
	}

	return retorno
}

//OBTENEMOS EL NOMBRE DE UN BLOQUE COMO STRING
func GetNameBloqueString(Name [12]uint8) string {
	cadena := ""

	for i := 0; i < len(Name); i++ {
		if Name[i] == 0 {
			break

		}
		cadena = cadena + string(Name[i])
	}

	return cadena
}

//AGREGAMOS UN NUEVO VALOR AL BLOQUE CARPETA
func AddContenidoBloqueCarpeta(Name string, Inodo int32, bloque BloqueCarpeta) BloqueCarpeta {
	for i := 0; i < 4; i++ {
		if bloque.Contenido[i].Apuntador == -1 {
			bloque.Contenido[i] = Contenido{Name: GetNameBloque(Name), Apuntador: Inodo}
			break
		}
	}
	return bloque
}
