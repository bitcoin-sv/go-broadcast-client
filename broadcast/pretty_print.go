package broadcast

import (
	"encoding/json"
	"fmt"
)

func PrettyPrint(name string, v interface{}) {
	vJson, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		fmt.Errorf("Error: %s", err.Error())
	}

	fmt.Printf("%s: %s\n", name, vJson)
}
