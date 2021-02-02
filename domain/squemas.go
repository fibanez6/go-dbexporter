package domain

type Device struct {
	Name          string `db:name`
	LastIpAddress string `db:lastipaddress`
}

type Monitor struct {
	SerialNumber string  `db:serialnumber`
	Resolution   float64 `db:resolution`
}
