//go:build js && wasm

package main

import (
	"bytes"
	"encoding/base64"
	"encoding/gob"
	"syscall/js"

	"github.com/fentec-project/gofe/abe"
)

func serialize(v interface{}) (string, error) {
	var buf bytes.Buffer
	if err := gob.NewEncoder(&buf).Encode(v); err != nil {
		return "", err
	}
	return base64.StdEncoding.EncodeToString(buf.Bytes()), nil
}

func deserialize(b64 string, v interface{}) error {
	data, err := base64.StdEncoding.DecodeString(b64)
	if err != nil {
		return err
	}
	return gob.NewDecoder(bytes.NewReader(data)).Decode(v)
}

// 1. Client-side Encryption (using Public Key + Policy)
func JSEncrypt(this js.Value, args []js.Value) interface{} {
	// Args: pubKeyB64, message, policyString
	if len(args) < 3 {
		return js.ValueOf(map[string]interface{}{"error": "missing arguments"})
	}

	pk := new(abe.FAMEPubKey)
	if err := deserialize(args[0].String(), pk); err != nil {
		return js.ValueOf(map[string]interface{}{"error": "invalid public key: " + err.Error()})
	}

	msg := args[1].String()
	policy := args[2].String()

	msp, err := abe.BooleanToMSP(policy, false)
	if err != nil {
		return js.ValueOf(map[string]interface{}{"error": "invalid policy: " + err.Error()})
	}

	a := abe.NewFAME()
	cipher, err := a.Encrypt(msg, msp, pk)
	if err != nil {
		return js.ValueOf(map[string]interface{}{"error": err.Error()})
	}

	cipherB64, _ := serialize(cipher)
	return js.ValueOf(map[string]interface{}{"cipher": cipherB64})
}

// 2. Client-side Decryption (using Attribute Keys + Public Key)
func JSDecrypt(this js.Value, args []js.Value) interface{} {
	// Args: cipherB64, attribKeysB64, pubKeyB64
	if len(args) < 3 {
		return js.ValueOf(map[string]interface{}{"error": "missing arguments"})
	}

	cipher := new(abe.FAMECipher)
	deserialize(args[0].String(), cipher)

	keys := new(abe.FAMEAttribKeys)
	deserialize(args[1].String(), keys)

	pk := new(abe.FAMEPubKey)
	deserialize(args[2].String(), pk)

	a := abe.NewFAME()
	msg, err := a.Decrypt(cipher, keys, pk)
	if err != nil {
		return js.ValueOf(map[string]interface{}{"error": "decryption failed: " + err.Error()})
	}

	return js.ValueOf(map[string]interface{}{"message": msg})
}

func main() {
	js.Global().Set("fameEncrypt", js.FuncOf(JSEncrypt))
	js.Global().Set("fameDecrypt", js.FuncOf(JSDecrypt))
	select {}
}