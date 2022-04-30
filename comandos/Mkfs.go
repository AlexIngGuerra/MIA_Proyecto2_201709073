package comandos

import (
	"MIA/structs"
	"bytes"
	"encoding/binary"
	"fmt"
	"os"
	"strconv"
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

//EJECUTAR EL COMANDO
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

//VERIFICAR ERRORES DE PARAMETROS Y OTROS
func (self Mkfs) tieneErrores() bool {
	errores := false
	if self.Id == "" {
		errores = true
		fmt.Println("Error: El parametro -id es obligatorio")
	}

	return errores
}

//FORMATEO FULL, LIMPIA TODO EL DISCO
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

	fmt.Println("Inicio: " + strconv.Itoa(int(part.Start)) + " , Size: " + strconv.Itoa(int(part.Size)))

	superBloque := self.GenerarSuperBloque(n, part.Start)
	self.limpiarEspacioDisco(archivo, part.Start, part.Size)

	//Creamos el usuario root
	self.crearUsuarioRoot(archivo, superBloque, n, part.Start)
}

//FORMATEO FAST, LIMPIA SOLAMENTE LOS BITMAPS
func (self Mkfs) formateoFast(mount ParticionMontada) {
}

//GENERAR UN SUPERBLOQUE DEFAULT
func (self Mkfs) GenerarSuperBloque(n int32, inicioPart int64) structs.SuperBloque {
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
	superBloque.Bm_block_start = inicioPart + int64(unsafe.Sizeof(superBloque))
	superBloque.Bm_inode_start = superBloque.Bm_block_start + 3*int64(n)
	superBloque.Block_start = superBloque.Bm_inode_start + int64(n)
	superBloque.Inode_start = superBloque.Block_start + 3*int64(n)*int64(64)
	//PRIMEROS LIBRES
	superBloque.First_bloc = superBloque.Block_start
	superBloque.First_ino = superBloque.Inode_start

	return superBloque
}

//LIMPIAR UN ESPACIO DE DISCO CON CERO
func (self Mkfs) limpiarEspacioDisco(archivo *os.File, inicio int64, size int32) {
	var cero uint8
	cero = 0

	archivo.Seek(inicio, 0)
	buffer := bytes.Buffer{}
	binary.Write(&buffer, binary.BigEndian, &cero)
	for i := int32(0); i < size; i++ {
		structs.EscribirArchivo(archivo, buffer.Bytes())
	}
}

//GENERAR LA CARPETA Y ARCHIVO USERS.TXT PARA LOS USUARIOS
func (self Mkfs) crearUsuarioRoot(archivo *os.File, superbloque structs.SuperBloque, n int32, partStart int64) {
	self.imprimirSB(superbloque, n)

	//Paso 1: Generar Inodo Raiz "/"
	inodoR := structs.NewInodo(1, 1, 0) //User:1, Grupo:1, Tipo:0
	inodoU := structs.NewInodo(1, 1, 1) //User:1, Grupo:1, Tipo:1

	bloqueC := structs.NewPrimerBloqueCarpeta(0, 0)
	bloqueC.Contenido[2] = structs.Contenido{Name: structs.GetNameBloque("users.txt"), Apuntador: 1}
	inodoR.Block[0] = 0

	root := "1, G, root\n1, U, root, root, 123\n"
	bloquesA := structs.EscribirTextoEnBloques(root)
	inodoU.Block[0] = 1

	superbloque = structs.EscribirInodo(archivo, superbloque, inodoR, n)
	superbloque = structs.EscribirInodo(archivo, superbloque, inodoU, n)
	superbloque = structs.EscribirBloqueC(archivo, superbloque, bloqueC, n)
	superbloque = structs.EscribirBloqueA(archivo, superbloque, bloquesA[0], n)

	structs.EscribirSuperBloque(archivo, superbloque, partStart)

	self.imprimirSB(superbloque, n)
}

func (self Mkfs) imprimirSB(sb structs.SuperBloque, n int32) {
	fmt.Print("N: ")
	fmt.Println(n)

	fmt.Print("InodesCount: ")
	fmt.Println(sb.Inodes_count)

	fmt.Print("BlockCount: ")
	fmt.Println(sb.Blocks_count)

	fmt.Print("FreeInodes: ")
	fmt.Println(sb.Free_inodes_count)

	fmt.Print("FreeBlocks: ")
	fmt.Println(sb.Free_blocks_count)

	fmt.Print("FirstInode: ")
	fmt.Println(sb.First_ino)

	fmt.Print("FirsBlock: ")
	fmt.Println(sb.First_bloc)

	fmt.Print("bmInodeStart: ")
	fmt.Println(sb.Bm_inode_start)

	fmt.Print("bmBloqueStart: ")
	fmt.Println(sb.Bm_block_start)

	fmt.Print("InodeStart: ")
	fmt.Println(sb.Inode_start)

	fmt.Print("BlockStart: ")
	fmt.Println(sb.Block_start)

	fmt.Print("\n")
}
