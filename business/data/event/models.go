package event

type Event struct {
	ID        string `db:"id" json:"id"`
	Title     string `db:"title" json:"title"`
	Location  string `db:"location" json:"location"`
	EventDate string `db:"event_date" json:"when"`
}
