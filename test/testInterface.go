package test

import (
	"AutoMapper/test/bo"
	"AutoMapper/test/dto"
)

// SensorMapper
// @mapper
type SensorMapper interface {

	//BoToDto
	//@mapping(from="sensor.ComplicatedObject.BoSimple", to="sensor.ComplicatedObject.DtoSimple")
	BoToDto(sensor bo.Sensor) dto.Sensor
}

type (
	//TestMapper
	// @mapper
	TestMapper interface {
		//BoToDto
		//@mapping(from="sensor.ComplicatedObject.BoSimple", to="sensor.ComplicatedObject.DtoSimple")
		BoToDto(sensor bo.Sensor) dto.Sensor
	}
)
