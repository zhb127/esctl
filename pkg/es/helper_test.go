package es

import (
	"esctl/pkg/config/dotenv"
	"esctl/test/data/pkg/log"
	"testing"

	"github.com/stretchr/testify/assert"
)

func MockHelperConfig() HelperConfig {
	config := HelperConfig{}
	if err := dotenv.Decode(&config); err != nil {
		panic(err)
	}
	return config
}

func Test(t *testing.T) {
	logHelper := log.MockHelper()
	helperConfig := MockHelperConfig()

	t.Run("NewHelper", func(t *testing.T) {
		helperInst, err := NewHelper(helperConfig, logHelper)
		if err != nil {
			t.Fatal(err)
		}

		assert.Nil(t, err)
		assert.NotNil(t, helperInst)
		assert.IsType(t, &helper{}, helperInst)
	})
}
