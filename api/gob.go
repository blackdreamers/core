package api

import "encoding/gob"

func GobModels(models ...interface{}) {
	for _, m := range models {
		gob.Register(m)
	}
}
