package test

import (
	"GenMapping/test/bo"
	bo2 "GenMapping/test/double_name/bo"
	"GenMapping/test/dto"
)

// SensorMapper
// @mapper
type SensorMapper interface {

	//BoToDto
	//@translate(from="sensor.ComplicatedObject.BoSimple", to="ComplicatedObject.DtoSimple")
	BoToDto(sensor bo.Sensor) dto.Sensor

	//BoToDtoComplex
	//@translate(from="sensor.ComplicatedObject.BoSimple", to="ComplicatedObject.DtoSimple")
	//@translate(from="constOverride", to="PropertyArray.DtoSimple")
	BoToDtoComplex(sensor bo.Sensor, constOverride int) dto.Sensor
}

// PropertyMapper
// @mapper
type PropertyMapper interface {

	//BoToDto
	//@translate(from="property.BoSimple", to="DtoSimple")
	BoToDto(property bo.Property) dto.Property

	//BoToDtoFixed
	//@translate(from="property.BoSimple", to="DtoSimple")
	//@translate(from="test", to="DoubleTest")
	BoToDtoFixed(property bo.Property, test bo2.DoubleNameTest) dto.Property
}

type (
	//TestMapper
	// @mapper
	TestMapper interface {
		//BoToDto
		//@translate(from="sensor.ComplicatedObject.BoSimple", to="ComplicatedObject.DtoSimple")
		BoToDto(sensor bo.Sensor) dto.Sensor
	}
)
