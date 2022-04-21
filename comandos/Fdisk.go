package comandos

import (
	"MIA/structs"
	"bytes"
	"encoding/binary"
	"fmt"
	"log"
	"os"
	"unsafe"
)

type Fdisk struct {
	Size int
	Unit string
	Path string
	Type string
	Fit  string
	Name string
}

func NewFdisk() Fdisk {
	return Fdisk{Size: 0, Unit: "K", Path: "", Type: "P", Fit: "WF", Name: ""}
}

func (self Fdisk) Ejecutar() {
	if self.tieneErrores() {
		return
	}
	var mbr structs.Mbr
	mbr = structs.GetMbr(self.Path)

	if mbr.Tamano <= 0 {
		fmt.Println("Error: El Mbr no es funcional")
		return
	}

	if self.Type == "P" || self.Type == "E" {
		self.crearParticionPrimaria(mbr)
	} else if self.Type == "L" {
		self.crearParticionLogica(mbr)
	}

}

func (self Fdisk) tieneErrores() bool {
	errores := false
	return errores
}

func (self Fdisk) crearParticionPrimaria(mbr structs.Mbr) {

	if self.tieneErroresParticionPrimExt(mbr) {
		return
	}

	var particion structs.Partition

	particion.Size = structs.GetSize(self.Size, self.Unit)
	particion.Fit = structs.GetFit(self.Fit)
	particion.Name = structs.GetName(self.Name)
	particion.Type = structs.GetType(self.Type)
	particion.Status = 0
	particion.Start = self.getStartPrimaria(mbr)

	for i := 0; i < 4; i++ {
		if mbr.Particion[i].Size == 0 {
			mbr.Particion[i] = particion
			break
		}
	}

	archivo, err := os.OpenFile(self.Path, os.O_RDWR, 0777)
	defer archivo.Close()
	if err != nil {
		log.Fatal(err)
	}

	var bufferMbr bytes.Buffer
	binary.Write(&bufferMbr, binary.BigEndian, &mbr)
	structs.EscribirArchivo(archivo, bufferMbr.Bytes())

}

func (self Fdisk) tieneErroresParticionPrimExt(mbr structs.Mbr) bool {
	errores := false
	if int(structs.GetSize(self.Size, self.Unit)) > structs.GetEspacioLibreMbr(mbr) {
		errores = true
		fmt.Println("Error: No hay suficiente espacio en el disco para crear la particion")
	}

	hayEspacio := false
	for i := 0; i < 4; i++ {
		if mbr.Particion[i].Size == 0 {
			hayEspacio = true
		}
	}

	if !hayEspacio {
		errores = true
		fmt.Println("Error: No se pueden crear más particiones en el disco")
	}

	return errores
}

func (self Fdisk) getStartPrimaria(mbr structs.Mbr) int64 {
	start := int64(unsafe.Sizeof(mbr))
	for i := 0; i < 4; i++ {
		start += mbr.Particion[i].Size
	}
	return start
}

func (self Fdisk) crearParticionLogica(mbr structs.Mbr) {

}
