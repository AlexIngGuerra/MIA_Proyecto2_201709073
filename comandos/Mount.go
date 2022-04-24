package comandos

import (
	"MIA/structs"
	"fmt"
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

	var nodo ParticionMontada
	nodo.Path = self.Path
	nodo.Name = structs.GetName(self.Name)

	mbr := structs.GetMbr(self.Path)
	if mbr.Tamano <= 0 {
		fmt.Println("Error: No se puede utilizar el disco solicitado")
		return
	}

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

func (self Mount) generarId() {
	id := "73"
	numero := int(0)   // Disco
	letra := uint8(65) // Particion
	for i := 0; i < len(Montados); i++ {
		if self.Path == Montados[i].Path {
			numero = Montados[i].Numero
			letra++
		}
	}

	id = id + strconv.Itoa(numero) + string(letra)
	fmt.Print(id)
}
