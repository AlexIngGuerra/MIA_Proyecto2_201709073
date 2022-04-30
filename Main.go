package main

import (
	"MIA/analizador"
	"MIA/structs"
	"fmt"
	"strconv"
	"unsafe"
)

func main() {
	fmt.Println("WALTER ALEXANDER GUERRA DUQUE 201709073")

	initPrueba()

	fmt.Println("----- DATOS IMPORTANTES -----")

	fmt.Println("MBR size: " + strconv.Itoa(int(unsafe.Sizeof(structs.Mbr{}))))
	fmt.Println("EBR size: " + strconv.Itoa(int(unsafe.Sizeof(structs.Ebr{}))))

	fmt.Println("SuperBloque size: " + strconv.Itoa(int(unsafe.Sizeof(structs.SuperBloque{}))))
	fmt.Println("Inodo size: " + strconv.Itoa(int(unsafe.Sizeof(structs.Inodo{}))))
	fmt.Println("BloqueA size: " + strconv.Itoa(int(unsafe.Sizeof(structs.BloqueArchivo{}))))
	fmt.Println("BloqueC size: " + strconv.Itoa(int(unsafe.Sizeof(structs.BloqueCarpeta{}))))

	fmt.Print("-----------------------------\n\n")

}

func initPrueba() {
	analizador := analizador.NewAnalizador()
	var comando string
	comando = "exec    -path=\"C:/Users/WALTER/Dropbox/Mi PC (DESKTOP-DUSC6PO)/Documents/1/Entradas/entrada2.script\"  #Comentario :D"
	analizador.Analizar(comando)
}

func initProyecto() {
	analizador := analizador.NewAnalizador()
	var comando string
	for true {
		fmt.Scanln(&comando)
		analizador.Analizar(comando)
		if comando == "salir" {
			break
		}
	}
}
