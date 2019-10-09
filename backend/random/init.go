package random

import (
	"math/rand"
	"time"
)

func init() {
	rand.Seed(time.Now().UTC().Unix())
}
