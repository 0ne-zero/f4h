package public_struct

type RequestBasicInformation struct {
	IP     string
	Path   string
	Method string
}

type ProductForCartItems struct {
	ID        int
	Name      string
	Price     float64
	Quantity  int
	ImagePath string
}
