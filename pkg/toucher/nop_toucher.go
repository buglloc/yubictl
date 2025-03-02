package toucher

import "errors"

var _ Toucher

type NopToucher struct{}

func NewNopToucher() *NopToucher {
	return &NopToucher{}
}

func (t *NopToucher) Touch(_ uint8) error {
	return errors.New("not implemented")
}
