package event

type Ticket struct {
	Title       string
	Price       int
	Description string
}

func NewTicket(title string, price int, description string) Ticket {
	return Ticket{
		Title:       title,
		Price:       price,
		Description: description,
	}
}
