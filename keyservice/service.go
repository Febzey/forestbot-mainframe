package keyservice

import (
	"crypto/rand"
	"crypto/sha256"
	"database/sql"
	"encoding/hex"
	"errors"
	"fmt"
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

	//bot-client, client, etc.
	TokenType string
}

// A structure for our api key service.
// perhaps we will use this package for some security features.
// if a part of our client needs to save, read, check an api key, we
// should do it through this service.
type APIKeyService struct {
	//Cache for used keys
	keyArr map[string]APIkey

	//mutex to keep us in sync
	mu *sync.Mutex

	// a struct to our database
	Db *sql.DB
}

// a new service for our api keys and perhaps some more security features in the future?
func NewAPIKeyService(db *sql.DB) *APIKeyService {
	return &APIKeyService{
		keyArr: make(map[string]APIkey),
		mu:     &sync.Mutex{},
		Db:     db,
	}
}

// Func to generate our api keys
// when the key is generated, it will be automatically saved to the database
// we return the api key in plaintext to give to the client. after that its always encrypted.
func (s *APIKeyService) NewApiKey(read, write bool, ownerEmail string, rateLimit int, TokenType string) (string, error) {

	key := generateRandomKey()
	keyEncryped, _ := s.EncryptAPIKey(key)
	createdAt := time.Now().UnixNano() / int64(time.Millisecond)

	apikey := APIkey{
		Key:         keyEncryped,
		OwnerEmail:  ownerEmail,
		CreatedAt:   createdAt,
		UpdatedAt:   createdAt,
		Permissions: APIPermissions{Read: read, Write: write},
		RateLimit:   rateLimit,
		TokenType:   TokenType,
	}

	err := s.saveKeyToDatabase(apikey)
	if err != nil {
		return "", err
	}

	return key, nil
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
// Will check our local key map if not present will search the database.
func (s *APIKeyService) GetAndVerifyAPIKey(plainTextKey string) (APIkey, bool) {
	s.mu.Lock()
	defer s.mu.Unlock()

	encryptedKey, err := s.EncryptAPIKey(plainTextKey)
	if err != nil {
		return APIkey{}, false
	}

	var key APIkey

	cachedKey, ok := s.keyArr[encryptedKey]
	if !ok {
		key, err := s.retrieveKeyFromDatabase(encryptedKey)
		if err != nil {
			return APIkey{}, false
		}

		s.keyArr[key.Key] = key

		return key, true
	}

	key = cachedKey

	return key, true
}

// A function to get a key from the database using the encrypted key.
func (s *APIKeyService) retrieveKeyFromDatabase(encryptedKey string) (APIkey, error) {

	query := `
        SELECT Api_key, OwnerEmail, CreatedAt, UpdatedAt, ReadPermission, WritePermission, RateLimit, TokenType
        FROM api_keys
        WHERE Api_key = ?;
    `

	var apiKey struct {
		Api_key         string `json:"Api_key"`
		OwnerEmail      string `json:"OwnerEmail"`
		CreatedAt       int64  `json:"CreatedAt"`
		UpdatedAt       int64  `json:"UpdatedAt"`
		ReadPermission  int    `json:"ReadPermission"`
		WritePermission int    `json:"WritePermission"`
		RateLimit       int    `json:"RateLimit"`
		TokenType       string `json:"TokenType"`
	}

	err := s.Db.QueryRow(query, encryptedKey).Scan(
		&apiKey.Api_key,
		&apiKey.OwnerEmail,
		&apiKey.CreatedAt,
		&apiKey.UpdatedAt,
		&apiKey.ReadPermission,
		&apiKey.WritePermission,
		&apiKey.RateLimit,
		&apiKey.TokenType,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return APIkey{}, fmt.Errorf("key not found.")
		}

		return APIkey{}, fmt.Errorf("error getting key: %w", err)
	}

	var read bool = false
	var write bool = false

	if apiKey.ReadPermission == 1 {
		read = true
	}
	if apiKey.WritePermission == 1 {
		write = true
	}

	builtApiKey := APIkey{
		Key:        apiKey.Api_key,
		OwnerEmail: apiKey.OwnerEmail,
		CreatedAt:  apiKey.CreatedAt,
		UpdatedAt:  apiKey.UpdatedAt,
		Permissions: APIPermissions{
			Read:  read,
			Write: write,
		},
		RateLimit: apiKey.RateLimit,
		TokenType: apiKey.TokenType,
	}

	return builtApiKey, nil
}

// A function to save a newly created key to the database
func (s *APIKeyService) saveKeyToDatabase(key APIkey) error {
	currentTime := time.Now()
	nanoseconds := currentTime.UnixNano()
	milliseconds := nanoseconds / int64(time.Millisecond)

	tableQuery := `
	CREATE TABLE IF NOT EXISTS api_keys (
		Api_key VARCHAR(255) NOT NULL,
		OwnerEmail VARCHAR(255) NOT NULL,
		CreatedAt BIGINT NOT NULL,
		UpdatedAt BIGINT NOT NULL,
		ReadPermission TINYINT NOT NULL,
		WritePermission TINYINT NOT NULL,
		RateLimit INT NOT NULL,
		TokenType VARCHAR(255) NOT NULL,
		PRIMARY KEY (Api_key)
	  );
	`

	insertQuery := `
	INSERT INTO api_keys (Api_key, OwnerEmail, CreatedAt, UpdatedAt, ReadPermission, WritePermission, RateLimit, TokenType)
	VALUES (?, ?, ?, ?, ?, ?, ?, ?);
	`

	if _, err := s.Db.Exec(tableQuery); err != nil {
		return fmt.Errorf("error checking/creating table: %w", err)
	}

	if _, err := s.Db.Exec(insertQuery, key.Key, key.OwnerEmail, milliseconds, milliseconds, key.Permissions.Read, key.Permissions.Write, key.RateLimit, key.TokenType); err != nil {
		return fmt.Errorf("error inserting API key: %w", err)
	}

	return nil
}

// Generate the api key.
func generateRandomKey() string {
	bytes := make([]byte, 16)
	if _, err := rand.Read(bytes); err != nil {
		panic(err)
	}
	return hex.EncodeToString(bytes)
}
