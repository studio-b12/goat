package util

import "io"

// ReadReaderToString takes a reader and an error
// and returns the passed error back if it is
// not nil.
//
// Otherwise, if r is not nil, the contents of
// r are read and returned as string.
//
// This function is designed to be used as a read
// wrapper for a reader getter. For example:
//
//	reader, err := r.Reader()
//	data, err := util.ReadReaderToString(reader, err)
//
// Or short:
//
//	data, err := util.ReadReaderToString(r.Reader())
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
