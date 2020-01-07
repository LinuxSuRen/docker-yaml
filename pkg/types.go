package pkg

type Application struct {
	Name string
	Image string
	Restart string
	Exposes []Expose
	Volumes []Volume
}

type Expose struct {
	Host int
	Container int
}

type Volume struct {
	Host string
	Container string
}