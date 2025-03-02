package toucher

type Toucher interface {
	Touch(pin uint8) error
}
