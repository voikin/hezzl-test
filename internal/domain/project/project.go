package project

type Project struct {
	ID   int    `db:"id"`
	Name string `db:"name"`
}
