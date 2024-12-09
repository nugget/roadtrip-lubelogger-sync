package roadtrip

// Road Trip Data File version 1500,en

// FUEL RECORDS
type Fuel struct {
	Odometer     int     `csv:"Odometer (mi)"`
	TripDistance int     `csv:"Trip Distance,omitempty"`
	Date         string  `csv:"Date"`
	FillAmount   float64 `csv:"Fill Amount,omitempty"`
	FillUnits    string  `csv:"Fill Units"`
	PricePerUnit float64 `csv:"Price per Unit,omitempty"`
	TotalPrice   float64 `csv:"Total Price,omitempty"`
	PartialFill  bool    `csv:"Partial Fill,omitempty"`
	MPG          float64 `csv:"MPG,omitempty"`
	Note         string  `csv:"Note"`
	Octane       string  `csv:"Octane"`
	Location     string  `csv:"Location"`
	Payment      string  `csv:"Payment"`
	Conditions   string  `csv:"Conditions"`
	Reset        string  `csv:"Reset"`
	Categories   string  `csv:"Categories"`
	Flags        string  `csv:"Flags"`
	CurrencyCode int     `csv:"Currency Code,omitempty"`
	CurrencyRate int     `csv:"Currency Rate,omitempty"`
	Latitude     float64 `csv:"Latitude,omitempty"`
	Longitude    float64 `csv:"Longitude,omitempty"`
	ID           int     `csv:"ID,omitempty"`
	FuelEconomy  string  `csv:"Trip Comp Fuel Economy"`
	AvgSpeed     string  `csv:"Trip Comp Avg. Speed"`
	Temperature  float64 `csv:"Trip Comp Temperature,omitempty"`
	DriveTime    string  `csv:"Trip Comp Drive Time"`
	TankNumber   int     `csv:"Tank Number,omitempty"`
}

// MAINTENANCE RECORDS
type Maintenance struct {
	Description          string  `csv:"Description"`
	Date                 string  `csv:"Date"`
	Odometer             int     `csv:"Odometer (mi.),omitempty"`
	Cost                 float64 `csv:"Cost,omitempty"`
	Note                 string  `csv:"Note"`
	Location             string  `csv:"Location"`
	Type                 string  `csv:"Type"`
	Subtype              string  `csv:"Subtype"`
	Payment              string  `csv:"Payment"`
	Categories           string  `csv:"Categories"`
	ReminderInterval     string  `csv:"Reminder Interval"`
	ReminderDistance     string  `csv:"Reminder Distance"`
	Flags                string  `csv:"Flags"`
	CurrencyCode         int     `csv:"Currency Code,omitempty"`
	CurrencyRate         int     `csv:"Currency Rate,omitempty"`
	Latitude             float64 `csv:"Latitude,omitempty"`
	Longitude            float64 `csv:"Longitude,omitempty"`
	ID                   int     `csv:"ID,omitempty"`
	NotificationInterval string  `csv:"Notification Interval"`
	NotificationDistance string  `csv:"Notification Distance"`
}

// ROAD TRIPS
type RoadTrip struct {
	Name          string `csv:"Name"`
	StartDate     string `csv:"Start Date"`
	StartOdometer int    `csv:"Start Odometer (mi.),omitempty"`
	EndDate       string `csv:"End Date"`
	EndOdometer   int    `csv:"End Odometer,omitempty"`
	Note          string `csv:"Note"`
	Distance      int    `csv:"Distance,omitempty"`
	ID            int    `csv:"ID,omitempty"`
	Type          string `csv:"Type"`
	Categories    string `csv:"Categories"`
	Flags         string `csv:"Flags"`
}

// VEHICLE
type Vehicle struct {
	Name                string `csv:"Name"`
	Odometer            string `csv:"Odometer"`
	Units               string `csv:"Units"`
	Notes               string `csv:"Notes"`
	TankCapacity        string `csv:"Tank Capacity"`
	Tank1Units          string `csv:"Tank Units"`
	HomeCurrency        string `csv:"Home Currency"`
	Flags               string `csv:"Flags"`
	IconID              string `csv:"IconID"`
	FuelUnits           string `csv:"FuelUnits"`
	TripCompUnits       string `csv:"TripComp Units"`
	TripCompSpeed       string `csv:"TripComp Speed"`
	TripCompTemperature string `csv:"TripComp Temperature"`
	TripCompTimeEnabled string `csv:"TripComp Time Enabled"`
	OdometerShift       string `csv:"Odometer Shift"`
	Tank1Type           string `csv:"Tank 1 Type"`
	Tank2Type           string `csv:"Tank 2 Type"`
	Tank2Units          string `csv:"Tank 2 Units"`
}

// TIRE LOG
type Tire struct {
	Name           string `csv:"Name"`
	StartDate      string `csv:"Start Date"`
	StartOdometer  int    `csv:"Start Odometer (mi.),omitempty"`
	Size           string `csv:"Size"`
	SizeCorrection string `csv:"Size Correction"`
	Distance       int    `csv:"Distance,omitempty"`
	Age            string `csv:"Age"`
	Note           string `csv:"Note"`
	Flags          string `csv:"Flags"`
	ID             int    `csv:"ID,omitempty"`
	ParentID       int    `csv:"ParentID,omitempty"`
}

// VALUATIONS
type Valuation struct {
	Type     string `csv:"Type"`
	Date     string `csv:"Date"`
	Odometer int    `csv:"Odometer,omitempty"`
	Price    string `csv:"Price"`
	Notes    string `csv:"Notes"`
	Flags    string `csv:"Flags"`
}
