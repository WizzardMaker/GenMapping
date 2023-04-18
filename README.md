# Go - AutoMapper
This project is a struct to struct mapping code generator for Go.

## Concept
AutoMapper analyzes a Go module and searches for interfaces which describe desired mapping operations. 
AutoMapper generates mapper based on those interfaces. The generation can be configured with tag documentation

### Tags
Tags are comment annotations which instruct AutoMapper how to generate the mapper functions.

Documentation for Tags can be found in `documentation/Tags.mnd`

### Syntax
A mapper is a simple Go interface annotated with the `@mapper` tag. This tag tells AutoMapper, that this interface describes mapping functions.

Mapping functions have the following rules:
- Only 1 return is allowed _(error handling is planned in a future update)_
- It has to be exported

## Examples

Structures used in the examples can be found under `test/`

Given a mapper interface like this:
```go
// SensorMapper
// @mapper
type SensorMapper interface {

	//BoToDto
	//@translate(from="sensor.ComplicatedObject.BoSimple", to="ComplicatedObject.DtoSimple")
	BoToDto(sensor bo.Sensor) dto.Sensor
}
```

AutoMapper will generate a mapper function like this:
```go
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
```

## To Do:
The following things are still not done (*) or are yet untested (°)
- Pointer Type check*
  - Mapping to pointers is still missing some type checks and eventual (de-)referencing between pointer and non pointer
- Maps*
  - `map[X]Y` mapping is not yet supported
- Ignored Fields*
  - Some fields should not be touched during mapping (Ids from database objects for example), thus needing a "ignore" tag for such fields
- Tag "inheritance" or groups*
  - Some mappers can share tags - allowing for referencable groups minimizes possible code duplications  
- Custom mapper file paths°
  - These should work, but some edge cases could break the resulting output