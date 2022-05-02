package comandos

import (
	"MIA/structs"
	"bytes"
	"encoding/binary"
	"fmt"
	"os"
	"strings"
)

type Mkdir struct {
	Path string
	P    bool
}

func NewMkdir() Mkdir {
	return Mkdir{Path: "", P: false}
}

//EJECUTAR EL COMANDO
func (self Mkdir) Ejecutar() {
	fmt.Println("Mkdir ejecutando")
	if self.tieneErrores() {
		return
	}

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

	if !self.P {
		self.CrearCarpetaSinPadres(archivo, superBloque, n, part.Start)
	}

	fmt.Print("\n")
}

//VERIFICAR SI EL COMANDO SE PUEDE EJECUTAR
func (self Mkdir) tieneErrores() bool {
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

//CREA LA CARPETA SOLICITADA (NO CREA LAS CARPETAS PADRES SI NO EXISTEN)
func (self Mkdir) CrearCarpetaSinPadres(archivo *os.File, superbloque structs.SuperBloque, n int32, partStart int64) {

	carpetas := strings.Split(self.Path, "/")
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
		return
	}

	if structs.ExisteNombreInodo(archivo, superbloque, inodo, carpetaNueva) {
		fmt.Println("Error: La carpeta que solicita crear ya existe")
		return
	}

	num := structs.BuscarBitLibre(archivo, int64(superbloque.Bm_inode_start), superbloque.Inodes_count)

	inodo, superbloque = structs.AgregarNuevoApuntador(archivo, superbloque, inodo, numInoActual, carpetaNueva, num, numInoAnterior)

	//Actualizamos el inodo anterior
	archivo.Seek(apuntador, 0)
	buffer := bytes.Buffer{}
	binary.Write(&buffer, binary.BigEndian, &inodo)
	structs.EscribirArchivo(archivo, buffer.Bytes())

	//Creamos el inodo nuevo
	inodoNuevo := structs.NewInodo(int32(Logeado.Gid), int32(Logeado.Uid), 0)
	superbloque = structs.EscribirInodo(archivo, superbloque, inodoNuevo)
	structs.EscribirSuperBloque(archivo, superbloque, partStart)
	fmt.Println("Carpeta creada")
}
