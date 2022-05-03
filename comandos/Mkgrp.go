package comandos

import (
	"MIA/structs"
	"bytes"
	"encoding/binary"
	"fmt"
	"os"
)

type Mkgrp struct {
	Name string
}

func NewMkgrp() Mkgrp {
	return Mkgrp{Name: ""}
}

//EJECUTAR EL COMANDO
func (self Mkgrp) Ejecutar() {
	if self.tieneErrores() {
		return
	}
	fmt.Println("Creando nuevo grupo:")

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

	var grupoNuevo structs.Grupo
	grupoNuevo.Id = len(groups) + 1
	grupoNuevo.Name = self.Name

	groups = structs.AddGrupo(groups, grupoNuevo)
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

//VERIFICAR SI SE PUEDE EJECUTAR EL COMANDO
func (self Mkgrp) tieneErrores() bool {
	errores := false
	if self.Name == "" {
		errores = true
		fmt.Println("Error: El parametro -name es obligatorio")
	}

	if Logeado.Gid != 1 && Logeado.Uid != 1 {
		errores = true
		fmt.Println("Error: Debe haber ingresado como usuario root")
	}

	return errores
}

func (self Mkgrp) EliminarGrupo() {
	if self.Name == "" {
		fmt.Println("Error: El parametro name es obligatorio")
		return
	}
	fmt.Println("Eliminando Grupo...")
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
		if groups[i].Name == self.Name {
			groups[i].Id = 0

			for u := 0; u < len(groups[i].Users); u++ {
				groups[i].Users[u].Id = 0
			}
			eliminado = true
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
