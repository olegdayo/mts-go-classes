package main

var conf *Config

func init() {
	var err error
	conf, err = InitConfig()
	if err != nil {
		panic(err)
	}
}

func main() {
}
