package comandos

import "fmt"

type Mkfile struct {
	Path string
	R    bool
	Size int
	Cont string
}

func NewMkfile() Mkfile {
	return Mkfile{Path: "", R: false, Size: 0, Cont: ""}
}

func (self Mkfile) Ejecutar() {
	if self.tieneErrores() {
		return
	}
	fmt.Println("Iniciando creacion de archivo")

}

//Verificar Errores
func (self Mkfile) tieneErrores() bool {
	errores := false
	if self.Path == "" {
		fmt.Println("Error: El parametro path es obligatorio")
		errores = true
	}

	if Logeado.Gid == 0 && Logeado.Uid == 0 {
		fmt.Println("Error: Usted debe iniciar sesión para ejecutar este comando")
		errores = true
	}
	return errores
}
