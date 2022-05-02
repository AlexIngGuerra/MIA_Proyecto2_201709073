package comandos

import (
	"MIA/structs"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"strconv"
	"unsafe"
)

type Rep struct {
	Name string
	Path string
	Id   string
	Ruta string
}

func NewRep() Rep {
	return Rep{Name: "", Path: "", Id: "", Ruta: ""}
}

//EJECUTAR EL COMANDO
func (self Rep) Ejecutar() {
	if self.tieneErrores() {
		return
	}

	mount := GetMount(self.Id)
	if mount.Id == "" {
		fmt.Println("Error: El id solicitado no corresponde a ninguna particion montada")
		return
	}

	var Mkdisk Mkdisk
	DotPath := Mkdisk.GetDir(self.Path) + "/reporte.dot"
	Mkdisk.CrearCarpetaSiNoExiste(Mkdisk.GetDir(self.Path)) //Crear las carpetas
	if Mkdisk.CrearArchivo(DotPath) {

		if self.Name == "disk" {
			self.repDisk(DotPath, mount)
		} else if self.Name == "tree" {
			self.repTree(DotPath, mount)
		} else if self.Name == "file" {

		}

	}
}

//VERIFICAR ERRORES DE PARAMETROS Y OTROS
func (self Rep) tieneErrores() bool {
	errores := false
	if self.Path == "" {
		errores = true
		fmt.Println("Error: El parametro -path es obligatorio")
	}

	if self.Name == "" {
		errores = true
		fmt.Println("Error: El parametro -name es obligatorio")
	}

	if self.Id == "" {
		errores = true
		fmt.Println("Error: El parametro -id es obligatorio")
	}

	return errores
}

//OBTENER LA PARTICION MONTADA
func GetMount(Id string) ParticionMontada {
	for i := 0; i < len(Montados); i++ {
		if Montados[i].Id == Id {
			return Montados[i]
		}
	}
	return ParticionMontada{}
}

//GENERAR REPORTE DISK
func (self Rep) repDisk(DotPath string, mount ParticionMontada) {
	mbr := structs.GetMbr(mount.Path)
	archivo, err := os.OpenFile(mount.Path, os.O_RDWR, 0777)
	defer archivo.Close()

	if err != nil {
		fmt.Println("Error: No se ha podido abrir el archivo")
		log.Fatal(err)
	}

	if mbr.Size <= 0 {
		fmt.Println("Error: El Mbr no es funcional")
		return
	}

	contenido := "digraph G {\n"
	contenido = contenido + "node_A [shape=record    label=\"MBR"
	porcentaje := float64(100)

	for i := 0; i < 4; i++ {
		particion := mbr.Particion[i]

		if particion.Start > 0 {

			if particion.Type == 'P' {
				contenido += "|Particion Primaria"
				numero := float64(particion.Size) / float64(mbr.Size) * 100
				porcentaje = porcentaje - numero
				var s string = strconv.FormatFloat(numero, 'f', 2, 64)
				contenido = contenido + "\\n" + s + string('%') + " del disco"
			}

			if particion.Type == 'E' {
				contenido += "|{Particion Extendida|{"
				extPorcentaje := float64(particion.Size) / float64(mbr.Size) * 100

				apuntador := mbr.Particion[i].Start
				ebrActual := structs.GetEbr(archivo, apuntador)
				if ebrActual.Size != 0 {
					contenido = contenido + "EBR| Particion Logica"
					numero := (float64(ebrActual.Size) + float64(unsafe.Sizeof(ebrActual))) / float64(mbr.Size) * 100
					porcentaje = porcentaje - numero
					extPorcentaje = extPorcentaje - numero
					var s string = strconv.FormatFloat(numero, 'f', 2, 64)
					contenido = contenido + "\\n" + s + string('%') + " del disco"
				}

				for ebrActual.Size != 0 {
					apuntador = ebrActual.Next
					ebrActual = structs.GetEbr(archivo, apuntador)
					if ebrActual.Size != 0 {
						contenido = contenido + "|EBR| Particion Logica"
						numero := (float64(ebrActual.Size) + float64(unsafe.Sizeof(ebrActual))) / float64(mbr.Size) * 100
						porcentaje = porcentaje - numero
						extPorcentaje = extPorcentaje - numero
						var s string = strconv.FormatFloat(numero, 'f', 2, 64)
						contenido = contenido + "\\n" + s + string('%') + " del disco"
					}
				}

				if extPorcentaje > 0 {
					contenido += "|Libre"
					var s string = strconv.FormatFloat(extPorcentaje, 'f', 2, 64)
					contenido = contenido + "\\n" + s + string('%') + "del disco"
					porcentaje = porcentaje - extPorcentaje
				}

				contenido += "}}"
			}

		}
	}

	if porcentaje > 0 {
		contenido += "|Libre"
		var s string = strconv.FormatFloat(porcentaje, 'f', 2, 64)
		contenido = contenido + "\\n" + s + string('%') + "del disco"
	}

	contenido += "\"];\n"
	contenido += "}"

	archivo2, err := os.OpenFile(DotPath, os.O_RDWR, 0777)
	if err != nil {
		fmt.Println("Error: No se pudo abrir archivo")
		return
	}
	defer archivo2.Close()

	b := []byte(contenido)
	err2 := ioutil.WriteFile(DotPath, b, 0777)
	if err2 != nil {
		fmt.Println("Error: Error al escribir el archivo")
		return
	}

	cmd := exec.Command("dot", "-Tsvg", "-o", self.Path+".svg", DotPath)
	_, err3 := cmd.Output()
	if err3 != nil {
		fmt.Println("Error: Dot no pudo generar el svg: ", err)
	}
	fmt.Println("Reporte Generado: " + self.Path + ".svg")
	fmt.Print("\n")
}

