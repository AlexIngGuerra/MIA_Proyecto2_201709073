package comandos

import (
	"MIA/structs"
	"fmt"
	"log"
	"os"
	"strconv"
)

var Montados []ParticionMontada

type ParticionMontada struct {
	Id     string
	Path   string    //Disco
	Name   [20]uint8 //Particion
	Numero int
	Disco  uint8
}

type Mount struct {
	Path string
	Name string
}

func NewMount() Mount {
	return Mount{Path: "", Name: ""}
}

func (self Mount) Ejecutar() {
	if self.tieneErrores() {
		return
	}

	mbr := structs.GetMbr(self.Path)
	if mbr.Tamano <= 0 {
		fmt.Println("Error: No se puede utilizar el disco solicitado")
		return
	}

	if !self.existeParticion(mbr) {
		fmt.Println("Error: La partición no existe")
		return
	}

	if self.estaMontada(structs.GetName(self.Name)) {
		fmt.Println("Error: La partición ya ha sido montada")
		return
	}

	nodo := self.generarNodo()
	Montados = append(Montados, nodo)
	fmt.Println("Montada la particion con el id: " + nodo.Id)
}

func (self Mount) tieneErrores() bool {
	errores := false
	if self.Path == "" {
		errores = true
		fmt.Println("Error: El parametro -path es obligatorio")
	}

	if self.Name == "" {
		errores = true
		fmt.Println("Error: El parametro -name es obligatorio")
	}

	return errores
}

func (self Mount) existeParticion(mbr structs.Mbr) bool {
	archivo, err := os.OpenFile(self.Path, os.O_RDWR, 0777)
	defer archivo.Close()

	if err != nil {
		fmt.Println("Error: No se ha podido abrir el archivo")
		log.Fatal(err)
	}

	nombre := structs.GetName(self.Name)

	for i := 0; i < 4; i++ {
		if mbr.Particion[i].Name == nombre {
			return true
		}

		if mbr.Particion[i].Type == 'E' {
			apuntador := mbr.Particion[i].Start
			ebrActual := structs.GetEbr(archivo, apuntador)

			if ebrActual.Name == nombre {
				return true
			}

			for ebrActual.Next != 0 {
				apuntador = ebrActual.Next
				ebrActual = structs.GetEbr(archivo, apuntador)

				if ebrActual.Name == nombre {
					return true
				}
			}
		}
	}

	return false
}

func (self Mount) estaMontada(Name [20]uint8) bool {
	for i := 0; i < len(Montados); i++ {
		if Montados[i].Name == Name {
			return true
		}
	}
	return false
}

func (self Mount) generarNodo() ParticionMontada {
	var nuevo ParticionMontada

	id := "73"
	numero := int(1)   // Disco
	letra := uint8(65) // Particion
	for i := 0; i < len(Montados); i++ {
		if self.Path == Montados[i].Path {
			numero = Montados[i].Numero
			letra++
		}
	}

	id = id + strconv.Itoa(numero) + string(letra)

	nuevo.Id = id
	nuevo.Numero = numero
	nuevo.Disco = letra
	nuevo.Path = self.Path
	nuevo.Name = structs.GetName(self.Name)

	return nuevo
}
