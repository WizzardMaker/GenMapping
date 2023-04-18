package SensorMapper

import (
	"AutoMapper/test/bo"
	"AutoMapper/test/dto"
	"AutoMapper/test/mapper/PropertyMapper"
)

//SensorMapper:

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

func BoToDtoComplex(sensor bo.Sensor, constOverride int) (target dto.Sensor) {
	target.Simple = sensor.Simple
	target.SimplePtr = sensor.SimplePtr
	target.SimpleString = sensor.SimpleString
	target.ComplicatedObject = PropertyMapper.BoToDto(sensor.ComplicatedObject)
	target.ComplicatedObject.DtoSimple = sensor.ComplicatedObject.BoSimple
	//target.ComplicatedObject.DoubleTest is not directly mapped
	//target.ComplicatedObject.Time is not directly mapped
	for i0 := range target.PropertyArray {
		target.PropertyArray[i0] = PropertyMapper.BoToDto(sensor.ComplicatedObject)
		target.PropertyArray[i0].DtoSimple = constOverride
		//target.PropertyArray[i0].DoubleTest is not directly mapped
		//target.PropertyArray[i0].Time is not directly mapped
	}
	return
}
