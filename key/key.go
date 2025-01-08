// Copyright (C) aliacme. 2025-present.
//
// Created at 2025-01-07, by liasica

package key

import (
	"crypto/ecdsa"
	"crypto/x509"
	"encoding/pem"
)

func EncodePrivateKey(private *ecdsa.PrivateKey) []byte {
	x509Encoded, _ := x509.MarshalECPrivateKey(private)
	pemEncoded := pem.EncodeToMemory(&pem.Block{Type: "PRIVATE KEY", Bytes: x509Encoded})
	return pemEncoded
}

func DecodePrivateKey(b []byte) (*ecdsa.PrivateKey, error) {
	block, _ := pem.Decode(b)
	return x509.ParseECPrivateKey(block.Bytes)
}

func EncodePublicKey(public *ecdsa.PublicKey) []byte {
	x509EncodedPub, _ := x509.MarshalPKIXPublicKey(public)
	pemEncodedPub := pem.EncodeToMemory(&pem.Block{Type: "PUBLIC KEY", Bytes: x509EncodedPub})
	return pemEncodedPub
}

func DecodePublicKey(b []byte) (*ecdsa.PublicKey, error) {
	block, _ := pem.Decode(b)
	genericPublicKey, _ := x509.ParsePKIXPublicKey(block.Bytes)
	return genericPublicKey.(*ecdsa.PublicKey), nil
}
