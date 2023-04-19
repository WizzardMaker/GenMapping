package bo

import (
	"GenMapping/test/double_name/bo"
	"time"
)

type Property struct {
	BoSimple   int
	DoubleTest bo.DoubleNameTest
	Time       time.Time
}
