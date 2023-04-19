package dto

import (
	"GenMapping/test/double_name/bo"
	"time"
)

type Property struct {
	DtoSimple  int
	DoubleTest bo.DoubleNameTest
	Time       time.Time
}
