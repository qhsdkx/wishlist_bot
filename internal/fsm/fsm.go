package fsm

type StateStore interface {
	Set(chatID int64, state string) error
	Get(chatID int64) (string, error)
	Delete(chatID int64) error
}
