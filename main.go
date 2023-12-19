package main

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/dp-152/gogetsecret/provider/mapp"
)

type SecretRequest struct {
	Secrets []string
	Version string
}

type SecretEntry struct {
	Value *string
	Error *string
}

func main() {
	input, err := io.ReadAll(os.Stdin)

	if err != nil {
		fmt.Fprint(os.Stderr, fmt.Errorf("cannot read input: %w", err))
	}

	req := new(SecretRequest)
	err = json.Unmarshal(input, req)

	if err != nil {
		fmt.Fprint(os.Stderr, fmt.Errorf("cannot unmarshall input: %w", err))
	}

	var res []byte
	switch req.Version {
	case "1.0":
		res, err = handleRequestV1_0(req)
	default:
		err = fmt.Errorf("unknown payload version %s", req.Version)
	}

	if err != nil {
		fmt.Fprint(os.Stderr, fmt.Errorf("cannot fetch secrets: %w", err))
	}

	fmt.Println(res)
}

func handleRequestV1_0(req *SecretRequest) (res []byte, fatal error) {
	secretResults := make(map[string]*SecretEntry)
	for _, secret := range req.Secrets {
		var (
			value string
			err   error
		)

		src, ident, ok := strings.Cut(secret, ":")

		if !ok {
			src = "map"
			ident = secret
		}

		switch src {
		case "map":
			value, err, fatal = mapp.GetSecret(ident)
		default:
			err = fmt.Errorf("no secret provider found for source %s", src)
		}

		if fatal != nil {
			return
		}

		result := &SecretEntry{}
		if err != nil {
			errstr := err.Error()
			result.Error = &errstr
		} else {
			result.Value = &value
		}
	}

	res, fatal = json.Marshal(secretResults)

	return
}
