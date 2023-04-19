package TestMapper

import (
	"GenMapping/test/bo"
	"GenMapping/test/dto"
	"GenMapping/test/mapper/PropertyMapper"
)

//TestMapper:

func BoToDto(sensor bo.Sensor) (target dto.Sensor) {
	target.Simple = sensor.Simple
	target.SimplePtr = sensor.SimplePtr
	target.SimpleString = sensor.SimpleString
	target.ComplicatedObject = PropertyMapper.BoToDto(sensor.ComplicatedObject)
	target.ComplicatedObject.DtoSimple = sensor.ComplicatedObject.BoSimple
	//target.ComplicatedObject.DoubleTest is not directly mapped
	//target.ComplicatedObject.Time is not directly mapped
	for i0 := range target.PropertyArray {
		target.PropertyArray[i0] = PropertyMapper.BoToDto(sensor.ComplicatedObject)
		//target.PropertyArray[i0].DtoSimple is not directly mapped
		//target.PropertyArray[i0].DoubleTest is not directly mapped
		//target.PropertyArray[i0].Time is not directly mapped
	}
	return
}
