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
	//Buscamos el Mbr
	var mbr structs.Mbr
	mbr = structs.GetMbr(self.Path)

	if mbr.Tamano <= 0 {
		fmt.Println("Error: El Mbr no es funcional")
		return
	}

	//Pasamos al metodo para crear cada particion
	if self.Type == "P" || self.Type == "E" {
		self.crearParticionPrimaria(mbr)
	} else if self.Type == "L" {
		self.crearParticionLogica(mbr)
	}
	fmt.Print("\n")
}

func (self Fdisk) tieneErrores() bool {
	errores := false
	return errores
}

func (self Fdisk) crearParticionPrimaria(mbr structs.Mbr) {

	archivo, err := os.OpenFile(self.Path, os.O_RDWR, 0777)
	defer archivo.Close()

	if err != nil {
		fmt.Println("Error: No se ha podido abrir el archivo")
		log.Fatal(err)
	}

	//Creamos la particion
	var particion structs.Partition
	particion.Size = structs.GetSize(self.Size, self.Unit)
	particion.Fit = structs.GetFit(self.Fit)
	particion.Name = structs.GetName(self.Name)
	particion.Type = structs.GetType(self.Type)
	particion.Status = '0'
	particion.Start = self.getStartPrimaria(mbr)

	for i := 0; i < 4; i++ {
		if mbr.Particion[i].Size == 0 {
			mbr.Particion[i] = particion
			break
		}
	}

	archivo.Seek(0, 0)
	var buffer bytes.Buffer
	binary.Write(&buffer, binary.BigEndian, &mbr)
	structs.EscribirArchivo(archivo, buffer.Bytes())

	for i := 0; i < 4; i++ {
		fmt.Println("###############################")
		fmt.Println(structs.UintToString(mbr.Particion[i].Name))
		fmt.Println(mbr.Particion[i].Size)
		fmt.Println(mbr.Particion[i].Start)
		fmt.Println(string(mbr.Particion[i].Fit))
		fmt.Println(string(mbr.Particion[i].Type))
		fmt.Println("###############################")
	}

	if particion.Type == 'E' {
		var ebrNull structs.Ebr
		ebrNull.Name = structs.GetName("EBRNULO")

		archivo.Seek(particion.Start, 0)
		var buffer bytes.Buffer
		binary.Write(&buffer, binary.BigEndian, &ebrNull)
		structs.EscribirArchivo(archivo, buffer.Bytes())
	}

}

func (self Fdisk) tieneErroresParticionPrimExt(mbr structs.Mbr) bool {
	errores := false

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

	archivo, err := os.OpenFile(self.Path, os.O_RDWR, 0777)
	defer archivo.Close()

	if err != nil {
		fmt.Println("Error: No se ha podido abrir el archivo")
		log.Fatal(err)
	}

	//BUSCAR PARTICION EXTENDIDA
	apuntador := int64(-1)
	for i := 0; i < 4; i++ {
		if mbr.Particion[i].Type == 'E' {
			apuntador = mbr.Particion[i].Start
		}
	}

	ebrActual := structs.GetEbr(archivo, apuntador)
	fmt.Println(structs.UintToString(ebrActual.Name))
	for ebrActual.Next != 0 {
		apuntador = ebrActual.Next
		ebrActual = structs.GetEbr(archivo, apuntador)
	}
	fmt.Println(apuntador)

	var ebr structs.Ebr
	ebr.Status = '0'
	ebr.Fit = structs.GetFit(self.Fit)
	ebr.Start = apuntador + int64(unsafe.Sizeof(ebr))
	ebr.Size = structs.GetSize(self.Size, self.Unit)
	ebr.Name = structs.GetName(self.Name)
	ebr.Next = ebr.Start + ebr.Size

	fmt.Println("###############################")
	fmt.Println(structs.UintToString(ebr.Name))
	fmt.Println(ebr.Start)
	fmt.Println(unsafe.Sizeof(ebr))
	fmt.Println(ebr.Size)
	fmt.Println(string(ebr.Fit))
	fmt.Println("###############################")

	archivo.Seek(apuntador, 0)
	var buffer bytes.Buffer
	binary.Write(&buffer, binary.BigEndian, &ebr)
	structs.EscribirArchivo(archivo, buffer.Bytes())

}

func (self Fdisk) tieneErroresParticionLogica(mbr structs.Mbr) bool {
	errores := false
	return errores
}
