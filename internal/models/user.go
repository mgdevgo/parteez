package models

type User interface {
}

type Users interface {
	Add(user User) error
	Find(id string) (User, error)
}

type DJ interface {
	User
	Find(id string) (int, error)
}

type SQLUser interface {
	User
	Save() error
}

type SQLUsers interface {
	Add(user SQLUser) error
}

// type User interface {
// 	AddToShoppingCart(product Product)
// 	//
// 	// some additional methods
// 	//
// }
//
// type LoggedInUser interface {
// 	User
// 	SaveEventToLibrary(eventID int) error
// 	//
// 	// some additional methods
// 	//
// }
//
// type PremiumUser interface {
// 	LoggedInUser
// 	HasDiscountFor(product Product) bool
// 	//
// 	// some additional methods
// 	//
// }
