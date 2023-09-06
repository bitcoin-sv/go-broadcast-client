package broadcast

import (
	"encoding/json"
	"fmt"
)

// PrettyPrint prints JSON in a friendly format
func PrettyPrint(name string, v interface{}) {
	vJson, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		fmt.Printf("Error: %+v\n", err)
	}

	fmt.Printf("%s: %s\n", name, vJson)
}
