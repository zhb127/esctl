package log

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test(t *testing.T) {
	helperConfig, err := DefaultHelperConfig()
	if err != nil {
		t.Fatal(err)
	}

	t.Run("helper.mergeFields", func(t *testing.T) {
		helperInst, err := NewHelper(helperConfig)
		if err != nil {
			t.Fatal(err)
		}

		helperInst.withGlobalFields = map[string]interface{}{
			"a": "a",
			"b": 2,
		}

		helperInst.withFields = map[string]interface{}{
			"a": map[string]interface{}{
				"aa": 11,
			},
			"c": "c",
		}

		customFields := map[string]interface{}{
			"d": "4",
			"b": "b",
		}

		expected := map[string]interface{}{
			"a": map[string]interface{}{
				"aa": 11,
			},
			"b": "b",
			"c": "c",
			"d": "4",
		}

		actual := helperInst.mergeFileds(customFields)

		assert.Equal(t, expected["a"], actual["a"])
		assert.Equal(t, expected["b"], actual["b"])
		assert.Equal(t, expected["c"], actual["c"])
		assert.Equal(t, expected["d"], actual["d"])
	})

	t.Run("helper.SetWithField", func(t *testing.T) {
		helperInst, err := NewHelper(helperConfig)
		if err != nil {
			t.Fatal(err)
		}

		helperInst.SetWithField("demo", "123")

		expected := map[string]interface{}{
			"demo": "123",
		}

		actual := helperInst.withFields

		assert.Equal(t, expected, actual)
	})

	t.Run("helper.SetWithGlobalField", func(t *testing.T) {
		helperInst, err := NewHelper(helperConfig)
		if err != nil {
			t.Fatal(err)
		}

		helperInst.SetWithGlobalField("demo", "123")

		helperChildInst := helperInst.NewChild()

		expectedPtr := fmt.Sprintf("%p", helperInst.withGlobalFields)
		actualPtr := fmt.Sprintf("%p", helperChildInst.withGlobalFields)
		assert.Equal(t, expectedPtr, actualPtr)

		expectedVal := map[string]interface{}{
			"demo": "123",
		}
		actualVal := helperChildInst.withGlobalFields
		assert.Equal(t, expectedVal, actualVal)
	})

	t.Run("NewHelper", func(t *testing.T) {
		helperInst, err := NewHelper(helperConfig)
		if err != nil {
			t.Fatal(err)
		}

		assert.Nil(t, err)
		assert.NotNil(t, helperInst)
		assert.IsType(t, &helper{}, helperInst)

		t.Run("helper.LogLevel", func(t *testing.T) {
			logLevel := helperInst.LogLevel()
			assert.Equal(t, LOG_LEVEL_DEBUG, logLevel)
		})

		t.Run("helper.LogLevelNum", func(t *testing.T) {
			logLevelNum := helperInst.LogLevelNum()
			assert.Equal(t, LOG_LEVEL_DEBUG_NUM, logLevelNum)
		})

		t.Run("helper.NewChild", func(t *testing.T) {
			childHelperInst := helperInst.NewChild()
			assert.NotSame(t, helperInst, childHelperInst)
		})
	})

}