//GENERAR EL REPORTE TREE
func (self Rep) repTree(DotPath string, mount ParticionMontada) {
	disco, err := os.OpenFile(mount.Path, os.O_RDWR, 0777)
	defer disco.Close()

	if err != nil {
		fmt.Println("Error: No se ha podido abrir el archivo")
		log.Fatal(err)
	}

	mbr := structs.GetMbr(mount.Path)
	part := structs.GetParticion(mount.Name, mbr, disco)
	superBloque := structs.LeerSuperBloque(disco, part.Start)
	n := structs.GetN(part.Size)

	contenido := "digraph H {\nrankdir=\"LR\"\n"

	//AQUI SE GENERA EL REPROTE

	apuntador := superBloque.Bm_inode_start
	for ap := int32(0); ap < n; ap++ {
		if structs.ExisteStructEnBM(disco, apuntador, n, int64(ap)) {
			inodo := structs.LeerInodo(disco, superBloque.Inode_start+int64(ap)*int64(superBloque.Inode_Size))
			contenido = contenido + self.graficarInodo(inodo, ap, superBloque, disco)
		}
	}

	contenido += "}"

	archivo2, err := os.OpenFile(DotPath, os.O_RDWR, 0777)
	if err != nil {
		fmt.Println("Error: No se pudo abrir archivo")
		return
	}
	defer archivo2.Close()

	b := []byte(contenido)
	err2 := ioutil.WriteFile(DotPath, b, 0777)
	if err2 != nil {
		fmt.Println("Error: Error al escribir el archivo")
		return
	}

	cmd := exec.Command("dot", "-Tsvg", "-o", self.Path+".svg", DotPath)
	_, err3 := cmd.Output()
	if err3 != nil {
		fmt.Println("Error: Dot no pudo generar el svg: ", err)
	}
	fmt.Println("Reporte Generado: " + self.Path + ".svg")
	fmt.Print("\n")
}

