package structs

import (
	"fmt"
	"log"
	"os"
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

//########## ARCHIVOS ####################

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
