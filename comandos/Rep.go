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

func (self Rep) Ejecutar() {
	if self.tieneErrores() {
		return
	}

	mount := self.getMount(self.Id)
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
		}

		cmd := exec.Command("dot", "-Tsvg", "-o", self.Path+".svg", DotPath)
		_, err := cmd.Output()
		if err != nil {
			fmt.Println("Error: Dot no pudo generar el svg: ", err)
		}

	}
}

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

func (self Rep) getMount(Id string) ParticionMontada {
	for i := 0; i < len(Montados); i++ {
		if Montados[i].Id == Id {
			return Montados[i]
		}
	}
	return ParticionMontada{}
}

func (self Rep) repDisk(DotPath string, mount ParticionMontada) {
	mbr := structs.GetMbr(mount.Path)
	archivo, err := os.OpenFile(mount.Path, os.O_RDWR, 0777)
	defer archivo.Close()

	if err != nil {
		fmt.Println("Error: No se ha podido abrir el archivo")
		log.Fatal(err)
	}

	if mbr.Tamano <= 0 {
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
				numero := float64(particion.Size) / float64(mbr.Tamano) * 100
				porcentaje = porcentaje - numero
				var s string = strconv.FormatFloat(numero, 'f', 2, 64)
				contenido = contenido + "\\n" + s + string('%') + " del disco"
			}

			if particion.Type == 'E' {
				contenido += "|{Particion Extendida|{"
				extPorcentaje := float64(particion.Size) / float64(mbr.Tamano) * 100

				apuntador := mbr.Particion[i].Start
				ebrActual := structs.GetEbr(archivo, apuntador)
				if ebrActual.Size != 0 {
					contenido = contenido + "EBR| Particion Logica"
					numero := (float64(ebrActual.Size) + float64(unsafe.Sizeof(ebrActual))) / float64(mbr.Tamano) * 100
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
						numero := (float64(ebrActual.Size) + float64(unsafe.Sizeof(ebrActual))) / float64(mbr.Tamano) * 100
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
}
