package server

type User struct {
	Username string `yaml:"username" bson:"username"`
	Password string `yaml:"password" bson:"password"`
}
