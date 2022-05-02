package structs

import (
	"bytes"
	"encoding/binary"
	"os"
)

//BUSCA EL NUMERO QUE CORRESPONDE AL PRIMER BIT LIBRE
func BuscarBitLibre(archivo *os.File, inicio int64, size int32) int32 {
	valor := int32(-1)

	archivo.Seek(inicio, 0)
	for i := int64(0); i < int64(size); i++ {
		caracter := LeerArchivo(archivo, 1)
		ascii := uint8(caracter[0])
		valor++
		if ascii == 0 {
			break
		}
		inicio++
	}
	return valor
}

//MARCA EL PRIMER BIT LIBRE COMO 1
func MarcarPrimerBitLibre(archivo *os.File, inicio int64, size int32) {

	archivo.Seek(inicio, 0)
	for i := int64(0); i < int64(size); i++ {
		caracter := LeerArchivo(archivo, 1)
		ascii := uint8(caracter[0])
		if ascii == 0 {
			break
		}
		inicio++
	}

	uno := uint8(1)
	archivo.Seek(inicio, 0)
	buffer := bytes.Buffer{}
	binary.Write(&buffer, binary.BigEndian, &uno)
	EscribirArchivo(archivo, buffer.Bytes())
}

//VERIFICA SI EXISTE EL STRUCT EN EL INODO
func ExisteStructEnBM(archivo *os.File, inicio int64, size int32, numStruct int64) bool {
	archivo.Seek(inicio+numStruct, 0)
	caracter := LeerArchivo(archivo, 1)
	ascii := uint8(caracter[0])
	if ascii == 1 {
		return true
	}
	return false
}
