package structs

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"os"
	"unsafe"
)

//FUNCION PARA LEER UN INODO
func LeerInodo(archivo *os.File, inicio int64) Inodo {
	inodo := Inodo{}

	archivo.Seek(inicio, 0)

	size := int(unsafe.Sizeof(inodo))
	data := LeerArchivo(archivo, size)
	buffer := bytes.NewBuffer(data)
	binary.Read(buffer, binary.BigEndian, &inodo)

	return inodo
}

//FUNCION PARA ESCRIBIR EL INODO EN EL ARCHIVOS
func EscribirInodo(archivo *os.File, superBloque SuperBloque, inodo Inodo) SuperBloque {
	if superBloque.Free_inodes_count < 1 {
		fmt.Println("Error: No se pueden crear más inodos")
		return superBloque
	}

	archivo.Seek(superBloque.First_ino, 0)
	buffer := bytes.Buffer{}
	binary.Write(&buffer, binary.BigEndian, &inodo)
	EscribirArchivo(archivo, buffer.Bytes())

	superBloque.Free_inodes_count -= 1                                                  // -1 Inodos libres
	superBloque.First_ino += int64(superBloque.Inode_Size)                              //+1 Inodo al inicio
	MarcarPrimerBitLibre(archivo, superBloque.Bm_inode_start, superBloque.Inodes_count) //marcamos el primer bit como libre

	return superBloque
}

//OBTIENE TODO LOS BLOQUES DE ARCHIVO EN UN INODO
func ObtenerBloquesArchivo(archivo *os.File, superBloque SuperBloque, inodo Inodo) []BloqueArchivo {
	var bloques []BloqueArchivo
	if inodo.Type != 1 {
		return bloques
	}

	for i := 0; i < len(inodo.Block); i++ {
		if inodo.Block[i] != -1 {
			bloc := LeerBloqueA(archivo, superBloque.Block_start+int64(inodo.Block[i])*int64(superBloque.Block_Size))
			bloques = append(bloques, bloc)
		}
	}

	return bloques
}

//PERMITE BUSCAR EL INODO SIGUIENTE SI EXISTE
func GetSiguienteInodo(archivo *os.File, superbloque SuperBloque, inoActual Inodo, nombre string) (Inodo, int64, int32) {
	if inoActual.Type == 1 {
		fmt.Println("Error: Se ha ingresado un inodo tipo archivo")
		return Inodo{}, -1, 0
	}

	for ap := 0; ap < 16; ap++ {
		if inoActual.Block[ap] != -1 {
			bloque := LeerBloqueC(archivo, superbloque.Block_start+int64(superbloque.Block_Size)*int64(inoActual.Block[ap]))
			name := GetNameBloqueString(GetNameBloque(nombre))
			coincide := false
			for c := 0; c < 4; c++ {

				for i := 0; i < len(name); i++ {
					coincide = name[i] == bloque.Contenido[c].Name[i]
				}

				if coincide {
					apuntador := superbloque.Inode_start + int64(superbloque.Inode_Size)*int64(bloque.Contenido[c].Apuntador)
					inode := LeerInodo(archivo, apuntador)
					return inode, apuntador, bloque.Contenido[c].Apuntador
				}

			}

		}
	}

	return Inodo{}, -1, 0
}

//VERIFICA SI EL NOMBRE BUSCADO EXISTE EN LA CARPETA ACTUAL
func ExisteNombreInodo(archivo *os.File, superbloque SuperBloque, inodoActual Inodo, nombre string) bool {

	for ap := 0; ap < 16; ap++ {

		if inodoActual.Block[ap] != -1 {
			bloque := LeerBloqueC(archivo, superbloque.Block_start+int64(superbloque.Block_Size)*int64(inodoActual.Block[ap]))
			for i := 0; i < 4; i++ {

				if bloque.Contenido[i].Name == GetNameBloque(nombre) {
					return true
				}

			}

		}

	}

	return false
}

//AGREAGA EL NUEVO APUNTADOR A LOS BLOQUES DEL INODO
func AgregarNuevoApuntador(archivo *os.File, superbloque SuperBloque, inoActual Inodo, numInoActual int32, inoName string, sigInodo int32, inoAnterior int32) (Inodo, SuperBloque) {

	bloqueLibre := -1
	apCrear := -1

	for ap := 0; ap < 16; ap++ {
		apuntador := superbloque.Block_start + int64(superbloque.Block_Size)*int64(inoActual.Block[ap])
		if inoActual.Block[ap] != -1 {
			bloque := LeerBloqueC(archivo, apuntador)

			for c := 0; c < 4; c++ {
				if bloque.Contenido[c].Apuntador == -1 {
					bloqueLibre = ap
				}
			}

			if bloqueLibre != -1 {
				break
			}

		} else {
			apCrear = ap
			break
		}
	}

	if bloqueLibre != -1 {

		apuntador := superbloque.Block_start + int64(superbloque.Block_Size)*int64(inoActual.Block[bloqueLibre])
		bloque := LeerBloqueC(archivo, apuntador)

		for c := 0; c < 4; c++ {
			if bloque.Contenido[c].Apuntador == -1 {

				bloque.Contenido[c] = Contenido{Name: GetNameBloque(inoName), Apuntador: sigInodo}

				archivo.Seek(apuntador, 0)
				buffer := bytes.Buffer{}
				binary.Write(&buffer, binary.BigEndian, &bloque)
				EscribirArchivo(archivo, buffer.Bytes())
				break

			}
		}

	} else {
		num := BuscarBitLibre(archivo, superbloque.Bm_block_start, superbloque.Blocks_count)
		if apCrear != 0 {
			bloque := NewBloqueCarpeta()
			bloque.Contenido[0] = Contenido{Name: GetNameBloque(inoName), Apuntador: sigInodo}

			superbloque = EscribirBloqueC(archivo, superbloque, bloque)
			inoActual.Block[apCrear] = num

		} else {
			bloque := NewPrimerBloqueCarpeta(numInoActual, inoAnterior)
			bloque.Contenido[2] = Contenido{Name: GetNameBloque(inoName), Apuntador: sigInodo}

			superbloque = EscribirBloqueC(archivo, superbloque, bloque)
			inoActual.Block[apCrear] = num
		}

	}

	return inoActual, superbloque
}
