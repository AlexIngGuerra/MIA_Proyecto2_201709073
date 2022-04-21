package comandos

import (
	"fmt"
	"os"
	"strings"
)

type Rep struct {
	Name string
	Path string
	Id   string
	Ruta string
}

func NewRep() Rep {
	return Rep{Name: "", Path: "", Id: "", Ruta: ""}
}

func (self Rep) Ejecutar() {

}

func (self Rep) getDir(Path string) string {
	cadena := strings.Split(Path, "/")
	dir := cadena[0]
	for i := 1; i < len(cadena)-1; i++ {
		if i < len(cadena) {
			dir = dir + "/" + cadena[i]
		}
	}
	return dir
}

func (self Rep) crearCarpetaSiNoExiste(Path string) {
	_, err := os.Stat(Path)
	if os.IsNotExist(err) {
		fmt.Println("Aviso: La carpeta o carpetas no existen, se procede a crearlas.")
		err = os.MkdirAll(Path, 0777)
		if err != nil {
			fmt.Println("Error: No se ha podido crear la carpeta.")
		}
	}
}

func (self Rep) crearArchivoSiNoExiste(Path string) bool {
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

func (self Rep) repDisk() {

}
