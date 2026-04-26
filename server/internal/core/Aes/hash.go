package Aes

import (
	"fmt"
	"lab2pepe/server/internal/core/md5MY"
)

const HashSize = 16 // MD5 = 16 байт

// AppendHash считает MD5 от данных и добавляет в конец
func AppendHash(data []byte) []byte {
	// копируем чтобы HASH не испортил оригинал
	dataCopy := make([]byte, len(data))
	copy(dataCopy, data)

	hash := md5MY.HASH(dataCopy)
	return append(data, hash...)
}

// VerifyAndStrip проверяет MD5 и возвращает данные без хэша
func VerifyAndStrip(data []byte) ([]byte, error) {
	if len(data) < HashSize {
		return nil, fmt.Errorf("data too short")
	}

	content := data[:len(data)-HashSize]

	// сохраняем копию хэша ДО вызова HASH
	receivedHash := make([]byte, HashSize)
	copy(receivedHash, data[len(data)-HashSize:])

	// теперь считаем ожидаемый хэш — он может затронуть data через append
	expectedHash := md5MY.HASH(content)

	fmt.Printf("receivedHash: %x\n", receivedHash)
	fmt.Printf("expectedHash: %x\n", expectedHash)

	for i := range expectedHash {
		if expectedHash[i] != receivedHash[i] {
			return nil, fmt.Errorf("hash mismatch: file corrupted")
		}
	}

	return content, nil
}
