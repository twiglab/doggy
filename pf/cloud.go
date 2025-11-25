package pf

import (
	"context"

	"github.com/twiglab/doggy/holo"
)

type DataHandler interface {
	HandleData(ctx context.Context, data UploadeData) error
}

type UploadeData struct {
	Common holo.Common
	Target holo.HumanMix
}
