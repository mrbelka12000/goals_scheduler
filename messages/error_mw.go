package messages

type (
	errorSender interface {
		SendError(error)
	}
)
