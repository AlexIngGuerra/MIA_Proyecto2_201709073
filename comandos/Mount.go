package comandos

type Mount struct {
	Path string
	Name string
}

func NewMount() Mount {
	return Mount{Path: "", Name: ""}
}

func (self Mount) Ejecutar() {

}
