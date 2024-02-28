package keyservice

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"sync"
	"time"

	"github.com/febzey/ForestBot-Mainframe/database"
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
	//the key here will be encrypted.
	Key string

	//Email or some way for us to contact the key owner.
	OwnerEmail string

	//When the key was created.
	CreatedAt int64

	//millisecond timestamp of when this key was last used.
	UpdatedAt int64

	//The permissions the key has.
	//Permissions are created when the key is made, and is a constant
	Permissions APIPermissions

	// number of messages this api key can write per hour
	RateLimit int
}

// A structure for our api key service.
// perhaps we will use this package for some security features.
// if a part of our client needs to save, read, check an api key, we
// should do it through this service.
type APIKeyService struct {
	//this map contains each api key currently available.
	KeyArr []APIkey

	//a mutext to keep our map in sync.
	Mu *sync.Mutex

	// a struct to our database
	Db *database.Database
}

// a new service for our api keys and perhaps some more security features in the future?
func NewAPIKeyService(db *database.Database) *APIKeyService {
	return &APIKeyService{
		KeyArr: make([]APIkey, 0),
		Mu:     &sync.Mutex{},
		Db:     db,
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

	s.KeyArr = append(s.KeyArr, *apikey)

	return apikey

}

//A function to encrypt an api key using sha256 encryption
func (s *APIKeyService) EncryptAPIKey(plainTextKey string) (string, error) {
	hasher := sha256.New()
	_, err := hasher.Write([]byte(plainTextKey))
	if err != nil {
		return "", errors.New("Error encrypting plaintext.")
	}

	hashedBytes := hasher.Sum(nil)
	hashedString := hex.EncodeToString(hashedBytes)

	return hashedString, nil

}

// A function to verify and get our API key
func (s *APIKeyService) GetAndVerifyAPIKey(plainTextKey string) (*APIkey, bool) {
	s.Mu.Lock()
	defer s.Mu.Unlock()

	encryptedKey, err := s.EncryptAPIKey(plainTextKey)
	if err != nil {
		return nil, false
	}

	localKey, ok := s.GetAPIKeyLocal(encryptedKey)
	if !ok {
		return nil, false
	}

	return localKey, true
}

// Return an unmutable api key from our keymap
// keys will be stored in the map encrypted.
func (s *APIKeyService) GetAPIKeyLocal(encryptedKey string) (*APIkey, bool) {
	for _, apiKey := range s.KeyArr {
		if apiKey.Key == encryptedKey {
			copyAPIKey := apiKey
			return &copyAPIKey, true
		}
	}

	return nil, false
}

// A function that will be called in the NewKeyService func.
// will load api keys stored in our database, to our local map.
func (s *APIKeyService) retrieveKeysFromDatabaseAndLoadLocalMap() {

}

// Generate the api key.
func generateRandomKey() string {
	bytes := make([]byte, 16)
	if _, err := rand.Read(bytes); err != nil {
		panic(err)
	}
	return hex.EncodeToString(bytes)
}
