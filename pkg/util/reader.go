package util

import "io"

func ReadReaderToString(r io.Reader, err error) (string, error) {
	if err != nil {
		return "", err
	}

	if r == nil {
		return "", nil
	}

	data, err := io.ReadAll(r)
	if err != nil {
		return "", err
	}

	return string(data), nil
}
