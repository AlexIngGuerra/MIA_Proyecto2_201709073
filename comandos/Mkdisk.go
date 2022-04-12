package comandos

type Mkdisk struct {
	Size int
	Fit  string
	Unit string
	Path string
}

func NewMkdisk() Mkdisk {
	return Mkdisk{Size: 0, Fit: "FF", Unit: "M", Path: ""}
}

func (self Mkdisk) Ejecutar() {

}