func (self Rep) graficarInodo(inodo structs.Inodo, numInodo int32, superBloque structs.SuperBloque, archivo *os.File) string {
	id := "I" + strconv.Itoa(int(numInodo))
	contenido := "I" + strconv.Itoa(int(numInodo)) + "[\nshape=plaintext\nlabel=<"
	contenido = contenido + "<table border='1' cellborder='1'>" //INICIO TABLA

	contenido = contenido + "<tr><td colspan=\"2\"> Inodo " + strconv.Itoa(int(numInodo)) + " </td></tr>"
	contenido = contenido + "<tr><td> Uid </td><td> " + strconv.Itoa(int(inodo.Uid)) + " </td></tr>"
	contenido = contenido + "<tr><td> Gid </td><td> " + strconv.Itoa(int(inodo.Gid)) + " </td></tr>"
	contenido = contenido + "<tr><td> Size </td><td> " + strconv.Itoa(int(inodo.Size)) + " </td></tr>"

	for i := 0; i < 16; i++ {
		contenido = contenido + "<tr><td> a" + strconv.Itoa(int(i)) + " </td><td port='a" + strconv.Itoa(int(i)) + "'> " + strconv.Itoa(int(inodo.Block[i])) + " </td></tr>"
	}

	contenido = contenido + "<tr><td> Type </td><td> " + strconv.Itoa(int(inodo.Type)) + " </td></tr>"
	contenido = contenido + "<tr><td> Perm </td><td> " + strconv.Itoa(int(inodo.Perm)) + " </td></tr>"

	contenido = contenido + "</table>\n>];\n" //FIN TABLA

	for i := 0; i < 16; i++ {
		if inodo.Block[i] != -1 {

			if inodo.Type == 0 {
				apuntador := superBloque.Block_start + int64(superBloque.Block_Size)*int64(inodo.Block[i])
				bloqueC := structs.LeerBloqueC(archivo, apuntador)
				contenido = contenido + self.graficarBloqueC(bloqueC, inodo.Block[i])
			} else if inodo.Type == 1 {
				apuntador := superBloque.Block_start + int64(superBloque.Block_Size)*int64(inodo.Block[i])
				bloqueA := structs.LeerBloqueA(archivo, apuntador)
				contenido = contenido + self.graficarBloqueA(bloqueA, inodo.Block[i])
			}

		}
	}

	//Graficar apuntadores
	for i := 0; i < 16; i++ {
		if inodo.Block[i] != -1 {
			contenido = contenido + id + ":a" + strconv.Itoa(i) + "->"
			contenido = contenido + "B" + strconv.Itoa(int(inodo.Block[i])) + "\n"
		}
	}

	return contenido
}

//GRAFICAR EL BLOQUE DE CARPETAS
func (self Rep) graficarBloqueC(bloque structs.BloqueCarpeta, numBloque int32) string {
	contenido := "B" + strconv.Itoa(int(numBloque)) + "[\nshape=plaintext\nlabel=<"
	contenido = contenido + "<table border='1' cellborder='1'>" //INICIO TABLA
	contenido = contenido + "<tr><td colspan=\"2\"> Bloque " + strconv.Itoa(int(numBloque)) + " </td></tr>"

	for i := 0; i < len(bloque.Contenido); i++ {
		cont := bloque.Contenido[i]
		contenido = contenido + "<tr><td> " + structs.GetNameBloqueString(cont.Name) + " </td><td port='a" + strconv.Itoa(int(i)) + "'> " + strconv.Itoa(int(cont.Apuntador)) + " </td></tr>"
	}

	contenido = contenido + "</table>\n>];\n" //FIN TABLA

	for i := 0; i < len(bloque.Contenido); i++ {
		cont := bloque.Contenido[i]
		name := structs.GetNameBloqueString(cont.Name)
		if name != "." && name != ".." && cont.Apuntador != -1 {
			contenido = contenido + "B" + strconv.Itoa(int(numBloque)) + ":a" + strconv.Itoa(i) + "->"
			contenido = contenido + "I" + strconv.Itoa(int(cont.Apuntador)) + "\n"
		}

	}

	return contenido
}

//GRAFICAR EL BLOQUE DE ARCHIVOS
func (self Rep) graficarBloqueA(bloque structs.BloqueArchivo, numBloque int32) string {
	contenido := "B" + strconv.Itoa(int(numBloque)) + "[\nshape=plaintext\nlabel=<"
	contenido = contenido + "<table border='1' cellborder='1'>" //INICIO TABLA
	contenido = contenido + "<tr><td colspan=\"1\"> Bloque " + strconv.Itoa(int(numBloque)) + " </td></tr>"

	txtBloque := ""
	numero := false
	for i := 0; i < 64; i++ {
		if numero {
			txtBloque = txtBloque + strconv.Itoa(int(bloque.Contenido[i]))
			continue
		}

		if bloque.Contenido[i] == 0 {
			txtBloque = txtBloque + strconv.Itoa(int(bloque.Contenido[i]))
			numero = true
		} else if bloque.Contenido[i] == '\n' {
			txtBloque = txtBloque + "\\n"
		} else {
			txtBloque = txtBloque + string(bloque.Contenido[i])
		}
	}
	contenido = contenido + "<tr><td>" + txtBloque + "</td></tr>"

	contenido = contenido + "</table>\n>];\n" //FIN TABLA
	return contenido
}

//Grafica
func (self Rep) repFile(mount ParticionMontada) {

}
