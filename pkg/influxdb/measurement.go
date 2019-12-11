package influxdb

type Field struct {
	Name string
	Type string
}

type Measurement struct {
	Database string
	Name     string
	Tags     []string
	Fields   []Field
}
