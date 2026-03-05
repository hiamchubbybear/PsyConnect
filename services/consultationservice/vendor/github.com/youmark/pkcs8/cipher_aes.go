package pkcs8

import (
	"crypto/aes"
	"encoding/asn1"
)

var (
	oidAES128CBC = asn1.ObjectIdentifier{2, 16, 840, 1, 101, 3, 4, 1, 2}
	oidAES128GCM = asn1.ObjectIdentifier{2, 16, 840, 1, 101, 3, 4, 1, 6}
	oidAES192CBC = asn1.ObjectIdentifier{2, 16, 840, 1, 101, 3, 4, 1, 22}
	oidAES192GCM = asn1.ObjectIdentifier{2, 16, 840, 1, 101, 3, 4, 1, 26}
	oidAES256CBC = asn1.ObjectIdentifier{2, 16, 840, 1, 101, 3, 4, 1, 42}
	oidAES256GCM = asn1.ObjectIdentifier{2, 16, 840, 1, 101, 3, 4, 1, 46}
)

func init() {
	RegisterCipher(oidAES128CBC, func() Cipher {
		return AES128CBC
	})
	RegisterCipher(oidAES128GCM, func() Cipher {
		return AES128GCM
	})
	RegisterCipher(oidAES192CBC, func() Cipher {
		return AES192CBC
	})
	RegisterCipher(oidAES192GCM, func() Cipher {
		return AES192GCM
	})
	RegisterCipher(oidAES256CBC, func() Cipher {
		return AES256CBC
	})
	RegisterCipher(oidAES256GCM, func() Cipher {
		return AES256GCM
	})
}


var AES128CBC = cipherWithBlock{
	ivSize:   aes.BlockSize,
	keySize:  16,
	newBlock: aes.NewCipher,
	oid:      oidAES128CBC,
}


var AES128GCM = cipherWithBlock{
	ivSize:   aes.BlockSize,
	keySize:  16,
	newBlock: aes.NewCipher,
	oid:      oidAES128GCM,
}


var AES192CBC = cipherWithBlock{
	ivSize:   aes.BlockSize,
	keySize:  24,
	newBlock: aes.NewCipher,
	oid:      oidAES192CBC,
}


var AES192GCM = cipherWithBlock{
	ivSize:   aes.BlockSize,
	keySize:  24,
	newBlock: aes.NewCipher,
	oid:      oidAES192GCM,
}


var AES256CBC = cipherWithBlock{
	ivSize:   aes.BlockSize,
	keySize:  32,
	newBlock: aes.NewCipher,
	oid:      oidAES256CBC,
}


var AES256GCM = cipherWithBlock{
	ivSize:   aes.BlockSize,
	keySize:  32,
	newBlock: aes.NewCipher,
	oid:      oidAES256GCM,
}
