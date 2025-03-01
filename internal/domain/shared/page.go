package shared

type Page struct {
	Number int
	Size   int
}

func NewPage(number, size int) Page {
	if number < 1 {
		number = 1
	}
	if size < 1 {
		size = 10
	}
	return Page{
		Number: number,
		Size:   size,
	}
}

func (p Page) Offset() int {
	return (p.Number - 1) * p.Size
}

func (p Page) Limit() int {
	return p.Size
}
