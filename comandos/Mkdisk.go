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

	self.crearCarpetaSiNoExiste(self.getDir(self.Path)) //Crear las carpetas
	if self.crearArchivoSiNoExiste(self.Path) {         //Si se pudo crear el archivo se crea el disco
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

func (self Mkdisk) getDir(Path string) string {
	cadena := strings.Split(Path, "/")
	dir := cadena[0]
	for i := 1; i < len(cadena)-1; i++ {
		if i < len(cadena) {
			dir = dir + "/" + cadena[i]
		}
	}
	return dir
}

func (self Mkdisk) crearCarpetaSiNoExiste(Path string) {
	_, err := os.Stat(Path)
	if os.IsNotExist(err) {
		fmt.Println("Aviso: La carpeta o carpetas no existen, se procede a crearlas.")
		err = os.MkdirAll(Path, 0777)
		if err != nil {
			fmt.Println("Error: No se ha podido crear la carpeta.")
		}
	}
}

func (self Mkdisk) crearArchivoSiNoExiste(Path string) bool {
	if _, err := os.Stat(Path); os.IsNotExist(err) {
		archivo, err := os.Create(Path)
		defer archivo.Close()
		if err != nil {
			fmt.Println("Error: No se ha podido crear el archivo")
		}
		return true
	}
	fmt.Println("Error: No se puede crear el disco porque el archivo ya existe.")
	return false
}
