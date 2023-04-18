package PropertyMapper

import (
	"AutoMapper/test/bo"
	bo2 "AutoMapper/test/double_name/bo"
	"AutoMapper/test/dto"
)

//PropertyMapper:

func BoToDto(property bo.Property) (target dto.Property) {
	target.DtoSimple = property.BoSimple
	target.DoubleTest = property.DoubleTest
	target.Time = property.Time
	return
}

func BoToDtoFixed(property bo.Property, test bo2.DoubleNameTest) (target dto.Property) {
	target.DtoSimple = property.BoSimple
	target.DoubleTest = test
	target.Time = property.Time
	return
}
