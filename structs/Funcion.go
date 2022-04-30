package structs

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"log"
	"math"
	"os"
	"unsafe"
)

//OBTENER EL FIT
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

//OBTENER EL TAMAÑO EN BASE A UNA UNIDAD DE MEDIDA
func GetSize(Size int, Unit string) int32 {
	if Unit == "M" {
		return int32(Size * 1024 * 1024)
	} else if Unit == "K" {
		return int32(Size * 1024)
	} else if Unit == "B" {
		return int32(Size)
	}
	return 0
}

//OBTENER LA FECHA COMO ARREGLO
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

//OBTENER EL NOMBRE COMO ARREGLO UINT
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

//OBTENER EL TIPO COMO UINT
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

//PASAR DE UINT A STRING
func UintToString(arreglo [20]uint8) string {
	cadena := ""
	for i := 0; i < 20; i++ {
		cadena += string(arreglo[i])
	}
	return cadena
}

//########## ARCHIVOS  #############################

//LEER EL ARCHIVO EN UN TAMAÑO ESPECIFICO
func LeerArchivo(archivo *os.File, numero int) []byte {
	bytes := make([]byte, numero)

	_, err := archivo.Read(bytes)
	if err != nil {
		fmt.Println("Error: No se ha podido leer en el disco")
		log.Fatal(err)
	}

	return bytes
}

//ESCRIBIR EN EL ARCHIVO
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

//OBTENER EL MBR DE UN ARCHIVO
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

//OBTENER EL EBR DE UN ARCHIVO
func GetEbr(archivo *os.File, Posicion int64) Ebr {
	archivo.Seek(Posicion, 0)
	var ebr Ebr

	size := int(unsafe.Sizeof(ebr))
	data := LeerArchivo(archivo, size)
	buffer := bytes.NewBuffer(data)
	binary.Read(buffer, binary.BigEndian, &ebr)

	return ebr
}

//OBTENER EL ESPACIO LIBRE DE UN MBR
func GetEspacioLibreMbr(mbr Mbr) int {
	valor := int(mbr.Size)

	for i := 0; i < 4; i++ {
		valor -= int(mbr.Particion[i].Size)
	}

	return valor
}

//############ ARCHIVOS #############################

//OBTENER EL VALOR N PARA EL MANEJO DE ARCHIVOSS
func GetN(Size int32) int32 {
	SB := float64(unsafe.Sizeof(SuperBloque{}))
	IN := float64(unsafe.Sizeof(Inodo{}))
	T := float64(Size - 1)
	N := (T - SB) / (196 + IN)
	result := math.Floor(N)
	return int32(result)
}

//OBTENER LA INFORMACION DE LA PARTICION
func GetParticion(Name [20]uint8, mbr Mbr, archivo *os.File) InfoPart {
	var info InfoPart

	for i := 0; i < 4; i++ {

		part := mbr.Particion[i]
		if Name == part.Name && part.Type != 'E' {
			info.Size = part.Size
			info.Start = part.Start
		}

		if part.Type == 'E' {
			apuntador := mbr.Particion[i].Start
			ebrActual := GetEbr(archivo, apuntador)

			for ebrActual.Next != 0 {
				if ebrActual.Name == Name {
					info.Size = ebrActual.Size
					info.Start = ebrActual.Start
				}
				apuntador = ebrActual.Next
				ebrActual = GetEbr(archivo, apuntador)
			}

		}
	}

	return info
}
