package dto

type Secret struct {
	User     string `json:"user"`
	Password string `json:"password"`
	Ip       string `json:"ip"`
	Port     string `json:"port"`
	Db       string `json:"db"`
	Server   string `json:"server"`
}
