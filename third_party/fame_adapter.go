package fame

/*
 * Copyright (c) 2019 ATOS
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 * http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */
import (
	"bytes"
	"encoding/base64"
	"encoding/gob"
	"encoding/hex"
	"fmt"
	"log"

	"github.com/fentec-project/gofe/abe"
)

// GenerateMasterKeys creates master keys and saves them to files
func GenerateMasterKeys(path string, debug bool) error {
	if debug {
		log.Println("GenerateMasterKeys: init")
	}

	a := abe.NewFAME()
	if err := WriteGob_pn(path, "fame.gob", a); err != nil {
		return fmt.Errorf("failed to write fame file: %w", err)
	}

	pubKey, secKey, err := a.GenerateMasterKeys()
	if err != nil {
		return fmt.Errorf("failed to generate master keys: %w", err)
	}

	if err := WriteGob_pn(path, "publicKey.gob", pubKey); err != nil {
		return fmt.Errorf("failed to write public key: %w", err)
	}
	if err := WriteGob_pn(path, "secretKey.gob", secKey); err != nil {
		return fmt.Errorf("failed to write secret key: %w", err)
	}

	return nil
}

// Encrypt encrypts a message with the given policy
func Encrypt(path string, msg string, policy string, debug bool) (string, error) {
	if debug {
		log.Printf("Encrypt: policy: %s", policy)
	}

	msp, err := abe.BooleanToMSP(policy, false)
	if err != nil {
		return "", fmt.Errorf("invalid policy: %w", err)
	}

	a := new(abe.FAME)
	if err := ReadGob_pn(path, "fame.gob", a); err != nil {
		return "", fmt.Errorf("failed to load fame: %w", err)
	}

	pk := new(abe.FAMEPubKey)
	if err := ReadGob_pn(path, "publicKey.gob", pk); err != nil {
		return "", fmt.Errorf("failed to load public key: %w", err)
	}

	cipher, err := a.Encrypt(msg, msp, pk)
	if err != nil {
		return "", fmt.Errorf("encryption failed: %w", err)
	}

	cipherBytes := WriteTCharA(cipher)
	return hex.EncodeToString(cipherBytes), nil
}

// GenerateAttribKeys generates attribute keys for given attributes
func GenerateAttribKeys(path string, pathGen string, gamma []string, debug bool) error {
	if debug {
		log.Printf("GenerateAttribKeys: gamma=%v", gamma)
	}

	a := new(abe.FAME)
	if err := ReadGob_pn(path, "fame.gob", a); err != nil {
		return fmt.Errorf("failed to load fame: %w", err)
	}

	sk := new(abe.FAMESecKey)
	if err := ReadGob_pn(path, "secretKey.gob", sk); err != nil {
		return fmt.Errorf("failed to load secret key: %w", err)
	}

	keys, err := a.GenerateAttribKeys(gamma, sk)
	if err != nil {
		return fmt.Errorf("key generation failed: %w", err)
	}

	if err := WriteGob_pn(pathGen, "genKey.gob", keys); err != nil {
		return fmt.Errorf("failed to save keys: %w", err)
	}

	return nil
}

// Decrypt decrypts a ciphertext using generated keys
func Decrypt(path string, pathGen string, cipherHex string, debug bool) (string, error) {
	if debug {
		log.Printf("Decrypt: cipher length=%d", len(cipherHex))
	}

	// Decode hex ciphertext
	cipherBytes, err := hex.DecodeString(cipherHex)
	if err != nil {
		return "", fmt.Errorf("invalid hex ciphertext: %w", err)
	}

	// Load required components
	a := new(abe.FAME)
	if err := ReadGob_pn(path, "fame.gob", a); err != nil {
		return "", fmt.Errorf("failed to load fame: %w", err)
	}

	pk := new(abe.FAMEPubKey)
	if err := ReadGob_pn(path, "publicKey.gob", pk); err != nil {
		return "", fmt.Errorf("failed to load public key: %w", err)
	}

	key := new(abe.FAMEAttribKeys)
	if err := ReadGob_pn(pathGen, "genKey.gob", key); err != nil {
		return "", fmt.Errorf("failed to load attribute keys: %w", err)
	}

	// Parse ciphertext
	cipher := new(abe.FAMECipher)
	ReadFCharA(cipherBytes, cipher)

	// Decrypt
	plaintext, err := a.Decrypt(cipher, key, pk)
	if err != nil {
		return "", fmt.Errorf("decryption failed: %w", err)
	}

	return plaintext, nil
}

func GenerateAttribKeysToBase64(path string, gamma []string, debug bool) (string, error) {
	if debug {
		log.Printf("GenerateAttribKeysToBase64: gamma=%v", gamma)
	}

	// 1. Load FAME parameters
	a := new(abe.FAME)
	if err := ReadGob_pn(path, "fame.gob", a); err != nil {
		return "", fmt.Errorf("failed to load fame: %w", err)
	}

	// 2. Load Master Secret Key
	sk := new(abe.FAMESecKey)
	if err := ReadGob_pn(path, "secretKey.gob", sk); err != nil {
		return "", fmt.Errorf("failed to load secret key: %w", err)
	}

	// 3. Generate attribute keys using FAME
	keys, err := a.GenerateAttribKeys(gamma, sk)
	if err != nil {
		return "", fmt.Errorf("key generation failed: %w", err)
	}

	// 4. Serialize struct to gob bytes and encode to Base64 string
	var buf bytes.Buffer
	if err := gob.NewEncoder(&buf).Encode(keys); err != nil {
		return "", fmt.Errorf("failed to encode attribute keys: %w", err)
	}

	return base64.StdEncoding.EncodeToString(buf.Bytes()), nil
}
