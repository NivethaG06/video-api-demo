package configloader

type Configurations struct {
	LocalServer struct {
		Host string `yaml:"host"`
		Port int    `yaml:"port"`
	} `yaml:"localserver"`
	DB struct {
		Host     string `yaml:"host"`
		Port     int    `yaml:"port"`
		Name     string `yaml:"db_name"`
		User     string `yaml:"db_user"`
		Password string `yaml:"db_password"`
	} `yaml:"db"`
	Redis struct {
		Host string `yaml:"host"`
		Port int    `yaml:"port"`
	} `yaml:"redis"`
}

type Localserver struct {
	host string
	port int
}

type Db struct {
	host     string
	port     int
	Name     string `yaml:"db_name"`
	User     string `yaml:"db_user"`
	Password string `yaml:"db_password"`
}

type Redis struct {
	host string
	port int
}
