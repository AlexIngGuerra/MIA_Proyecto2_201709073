package comandos

import (
	"fmt"
	"io/ioutil"
)

type Exec struct {
	Path string
}

func NewExec() Exec {
	return Exec{Path: ""}
}

//EJECUTAR COMANDO
func (self Exec) Ejecutar() string {
	contenido := ""
	if !self.tieneErrores() {
		archivo, _ := ioutil.ReadFile(self.Path)
		contenido = string(archivo)
		return contenido
	}
	return ""
}

//VERIFICAR ERRORES DE PARAMETROS
func (self Exec) tieneErrores() bool {
	errores := false
	if self.Path == "" {
		fmt.Println("Error: Exec debe tener el parametro Path")
		errores = true
	}
	return errores
}
