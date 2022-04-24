package structs

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"log"
	"os"
	"unsafe"
)

func GetFit(Fit string) uint8 {
	if Fit == "FF" {
		return 'F'
	} else if Fit == "BF" {
		return 'B'
	} else if Fit == "WF" {
		return 'W'
	}
	return 0
}

func GetSize(Size int, Unit string) int64 {
	if Unit == "M" {
		return int64(Size * 1024 * 1024)
	} else if Unit == "K" {
		return int64(Size * 1024)
	} else if Unit == "B" {
		return int64(Size)
	}
	return 0
}

func GetFecha(Fecha string) [20]uint8 {
	arreglo := []uint8(Fecha)
	var retorno [20]uint8
	for i := 0; i < 20; i++ {
		if i < len(arreglo) {
			retorno[i] = arreglo[i]
		} else {
			retorno[i] = 0
		}
	}
	return retorno
}

func GetName(Name string) [20]uint8 {
	var retorno [20]uint8
	arreglo := []uint8(Name)
	for i := 0; i < 20; i++ {
		if i < len(arreglo) {
			retorno[i] = arreglo[i]
		} else {
			retorno[i] = 0
		}
	}

	return retorno
}

func GetType(Type string) uint8 {
	if Type == "P" {
		return 'P'
	} else if Type == "E" {
		return 'E'
	} else if Type == "L" {
		return 'L'
	}
	return 0
}

func UintToString(arreglo [20]uint8) string {
	cadena := ""
	for i := 0; i < 20; i++ {
		cadena += string(arreglo[i])
	}
	return cadena
}

//########## ARCHIVOS  #############################

func LeerArchivo(archivo *os.File, numero int) []byte {
	bytes := make([]byte, numero)

	_, err := archivo.Read(bytes)
	if err != nil {
		fmt.Println("Error: No se ha podido leer en el disco")
		log.Fatal(err)
	}

	return bytes
}

func EscribirArchivo(archivo *os.File, bytes []byte) bool {
	_, err := archivo.Write(bytes)

	if err != nil {
		fmt.Println("Error: No se ha podido escribir en el disco")
		log.Fatal(err)
		return false
	}
	return true
}

//############ DISCOS #############################

func GetMbr(Path string) Mbr {
	archivo, err := os.Open(Path)
	defer archivo.Close()

	if err != nil {
		fmt.Println("Error: No se pudo abrir el archivo")
		log.Fatal(err)
	}

	archivo.Seek(0, 0)
	var mbr Mbr

	size := int(unsafe.Sizeof(mbr))
	data := LeerArchivo(archivo, size)
	buffer := bytes.NewBuffer(data)

	err = binary.Read(buffer, binary.BigEndian, &mbr)
	if err != nil {
		fmt.Println("Error: No se ha podido leer el mbr del disco")
		log.Fatal(err)
	}

	return mbr
}

func GetEbr(archivo *os.File, Posicion int64) Ebr {
	archivo.Seek(Posicion, 0)
	var ebr Ebr

	size := int(unsafe.Sizeof(ebr))
	data := LeerArchivo(archivo, size)
	buffer := bytes.NewBuffer(data)
	binary.Read(buffer, binary.BigEndian, &ebr)

	return ebr
}

func GetEspacioLibreMbr(mbr Mbr) int {
	valor := int(mbr.Tamano)

	for i := 0; i < 4; i++ {
		valor -= int(mbr.Particion[i].Size)
	}

	return valor
}
