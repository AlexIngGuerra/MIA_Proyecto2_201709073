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
func EscribirInodo(archivo *os.File, superBloque SuperBloque, inodo Inodo, n int32) SuperBloque {
	if superBloque.Free_inodes_count < 1 {
		fmt.Println("Error: No se pueden crear más inodos")
		return superBloque
	}

	archivo.Seek(superBloque.First_ino, 0)
	buffer := bytes.Buffer{}
	binary.Write(&buffer, binary.BigEndian, &inodo)
	EscribirArchivo(archivo, buffer.Bytes())

	superBloque.Free_inodes_count -= 1                           // -1 Inodos libres
	superBloque.First_ino += int64(superBloque.Inode_Size)       //+1 Inodo al inicio
	MarcarPrimerBitLibre(archivo, superBloque.Bm_inode_start, n) //marcamos el primer bit como libre

	return superBloque
}

func ObtenerBloquesArchivo(archivo *os.File, superBloque SuperBloque, inodo Inodo) []BloqueArchivo {
	var bloques []BloqueArchivo

	for i := 0; i < len(inodo.Block); i++ {
		if inodo.Block[i] != -1 {
			bloc := LeerBloqueA(archivo, superBloque.Block_start+int64(inodo.Block[i])*int64(superBloque.Block_Size))
			bloques = append(bloques, bloc)
		}
	}

	return bloques
}
