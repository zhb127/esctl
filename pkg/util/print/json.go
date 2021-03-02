package print

import (
	"encoding/json"
	"fmt"
)

func JSON(v interface{}) {
	if res, err := json.Marshal(v); err != nil {
		panic(err)
	} else {
		fmt.Println(string(res))
	}
}
