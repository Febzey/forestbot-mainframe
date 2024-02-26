package utils

import (
	"errors"
	"os"

	"github.com/febzey/ForestBot-Mainframe/types"
)

func CheckApiKey(api_key string) (types.APIPermissions, error) {
	apiKeyType := determineAPIKeyType(api_key)

	switch apiKeyType {
	case "read_write":
		return types.APIPermissions{Read: true, Write: true}, nil
	case "read_only":
		return types.APIPermissions{Read: true, Write: false}, nil
	default:
		return types.APIPermissions{Read: false, Write: false}, errors.New("invalid api key")
	}

}

func determineAPIKeyType(apiKey string) string {

	READ_WRITE_KEY := os.Getenv("READ_WRITE_KEY")
	READ_ONLY_KEY := os.Getenv("READ_ONLY_KEY")

	switch apiKey {
	case READ_WRITE_KEY:
		return "read_write"
	case READ_ONLY_KEY:
		return "read_only"
	default:
		return "invalid"
	}

}
