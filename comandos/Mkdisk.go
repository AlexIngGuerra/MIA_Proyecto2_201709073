package comandos

import (
	"MIA/structs"
	"bytes"
	"encoding/binary"
	"fmt"
	"log"
	"math/rand"
	"os"
	"strings"
	"time"
)

type Mkdisk struct {
	Size int
	Fit  string
	Unit string
	Path string
}

func NewMkdisk() Mkdisk {
	return Mkdisk{Size: 0, Fit: "FF", Unit: "M", Path: ""}
}

func (self Mkdisk) Ejecutar() {
	if self.tieneErrores() {
		return
	}

	self.CrearCarpetaSiNoExiste(self.GetDir(self.Path)) //Crear las carpetas
	if self.CrearArchivoSiNoExiste(self.Path) {
		//Si se pudo crear el archivo se crea el disco
		var mbr structs.Mbr
		mbr.Tamano = structs.GetSize(self.Size, self.Unit)
		mbr.Fit = structs.GetFit(self.Fit)
		mbr.Dks_Signature = rand.Int63()
		currentTime := time.Now()
		mbr.Fecha_Creacion = structs.GetFecha(currentTime.Format("2006-01-02 15:04:05"))

		discoCreado := true

		archivo, err := os.OpenFile(self.Path, os.O_RDWR, 0777)
		defer archivo.Close()
		if err != nil {
			log.Fatal(err)
			discoCreado = false
		}

		tamano := self.Size
		if self.Unit == "M" {
			tamano = tamano * 1024
		}

		//Escribir el MBR
		var bufferMbr bytes.Buffer
		binary.Write(&bufferMbr, binary.BigEndian, &mbr)
		discoCreado = structs.EscribirArchivo(archivo, bufferMbr.Bytes())

		//Llenar de cero el archivo binario
		var bufferCero bytes.Buffer
		var cero [1024]uint8
		binary.Write(&bufferCero, binary.BigEndian, &cero)

		for i := 0; i < tamano; i++ {
			discoCreado = structs.EscribirArchivo(archivo, bufferCero.Bytes())
		}

		if discoCreado {
			fmt.Println("Disco creado correctamente en: " + self.Path)
		}

	}
	fmt.Print("\n")
}

func (self Mkdisk) tieneErrores() bool {
	errores := false
	if self.Size < 1 {
		fmt.Println("Error: El parametro size debe recibir un numero entero positivo")
	}

	if self.Path == "" {
		fmt.Println("Error: El comando mkdisk debe recibir un path")
	}
	return errores
}

func (self Mkdisk) GetDir(Path string) string {
	cadena := strings.Split(Path, "/")
	dir := cadena[0]
	for i := 1; i < len(cadena)-1; i++ {
		if i < len(cadena) {
			dir = dir + "/" + cadena[i]
		}
	}
	return dir
}

func (self Mkdisk) CrearCarpetaSiNoExiste(Path string) {
	_, err := os.Stat(Path)
	if os.IsNotExist(err) {
		fmt.Println("Aviso: La carpeta o carpetas no existen, se procede a crearlas.")
		err = os.MkdirAll(Path, 0777)
		if err != nil {
			fmt.Println("Error: No se ha podido crear la carpeta.")
		}
	}
}

func (self Mkdisk) CrearArchivoSiNoExiste(Path string) bool {
	if _, err := os.Stat(Path); os.IsNotExist(err) {
		archivo, err := os.Create(Path)
		defer archivo.Close()
		if err != nil {
			fmt.Println("Error: No se ha podido crear el archivo")
			return false
		}
		return true
	}
	fmt.Println("Error: No se puede crear el disco porque el archivo ya existe.")
	return false
}

func (self Mkdisk) CrearArchivo(Path string) bool {
	archivo, err := os.Create(Path)
	defer archivo.Close()
	if err != nil {
		fmt.Println("Error: No se ha podido crear el archivo")
		return false
	}
	return true
}

func (self Mkdisk) EjecutarRmdisk() {
	if self.Path != "" {

		fmt.Println("Desea eliminar el disco: y/n")
		var comando string
		fmt.Scanln(&comando)

		if comando == "y" {
			err := os.Remove(self.Path)
			if err != nil {
				fmt.Println("Error: No se ha podido eliminar el disco")
				fmt.Println(err)
			} else {
				fmt.Println("Eliminado el disco en: " + self.Path)
			}
		}

	} else {
		fmt.Println("Error: El comando rmdisk necesita el parametro path")
	}
	fmt.Print("\n")
}
