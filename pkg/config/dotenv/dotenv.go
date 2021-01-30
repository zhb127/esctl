package dotenv

import (
	"github.com/joeshaw/envdecode"
	"github.com/joho/godotenv"
)

func Load(paths ...string) error {
	// example: paths=[]string{".env.example", ".env"}
	if err := godotenv.Load(paths...); err != nil {
		return err
	}
	return nil
}

func Decode(targetPtr interface{}) error {
	if err := envdecode.StrictDecode(targetPtr); err != nil {
		return err
	}
	return nil
}
