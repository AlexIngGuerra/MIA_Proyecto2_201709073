package comandos

import (
	"MIA/structs"
	"fmt"
	"os"
)

type Login struct {
	Usuario  string
	Password string
	Id       string
}

type Usuario struct {
	Usuario string
	Uid     int
	Gid     int
	Id      string
}

var Logeado Usuario

func NewLogin() Login {
	return Login{Usuario: "", Password: "", Id: ""}
}

//EJECUTAR COMANDO
func (self Login) Ejecutar() {
	if self.tieneErrores() {
		return
	}

	mount := GetMount(self.Id)
	if mount.Id == "" {
		fmt.Println("Error: El id solicitado no corresponde a ninguna particion montada")
		return
	}

	archivo, err := os.OpenFile(mount.Path, os.O_RDWR, 0777)
	defer archivo.Close()
	if err != nil {
		fmt.Println("Error: No se ha podido abrir el archivo")
		return
	}

	mbr := structs.GetMbr(mount.Path)
	part := structs.GetParticion(mount.Name, mbr, archivo)
	superBloque := structs.LeerSuperBloque(archivo, part.Start)
	n := structs.GetN(part.Size)

	self.buscarUsuario(archivo, superBloque, n)
	fmt.Print("\n")
}

//VERIFICAR ERRORES
func (self Login) tieneErrores() bool {
	errores := false
	if self.Usuario == "" {
		errores = true
		fmt.Println("Error: El parametro usuario es obligatorio")
	}

	if self.Password == "" {
		errores = true
		fmt.Println("Error: El parametro password es obligatorio")
	}

	if self.Id == "" {
		errores = true
		fmt.Println("Error: El parametro id es obligatorio")
	}

	if Logeado.Gid != 0 && Logeado.Uid != 0 {
		fmt.Println("Error: Usted debe cerrar sesión para ingresar con otro usuario")
		errores = true
	}

	return errores
}

//BUSCA EL USUARIO EN EL USERS.TXT Y LO LOGEA
func (self Login) buscarUsuario(archivo *os.File, superBloque structs.SuperBloque, n int32) {
	if !structs.ExisteStructEnBM(archivo, superBloque.Bm_inode_start, n, 0) {
		fmt.Println("Error: No se ha creado el inodo raiz, debe formatear el disco")
		return
	}

	inodoRaiz := structs.LeerInodo(archivo, superBloque.Inode_start)
	bloqueC := structs.LeerBloqueC(archivo, superBloque.Block_start+int64(inodoRaiz.Block[0])*64)
	inodoUsers := structs.LeerInodo(archivo, superBloque.Inode_start+int64(bloqueC.Contenido[2].Apuntador)*int64(superBloque.Inode_Size))

	usuarios := structs.LeerBloquesArchivo(structs.ObtenerBloquesArchivo(archivo, superBloque, inodoUsers))
	grupos := structs.GetGruposYUsuarios(usuarios)

	loginExitoso := false
	for i := 0; i < len(grupos); i++ {

		for n := 0; n < len(grupos[i].Users); n++ {

			if self.Usuario == grupos[i].Users[n].User &&
				self.Password == grupos[i].Users[n].Password {
				Logeado = Usuario{Usuario: self.Usuario, Uid: grupos[i].Users[n].Id, Gid: grupos[i].Id, Id: self.Id}
				loginExitoso = true
				break
			}

		}

		if loginExitoso {
			break
		}

	}

	if loginExitoso {
		fmt.Println("Usted ha ingresado con el usuario " + Logeado.Usuario)
	} else {
		Logeado = Usuario{}
		fmt.Println("Error: No se ha podio ingresaro con " + self.Usuario)
	}
}

func (self Login) Logout() {
	if Logeado.Gid == 0 && Logeado.Uid == 0 {
		fmt.Println("Error: No ha ingresado con ningun usuario")
		return
	}

	Logeado = Usuario{}
	fmt.Println("Sesión finalizada")
}
