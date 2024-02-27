package utils

import (
	"crypto/rand"
	"encoding/hex"
	"sync"
	"time"
)

/******

Within this file we declare functions to generate client API keys,
we will either manually generate API keys through an HTTP endpoint.
Api Keys will be generated when a new client wishes to connect to our server
api keys will be saved in the database along with a way to contact the key owner,
when the key was generated, the permissions that the api key posseses.

******/

// Permissions for the api key.
type APIPermissions struct {
	//if read is true, then client is able to read every message
	//that is passed through the event proccessor.
	Read bool

	//if write permission is true, then our
	//event proccessor will accept client events.
	Write bool
}

// This is the structure all of our API keys will follow.
type APIkey struct {
	//The actual api key
	Key string

	//Email or some way for us to contact the key owner.
	OwnerEmail string

	//When the key was created.
	CreatedAt int64

	//The permissions the key has.
	Permissions APIPermissions
}

// A structure for our api key service.
// perhaps we will use this package for some security features.
// if a part of our client needs to save, read, check an api key, we
// should do it through this service.
type APIKeyService struct {
	//this map contains each api key currently available.
	KeyMap map[string]APIkey

	//a mutext to keep our map in sync.
	Mu sync.Mutex
}

// a new service for our api keys and perhaps some more security features in the future?
func NewAPIKeyService() *APIKeyService {
	return &APIKeyService{
		KeyMap: make(map[string]APIkey),
	}
}

// Func to generate our api keys
// when the key is generated, it will be automatically saved to the keymap
func (s *APIKeyService) NewApiKey(read, write bool, ownerEmail string) *APIkey {

	key := generateRandomKey()
	createdAt := time.Now().UnixNano() / int64(time.Millisecond)

	s.Mu.Lock()
	defer s.Mu.Unlock()

	apikey := &APIkey{
		Key:         key,
		OwnerEmail:  ownerEmail,
		CreatedAt:   createdAt,
		Permissions: APIPermissions{Read: read, Write: write},
	}

	s.KeyMap[key] = *apikey

	return apikey

}

// Return an unmutable api key from our keymap
func (s *APIKeyService) GetAPIKey(key string) (*APIkey, bool) {
	s.Mu.Lock()
	defer s.Mu.Unlock()

	apiKey, ok := s.KeyMap[key]
	copyAPIKey := apiKey
	return &copyAPIKey, ok
}

// Generate the api key.
func generateRandomKey() string {
	bytes := make([]byte, 16)
	if _, err := rand.Read(bytes); err != nil {
		panic(err)
	}
	return hex.EncodeToString(bytes)
}
