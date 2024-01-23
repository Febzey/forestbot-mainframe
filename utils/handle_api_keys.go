package utils

import "os"

type APIPermissions struct {
	Read  bool
	Write bool
}

func CheckApiKey(api_key string) APIPermissions {
	apiKeyType := determineAPIKeyType(api_key)

	switch apiKeyType {
	case "read_write":
		return APIPermissions{Read: true, Write: true}
	case "read_only":
		return APIPermissions{Read: true, Write: false}
	default:
		return APIPermissions{Read: false, Write: false}
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
