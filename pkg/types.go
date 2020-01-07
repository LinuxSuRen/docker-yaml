package pkg

type Application struct {
	Name string
	Image string
	Restart bool
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