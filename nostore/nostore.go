package nostore

import (
	"bytes"
	"compress/gzip"
	"encoding/base64"
	"io/ioutil"
	"strconv"
	"time"

	"github.com/gtank/cryptopasta"
)

var encryptionKey = cryptopasta.NewEncryptionKey()

func Encode(data []byte, until time.Time) (string, error) {
	expiresAt := strconv.FormatInt(until.Unix(), 10)

	compressed, err := compress(data)
	if err != nil {
		return "", err
	}

	encrypted, err := cryptopasta.Encrypt(
		append([]byte(expiresAt+"."), compressed...),
		encryptionKey,
	)
	if err != nil {
		return "", err
	}

	return base64Encode(encrypted), nil
}

func Decode(encoded string) ([]byte, time.Time, bool, error) {
	output := []byte{}
	now := time.Now().Unix()
	expiresAt := time.Now()

	decoded, err := base64Decode(encoded)
	if err != nil {
		return output, expiresAt, false, err
	}

	decrypted, err := cryptopasta.Decrypt([]byte(decoded), encryptionKey)
	if err != nil {
		return output, expiresAt, false, err
	}

	values := bytes.SplitN(decrypted, []byte("."), 2)

	until, err := strconv.ParseInt(string(values[0]), 10, 64)
	if err != nil {
		return output, expiresAt, false, err
	}

	expiresAt = time.Unix(until, 0)

	if now > until {
		return output, expiresAt, true, nil
	}

	data, err := decompress(values[1])
	if err != nil {
		return output, expiresAt, false, err
	}

	return data, expiresAt, false, nil
}

func compressToString(data []byte) (string, error) {
	var buf bytes.Buffer

	gzipWriter := gzip.NewWriter(&buf)
	if _, err := gzipWriter.Write(data); err != nil {
		return "", err
	}

	if err := gzipWriter.Flush(); err != nil {
		return "", err
	}

	if err := gzipWriter.Close(); err != nil {
		return "", err
	}

	return base64Encode(buf.Bytes()), nil
}

func compress(data []byte) ([]byte, error) {
	var buf bytes.Buffer
	output := []byte{}

	gzipWriter := gzip.NewWriter(&buf)
	if _, err := gzipWriter.Write(data); err != nil {
		return output, err
	}

	if err := gzipWriter.Flush(); err != nil {
		return output, err
	}

	if err := gzipWriter.Close(); err != nil {
		return output, err
	}

	return buf.Bytes(), nil
}

func decompress(data []byte) ([]byte, error) {
	output := []byte{}

	readableBytes := bytes.NewReader(data)
	gzipReader, err := gzip.NewReader(readableBytes)
	if err != nil {
		return output, err
	}

	return ioutil.ReadAll(gzipReader)
}

func base64Encode(data []byte) string {
	return base64.StdEncoding.EncodeToString(data)
}

func base64Decode(data string) ([]byte, error) {
	return base64.StdEncoding.DecodeString(data)
}
