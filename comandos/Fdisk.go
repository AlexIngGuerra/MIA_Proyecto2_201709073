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
	if self.Size <= 0 {
		errores = true
		fmt.Println("Error: El size tiene que ser mayor a cero")
	}

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

func (self Fdisk) crearParticionPrimaria(mbr structs.Mbr) {
	if self.tieneErroresParticionPrimExt(mbr) {
		return
	}

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
	particion.Start = self.getStartPrimaria(mbr) + 1

	for i := 0; i < 4; i++ {
		if mbr.Particion[i].Size == 0 {
			mbr.Particion[i] = particion
			break
		}
	}
	particionCreada := false
	ebrNuloCreado := true

	archivo.Seek(0, 0)
	var buffer bytes.Buffer
	binary.Write(&buffer, binary.BigEndian, &mbr)
	particionCreada = structs.EscribirArchivo(archivo, buffer.Bytes())

	if particion.Type == 'E' {
		var ebrNull structs.Ebr
		ebrNull.Name = structs.GetName("EBRNULO")

		archivo.Seek(particion.Start, 0)
		var buffer bytes.Buffer
		binary.Write(&buffer, binary.BigEndian, &ebrNull)
		ebrNuloCreado = structs.EscribirArchivo(archivo, buffer.Bytes())
	}

	if particionCreada && ebrNuloCreado {
		fmt.Println("Se ha creado la particion: " + self.Name)
	}

}

func (self Fdisk) tieneErroresParticionPrimExt(mbr structs.Mbr) bool {
	errores := false
	//Revisar Nombre
	if self.nombreRepetido(mbr) {
		errores = true
		fmt.Println("Error: El Nombre de la particion esta repetido")
	}

	//Revisar espacio y si hay espacio libre para crear la particion
	espacioLibre := mbr.Tamano - 1
	libreParticion := false
	existeExtendida := false
	for i := 0; i < 4; i++ {
		espacioLibre -= mbr.Particion[i].Size
		if mbr.Particion[i].Size == 0 {
			libreParticion = true
		}

		if mbr.Particion[i].Type == 'E' {
			existeExtendida = true
		}
	}

	if espacioLibre < structs.GetSize(self.Size, self.Unit) {
		errores = true
		fmt.Println("Error: No hay suficiente espacio en el disco para crear la particion")
	}

	if !libreParticion {
		errores = true
		fmt.Println("Error: El máximo de particiones no logicas es de 4")
	}

	if existeExtendida && self.Type == "E" {
		errores = true
		fmt.Println("Error: Solo puede existir una sola particion extendida")
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
	if self.tieneErroresParticionLogica(mbr) {
		return
	}

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
	ebr.Start = apuntador + int64(unsafe.Sizeof(ebr)) + 1
	ebr.Size = structs.GetSize(self.Size, self.Unit)
	ebr.Name = structs.GetName(self.Name)
	ebr.Next = ebr.Start + ebr.Size

	//Creamos el Ebr de la particion
	particionCreada := false

	archivo.Seek(apuntador, 0)
	var buffer bytes.Buffer
	binary.Write(&buffer, binary.BigEndian, &ebr)
	particionCreada = structs.EscribirArchivo(archivo, buffer.Bytes())

	//Creamos el Ebr nulo para verificar la finalizacion
	apuntador = ebr.Next
	ebrNuloCreado := false

	var ebrNull structs.Ebr
	ebrNull.Name = structs.GetName("EBRNULO")

	archivo.Seek(apuntador, 0)
	var buffer2 bytes.Buffer
	binary.Write(&buffer2, binary.BigEndian, &ebrNull)
	ebrNuloCreado = structs.EscribirArchivo(archivo, buffer2.Bytes())

	if particionCreada && ebrNuloCreado {
		fmt.Println("Se ha creado la particion logica: " + self.Name)
	}

}

func (self Fdisk) tieneErroresParticionLogica(mbr structs.Mbr) bool {
	archivo, err := os.OpenFile(self.Path, os.O_RDWR, 0777)
	defer archivo.Close()
	if err != nil {
		fmt.Println("Error: No se ha podido abrir el archivo")
		log.Fatal(err)
	}

	errores := false
	//Revisar que el nombre no este repetido
	if self.nombreRepetido(mbr) {
		errores = true
		fmt.Println("Error: El Nombre de la particion esta repetido")
	}

	//Existe Particion Extendida
	pos := -1
	for i := 0; i < 4; i++ {
		if mbr.Particion[i].Type == 'E' {
			pos = i
			break
		}
	}

	if pos == -1 {
		errores = true
		fmt.Println("Error: No existe la particion extendida")

		// Si hay espacio suficiente para crear la particion
	} else {
		espacioLibre := mbr.Particion[pos].Start + mbr.Particion[pos].Size
		apuntador := mbr.Particion[pos].Start
		ebrActual := structs.GetEbr(archivo, apuntador)
		for ebrActual.Next != 0 {
			apuntador = ebrActual.Next
			ebrActual = structs.GetEbr(archivo, apuntador)
		}
		espacioLibre -= ebrActual.Start
		espacioLibre -= int64(unsafe.Sizeof(ebrActual))

		if espacioLibre < structs.GetSize(self.Size, self.Unit) {
			errores = true
			fmt.Println("Error: No hay suficiente espacio en la particion extendida")
		}
	}

	return errores
}

func (self Fdisk) nombreRepetido(mbr structs.Mbr) bool {
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
