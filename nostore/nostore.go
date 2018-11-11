package nostore

import (
	"bytes"
	"compress/gzip"
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"strconv"
	"strings"
	"time"

	"github.com/gtank/cryptopasta"
)

var encryptionKey = cryptopasta.NewEncryptionKey()

func Encode(data []byte, until time.Time) (string, error) {
	expiresAt := strconv.FormatInt(until.Unix(), 10)

	compressed, err := compressToString(data)
	if err != nil {
		return "", err
	}

	encrypted, err := cryptopasta.Encrypt(
		[]byte(fmt.Sprintf("%s.%s", expiresAt, compressed)),
		encryptionKey,
	)
	if err != nil {
		return "", err
	}

	return base64Encode(encrypted), nil
}

func Decode(encoded string) ([]byte, bool, error) {
	output := []byte{}
	now := time.Now().Unix()

	decoded, err := base64Decode(encoded)
	if err != nil {
		return output, false, err
	}

	decrypted, err := cryptopasta.Decrypt([]byte(decoded), encryptionKey)
	if err != nil {
		return output, false, err
	}

	values := strings.Split(string(decrypted), ".")

	until, err := strconv.ParseInt(string(values[0]), 10, 64)
	if err != nil {
		return output, false, err
	}

	if now > until {
		return output, true, nil
	}

	data, err := decompressString(values[1])
	if err != nil {
		return output, false, err
	}

	return data, false, nil
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

func decompressString(data string) ([]byte, error) {
	output := []byte{}

	decoded, err := base64Decode(data)
	if err != nil {
		return output, err
	}

	readableBytes := bytes.NewReader(decoded)
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
