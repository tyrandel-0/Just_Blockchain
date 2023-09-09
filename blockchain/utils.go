package blockchain

import (
	"bytes"
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/base64"
	"encoding/binary"
	"encoding/json"
)

func SerializeBlock(block *Block) string {
	jsonData, err := json.MarshalIndent(*block, "", "\t")
	if err != nil {
		return ""
	}
	return string(jsonData)
}

func GenerateRandomBytes(max uint) []byte {
	var slice []byte = make([]byte, max)
	_, err := rand.Read(slice)
	if err != nil {
		return nil
	}
	return slice
}

func StringPublic(pub *rsa.PublicKey) string {
	return base64.StdEncoding.EncodeToString(x509.MarshalPKCS1PublicKey(pub))
}

func HashSum(data []byte) []byte {
	hash := sha256.Sum256(data)
	return hash[:]
}

func ToBytes(num uint64) []byte {
	var data = new(bytes.Buffer)
	err := binary.Write(data, binary.BigEndian, num)
	if err != nil {
		return nil
	}

	return data.Bytes()
}

func Sign(priv *rsa.PrivateKey, data []byte) []byte {
	signature, err := rsa.SignPSS(rand.Reader, priv, crypto.SHA3_256, data, nil)
	if err != nil {
		return nil
	}

	return signature
}

func DeserializeBlock(data string) *Block {
	var block Block
	err := json.Unmarshal([]byte(data), &block)
	if err != nil {
		return nil
	}
	return &block
}
