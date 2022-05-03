package comandos

import (
	"MIA/structs"
	"bytes"
	"encoding/binary"
	"fmt"
	"os"
)

type Mkuser struct {
	Usuario  string
	Password string
	Group    string
}

func NewMkuser() Mkuser {
	return Mkuser{Usuario: "", Password: "", Group: ""}
}

func (self Mkuser) Ejecutar() {
	if self.tieneErrores() {
		return
	}
	fmt.Println("Creando nuevo usuario")

	mount := GetMount(Logeado.Id)
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

	posInodo := superBloque.Inode_start + int64(superBloque.Inode_Size)
	inodoUsers := structs.LeerInodo(archivo, posInodo)

	bloques := structs.ObtenerBloquesArchivo(archivo, superBloque, inodoUsers)
	groups := structs.GetGruposYUsuarios(structs.LeerBloquesArchivo(bloques))

	var user structs.Usuario
	user.Password = self.Password
	user.User = self.Usuario
	user.Id = 0

	groups = structs.AddUsuario(groups, user, self.Group)
	txt := structs.GroupsToString(groups)

	bloques = structs.EscribirTextoEnBloques(txt)

	for blc := 0; blc < len(bloques); blc++ {
		if blc >= 16 {
			break
		}

		numBloc := inodoUsers.Block[blc]
		num := structs.BuscarBitLibre(archivo, superBloque.Bm_block_start, superBloque.Blocks_count)

		if numBloc == -1 {
			superBloque = structs.EscribirBloqueA(archivo, superBloque, bloques[blc])
			inodoUsers.Block[blc] = num

		} else {

			archivo.Seek(superBloque.Block_start+int64(numBloc)*int64(superBloque.Block_Size), 0)
			buffer := bytes.Buffer{}
			binary.Write(&buffer, binary.BigEndian, &bloques[blc])
			structs.EscribirArchivo(archivo, buffer.Bytes())
		}

	}
	fmt.Println(txt)

	archivo.Seek(posInodo, 0)
	buffer := bytes.Buffer{}
	binary.Write(&buffer, binary.BigEndian, &inodoUsers)
	structs.EscribirArchivo(archivo, buffer.Bytes())

	structs.EscribirSuperBloque(archivo, superBloque, part.Start)
	fmt.Println("Finalizada la Ejecución del comando")
	fmt.Print("\n")
}

func (self Mkuser) tieneErrores() bool {
	errores := false
	if self.Usuario == "" {
		errores = true
		fmt.Println("Error: El parametro -name es obligatorio")
	}

	if self.Password == "" {
		errores = true
		fmt.Println("Error: El parametro -name es obligatorio")
	}

	if self.Group == "" {
		errores = true
		fmt.Println("Error: El parametro -name es obligatorio")
	}

	if Logeado.Gid != 1 && Logeado.Uid != 1 {
		errores = true
		fmt.Println("Error: Debe haber ingresado como usuario root")
	}

	return errores
}

func (self Mkuser) EliminarUsuario() {
	if self.Usuario == "" {
		fmt.Println("Error: El parametro name es obligatorio")
		return
	}

	fmt.Println("Eliminando Usuario...")
	mount := GetMount(Logeado.Id)
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

	posInodo := superBloque.Inode_start + int64(superBloque.Inode_Size)
	inodoUsers := structs.LeerInodo(archivo, posInodo)

	bloques := structs.ObtenerBloquesArchivo(archivo, superBloque, inodoUsers)
	groups := structs.GetGruposYUsuarios(structs.LeerBloquesArchivo(bloques))

	eliminado := false
	for i := 0; i < len(groups); i++ {

		for u := 0; u < len(groups[i].Users); u++ {
			if groups[i].Users[u].User == self.Usuario {
				groups[i].Users[u].Id = 0
				eliminado = true
			}
		}
	}

	if !eliminado {
		fmt.Println("Error: No se ha encontrado el grupo a eliminar")
		return
	}

	txt := structs.GroupsToString(groups)

	bloques = structs.EscribirTextoEnBloques(txt)

	for blc := 0; blc < len(bloques); blc++ {
		if blc >= 16 {
			break
		}

		numBloc := inodoUsers.Block[blc]
		num := structs.BuscarBitLibre(archivo, superBloque.Bm_block_start, superBloque.Blocks_count)

		if numBloc == -1 {
			superBloque = structs.EscribirBloqueA(archivo, superBloque, bloques[blc])
			inodoUsers.Block[blc] = num

		} else {

			archivo.Seek(superBloque.Block_start+int64(numBloc)*int64(superBloque.Block_Size), 0)
			buffer := bytes.Buffer{}
			binary.Write(&buffer, binary.BigEndian, &bloques[blc])
			structs.EscribirArchivo(archivo, buffer.Bytes())
		}

	}
	fmt.Println(txt)

	archivo.Seek(posInodo, 0)
	buffer := bytes.Buffer{}
	binary.Write(&buffer, binary.BigEndian, &inodoUsers)
	structs.EscribirArchivo(archivo, buffer.Bytes())

	structs.EscribirSuperBloque(archivo, superBloque, part.Start)
	fmt.Println("Finalizada la Ejecución del comando")
	fmt.Print("\n")
}
