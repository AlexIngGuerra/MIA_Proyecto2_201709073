package comandos

import (
	"MIA/structs"
	"bytes"
	"encoding/binary"
	"fmt"
	"os"
	"time"
	"unsafe"
)

type Mkfs struct {
	Id   string
	Type string
}

func NewMkfs() Mkfs {
	return Mkfs{Id: "", Type: "full"}
}

func (self Mkfs) Ejecutar() {
	fmt.Println("Ejecutando mkfs")
	if self.tieneErrores() {
		return
	}

	mount := GetMount(self.Id)
	if mount.Id == "" {
		fmt.Println("Error: El id solicitado no corresponde a ninguna particion montada")
		return
	}

	if self.Type == "full" {
		self.formateoFull(mount)
	} else if self.Type == "fast" {
		self.formateoFast(mount)
	}
	fmt.Print("\n")
}

func (self Mkfs) tieneErrores() bool {
	errores := false
	if self.Id == "" {
		errores = true
		fmt.Println("Error: El parametro -id es obligatorio")
	}

	return errores
}

//Formateo limpia todo el disco
func (self Mkfs) formateoFull(mount ParticionMontada) {

	archivo, err := os.OpenFile(mount.Path, os.O_RDWR, 0777)
	defer archivo.Close()
	if err != nil {
		fmt.Println("Error: No se ha podido abrir el archivo")
		return
	}

	mbr := structs.GetMbr(mount.Path)
	part := structs.GetParticion(mount.Name, mbr, archivo)
	n := structs.GetN(part.Size)
	apuntador := part.Start + int64(unsafe.Sizeof(structs.SuperBloque{})) //Lo dejamos al inicio del bitmap bloques
	apuntador = apuntador + 1
	superBloque := self.GenerarSuperBloque(n, apuntador)
	fmt.Println(apuntador)
	fmt.Println(superBloque.Inode_Size)

	self.crearUsuarioRoot(archivo)
}

//Formateo que solo limpia los bitmaps
func (self Mkfs) formateoFast(mount ParticionMontada) {

}

func (self Mkfs) GenerarSuperBloque(n int32, apuntador int64) structs.SuperBloque {
	var superBloque structs.SuperBloque

	superBloque.FileSystem_Type = 2
	// CONTEO DE INODOS Y BLOQUES
	superBloque.Inodes_count = n
	superBloque.Blocks_count = 3 * n
	//CONTEO DE INODOS Y BLOQUES LIBRES
	superBloque.Free_blocks_count = 3 * n
	superBloque.Free_inodes_count = n
	//TIEMPO
	currentTime := time.Now()
	superBloque.Mtime = structs.GetFecha(currentTime.Format("2006-01-02 15:04:05"))
	superBloque.Mnt_count = 1
	superBloque.Magic = 61267
	//INFO DE STRUCTS
	superBloque.Inode_Size = int32(unsafe.Sizeof(structs.Inodo{}))
	superBloque.Block_Size = 64
	//INICIO DE STRUCTS
	superBloque.Bm_block_start = apuntador
	superBloque.Bm_inode_start = apuntador + 3*int64(n) + 1
	superBloque.Block_start = apuntador + 3*int64(n) + int64(n) + 1
	superBloque.Inode_start = apuntador + 3*int64(n) + int64(n) + 3*int64(n)*64 + 1
	//PRIMEROS LIBRES
	superBloque.First_bloc = superBloque.Block_start
	superBloque.First_ino = superBloque.Inode_start

	return superBloque

}

func (self Mkdisk) limpiarEspacioDisco(archivo *os.File, inicio int64, size int64) {
	var cero uint8
	cero = '0'

	var buffer bytes.Buffer
	binary.Write(&buffer, binary.BigEndian, &cero)
	archivo.Seek(inicio, 0)
	for i := int64(0); i < size; i++ {
		structs.EscribirArchivo(archivo, buffer.Bytes())
	}

}

func (self Mkfs) crearUsuarioRoot(archivo *os.File) {
	root := "1, G, root\n1, U, root, root, 123\n"
	bloques := structs.EscribirBloqueArchivo((root))
	contenido := structs.LeerBloquesArchivo(bloques)
	grupos := structs.GetGruposYUsuarios(contenido)
	fmt.Println(grupos)
}
