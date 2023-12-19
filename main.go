package main

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
)

type SecretRequest struct {
	Secrets []string
	Version string
}

func main() {
	input, err := io.ReadAll(os.Stdin)

	if err != nil {
		panic(fmt.Errorf("cannot read input: %w", err))
	}

	req := new(SecretRequest)
	err = json.Unmarshal(input, req)

	if err != nil {
		panic(fmt.Errorf("cannot unmarshall input: %w", err))
	}

	var res string
	switch req.Version {
	default:
		err = fmt.Errorf("unknown payload version %s", req.Version)
	}

	if err != nil {
		panic(fmt.Errorf("cannot fetch secrets: %w", err))
	}

	fmt.Println(res)
}
