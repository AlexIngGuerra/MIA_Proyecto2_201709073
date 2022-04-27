package structs

type SuperBloque struct {
	FileSystem_Type   int32
	Inodes_count      int32
	Blocks_count      int32
	Free_blocks_count int32
	Free_inodes_coutn int32
	Mtime             [20]uint8
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
	Atime [20]uint8
	Ctime [20]uint8
	Mtime [20]uint8
	Block [16]int32
	Type  uint8
	Perm  int32
}

type BloqueCarpeta struct {
	Contenido [4]Contenido
}

type Contenido struct {
	Name  [12]uint8
	Inodo int32
}

type BloqueArchivo struct {
	Contenido [64]uint8
}
