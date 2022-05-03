package comandos

import (
	"MIA/structs"
	"bytes"
	"encoding/binary"
	"fmt"
	"io/ioutil"
	"os"
	"strconv"
	"strings"
)

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
	n := structs.GetN(part.Size)

	if !self.R {
		self.crearArchivoSinPadre(archivo, superBloque, n, part.Start, self.Path)
	} else {
		self.crearArchivosConPadres(archivo, superBloque, n, part.Start, self.Path)
	}

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

//Ajusta el contenido del archivo
func (self Mkfile) getContenido() string {
	contenido := ""
	if self.Cont == "" {
		for i := 0; i < self.Size; i++ {
			contenido += strconv.Itoa(i % 10)
		}
		return contenido
	} else {
		archivo, err := ioutil.ReadFile(self.Cont)
		if err != nil {
			fmt.Println("Error: al leer el archivo para el contenido")
			if self.Size != 0 {
				for i := 0; i < self.Size; i++ {
					contenido += strconv.Itoa(i % 10)
				}
			}
			return contenido
		}
		contenido = string(archivo)
		return contenido
	}
}

//CREA EL EL ARCHIVO SIN TOMAR EN CUENTA LOS PADRES
func (self Mkfile) crearArchivoSinPadre(archivo *os.File, superbloque structs.SuperBloque, n int32, partStart int64, Path string) structs.SuperBloque {
	carpetas := strings.Split(Path, "/")
	carpetaNueva := carpetas[len(carpetas)-1]

	apuntador := superbloque.Inode_start
	inodo := structs.LeerInodo(archivo, apuntador)
	numInoActual := int32(0)
	numInoAnterior := int32(0)

	for c := 1; c < len(carpetas)-1; c++ {
		ino, ap, num := structs.GetSiguienteInodo(archivo, superbloque, inodo, carpetas[c])
		inodo = ino
		apuntador = ap
		numInoAnterior = numInoActual
		numInoActual = num
	}

	if apuntador == -1 || inodo.Uid == 0 || inodo.Gid == 0 {
		fmt.Println("Error: carpetas padre faltantes")
		return superbloque
	}

	if structs.ExisteNombreInodo(archivo, superbloque, inodo, carpetaNueva) {
		fmt.Println("Error: La carpeta que solicita crear ya existe")
		return superbloque
	}

	num := structs.BuscarBitLibre(archivo, int64(superbloque.Bm_inode_start), superbloque.Inodes_count)

	inodo, superbloque = structs.AgregarNuevoApuntador(archivo, superbloque, inodo, numInoActual, carpetaNueva, num, numInoAnterior)

	//Actualizamos el inodo anterior
	archivo.Seek(apuntador, 0)
	buffer := bytes.Buffer{}
	binary.Write(&buffer, binary.BigEndian, &inodo)
	structs.EscribirArchivo(archivo, buffer.Bytes())

	//Creamos el inodo nuevo
	inodoNuevo := structs.NewInodo(int32(Logeado.Gid), int32(Logeado.Uid), 1)
	superbloque = structs.EscribirInodo(archivo, superbloque, inodoNuevo)
	//Escribimos su contenido

	bloques := structs.EscribirTextoEnBloques(self.getContenido())

	for blc := 0; blc < len(bloques); blc++ {
		if blc >= 16 {
			break
		}

		numBloc := inodoNuevo.Block[blc]
		num := structs.BuscarBitLibre(archivo, superbloque.Bm_block_start, superbloque.Blocks_count)

		if numBloc == -1 {
			superbloque = structs.EscribirBloqueA(archivo, superbloque, bloques[blc])
			inodoNuevo.Block[blc] = num

		} else {

			archivo.Seek(superbloque.Block_start+int64(numBloc)*int64(superbloque.Block_Size), 0)
			buffer := bytes.Buffer{}
			binary.Write(&buffer, binary.BigEndian, &bloques[blc])
			structs.EscribirArchivo(archivo, buffer.Bytes())
		}

	}
	inodoNuevo.Size = int32(len(bloques)) * 64
	posInodo := superbloque.Inode_start + int64(superbloque.Inode_Size)*int64(num)
	archivo.Seek(posInodo, 0)
	buffer = bytes.Buffer{}
	binary.Write(&buffer, binary.BigEndian, &inodoNuevo)
	structs.EscribirArchivo(archivo, buffer.Bytes())

	//Fin escribir contenido
	structs.EscribirSuperBloque(archivo, superbloque, partStart)
	fmt.Println("Archivo creada")
	return superbloque
}

func (self Mkfile) crearArchivosConPadres(archivo *os.File, superbloque structs.SuperBloque, n int32, partStart int64, Path string) {
	mkdir := Mkdir{}
	carpetas := strings.Split(Path, "/")
	dirActual := ""
	for i := 1; i < len(carpetas)-1; i++ {
		dirActual += "/" + carpetas[i]
		fmt.Println(dirActual)
		superbloque = mkdir.CrearCarpetaSinPadres(archivo, superbloque, n, partStart, dirActual)
		structs.EscribirSuperBloque(archivo, superbloque, partStart)
	}
	superbloque = self.crearArchivoSinPadre(archivo, superbloque, n, partStart, Path)
	structs.EscribirSuperBloque(archivo, superbloque, partStart)
}
