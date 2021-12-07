package page

type Page struct {
	Size int
	Num  int
}

func Default() Page {
	return Page{
		Size: 10,
		Num:  0,
	}
}
