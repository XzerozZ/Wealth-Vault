package line

type LineClient interface {
	SendTextMessage(to string, text string) error
}
