package dto

type Sensor struct {
	Simple            int
	SimplePtr         *int
	SimpleString      string
	ComplicatedObject Property
	PropertyArray     []Property
}
