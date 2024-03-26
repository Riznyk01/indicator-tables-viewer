package config

type Config struct {
	Host     string
	Port     string
	Path     string
	Username string
	DBName   string
	Password string
}

func NewConfig() *Config {
	return &Config{
		//Host: "localhost",
		//Host: "127.0.0.1",
		Host: "192.168.0.1",
		Port: "3050",
		//Path: "home/user/db",
		//Path: "C:/childrend",
		Path:     "D:/s/Temp/temp",
		Username: "sysdba",
		DBName:   "MEDSTAT.GDB",
		Password: "1",
		//Password: "masterkey",
	}
}
