package test

import (
	"AutoMapper/test/bo"
	"AutoMapper/test/dto"
)

// SensorMapper
// @mapper
type SensorMapper interface {

	//BoToDto
	//@translate(from="sensor.ComplicatedObject.BoSimple", to="target.ComplicatedObject.DtoSimple")
	BoToDto(sensor bo.Sensor) dto.Sensor
}

// PropertyMapper
// @mapper
type PropertyMapper interface {

	//BoToDto
	//@translate(from="property.BoSimple", to="target.DtoSimple")
	BoToDto(property bo.Property) dto.Property
}

type (
	//TestMapper
	// @mapper
	TestMapper interface {
		//BoToDto
		//@translate(from="sensor.ComplicatedObject.BoSimple", to="target.ComplicatedObject.DtoSimple")
		BoToDto(sensor bo.Sensor) dto.Sensor
	}
)
