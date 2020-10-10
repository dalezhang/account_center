package response

type ErrResponse struct {
	ErrMessage errMessage
	Status     int
}
type errMessage struct {
	Text string
	Code int
}
