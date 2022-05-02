package structs

import "time"

//STRUCTS UTILIZADOS PARA EL EXT2---------------------------------------------
type SuperBloque struct {
	FileSystem_Type   int32
	Inodes_count      int32
	Blocks_count      int32
	Free_blocks_count int32
	Free_inodes_count int32
	Mtime             [20]uint8
	Magic             int32
	Mnt_count         int32
	Inode_Size        int32
	Block_Size        int32
	First_ino         int64
	First_bloc        int64
	Bm_inode_start    int64
	Bm_block_start    int64
	Inode_start       int64
	Block_start       int64
}

type Inodo struct {
	Uid   int32
	Gid   int32
	Size  int32
	Atime [20]uint8 //Acceso sin modificar
	Ctime [20]uint8 //Fecha Creacion
	Mtime [20]uint8 //Fecha Modificaion
	Block [16]int32
	Type  uint8
	Perm  int32
}

type BloqueCarpeta struct {
	Contenido [4]Contenido
}

type Contenido struct {
	Name      [12]uint8
	Apuntador int32
}

type BloqueArchivo struct {
	Contenido [64]uint8
}

//CONSTRUCTORES -------------------------------------------------
func NewInodo(Uid int32, Gid int32, tipo uint8) Inodo {
	inodo := Inodo{}

	inodo.Uid = Uid
	inodo.Gid = Gid
	inodo.Size = 0

	currentTime := time.Now()
	inodo.Atime = GetFecha(currentTime.Format("2006-01-02 15:04:05"))
	inodo.Ctime = GetFecha(currentTime.Format("2006-01-02 15:04:05"))
	inodo.Mtime = GetFecha(currentTime.Format("2006-01-02 15:04:05"))

	inodo.Type = tipo

	for i := 0; i < 16; i++ {
		inodo.Block[i] = -1
	}

	inodo.Perm = 664

	return inodo
}

func NewBloqueCarpeta() BloqueCarpeta {
	carpeta := BloqueCarpeta{}
	for i := 0; i < 4; i++ {
		carpeta.Contenido[i].Apuntador = -1
	}
	return carpeta
}

func NewPrimerBloqueCarpeta(inodoActual int32, inodoAnterior int32) BloqueCarpeta {
	carpeta := BloqueCarpeta{}
	carpeta.Contenido[0] = Contenido{Name: GetNameBloque("."), Apuntador: inodoActual}
	carpeta.Contenido[1] = Contenido{Name: GetNameBloque(".."), Apuntador: inodoAnterior}

	for i := 2; i < 4; i++ {
		carpeta.Contenido[i].Apuntador = -1
	}

	return carpeta
}

//STRUCTS PARA EL MANEJO DE INSTRUCCIONES -------------------------
type Grupo struct {
	Id    int
	Name  string
	Users []Usuario
}

type Usuario struct {
	Id       int
	User     string
	Password string
}
