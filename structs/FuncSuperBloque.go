package structs

import (
	"bytes"
	"encoding/binary"
	"os"
	"unsafe"
)

//FUNCION PARA LEER EL SUPERBLOQUE
func LeerSuperBloque(archivo *os.File, inicio int64) SuperBloque {
	superBloque := SuperBloque{}

	archivo.Seek(inicio, 0)

	size := int(unsafe.Sizeof(superBloque))
	data := LeerArchivo(archivo, size)
	buffer := bytes.NewBuffer(data)
	binary.Read(buffer, binary.BigEndian, &superBloque)

	return superBloque
}

//FUNCION PARA ESCRIBIR EL SUPER BLOQUE EN EL ARCHIVO
func EscribirSuperBloque(archivo *os.File, superBloque SuperBloque, inicio int64) {

	archivo.Seek(inicio, 0)
	buffer := bytes.Buffer{}
	binary.Write(&buffer, binary.BigEndian, &superBloque)
	EscribirArchivo(archivo, buffer.Bytes())

}
