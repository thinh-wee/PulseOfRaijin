package keystore

import "crypto/sha256"

func GetLocalAESKey() []byte {
	// salt for aes key has creation by text "local-key-PulseOfFraijin"
	saltForAESKey := "local-key-PulseOfFraijin"
	// make aes key from saltForAESKey
	aesKey := sha256.Sum256([]byte(saltForAESKey))
	// aesKey is 32 bytes
	return aesKey[:]
}
