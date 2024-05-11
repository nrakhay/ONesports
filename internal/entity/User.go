package entity

type User struct {
	ID          int    `db:"id"`
	FirstName   string `db:"firstName"`
	LastName    string `db:"lastName"`
	Email       string `db:"email"`
	PhoneNumber string `db:"phoneNumber"`
}
