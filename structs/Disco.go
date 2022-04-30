package structs

type Mbr struct {
	Size           int32
	Fecha_Creacion [20]uint8
	Dks_Signature  int64
	Fit            uint8
	Particion      [4]Partition
}

type Partition struct {
	Status uint8
	Type   uint8
	Fit    uint8
	Start  int64 //Apuntador
	Size   int32
	Name   [20]uint8
}

type Ebr struct {
	Status uint8
	Fit    uint8
	Start  int64 //Apuntador
	Size   int32
	Next   int64 //Apuntador
	Name   [20]uint8
}

type InfoPart struct {
	Start int64 //Apuntador
	Size  int32
}
