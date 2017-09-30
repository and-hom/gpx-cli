package v2

import "fmt"

type Mode struct {
	preserveTracks   bool
	preserveSegments bool
}

func ParseMode(s string) (Mode, error) {
	switch s {
	case "n":
		return Mode{false, false}, nil
	case "st":
		return Mode{true, true}, nil
	case "ts":
		return Mode{true, true}, nil
	case "t":
		return Mode{true, false}, nil
	case "s":
		return Mode{false, true}, nil
	default:
		return Mode{}, fmt.Errorf("Unknown mode %s", s)
	}
}

const FILTER_MODE_FLAG string = "mode"
const FILTER_MODEL_FLAG string = "model"
const FILTER_MODEL_PARAMS_FLAG string = "model-params"

