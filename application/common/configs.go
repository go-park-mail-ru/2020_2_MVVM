package common

import "net/http"

type DBConfig struct {
	Host     string `yaml:"host"`
	Port     int    `yaml:"port"`
	User     string `yaml:"user"`
	Password string `yaml:"password"`
	Name     string `yaml:"name"`
}

type AuthCookieConfig struct {
	Key      string
	Path     string
	Domain   string
	MaxAge   int
	Secure   bool
	HttpOnly bool
	SameSite http.SameSite
}

type MicroserviceConfig struct {
	Host string `yaml:"host"`
	Port int    `yaml:"port"`
}
