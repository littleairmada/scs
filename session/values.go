package session

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
)

// ErrKeyNotFound is returned by operations on session values where the given
// key could not be found.
var ErrKeyNotFound = errors.New("key not found in session values")

// ErrTypeAssertionFailed is returned by operations on session values where the
// received value could not be type asserted or converted into the required type.
var ErrTypeAssertionFailed = errors.New("type assertion failed")

// GetString returns the string value for a given key from the session. If the
// key is not present a ErrKeyNotFound error is returned. If the value could not
// be type asserted to a string then a ErrTypeAssertionFailed error is returned.
func GetString(r *http.Request, key string) (string, error) {
	s, err := sessionFromContext(r)
	if err != nil {
		return "", err
	}

	s.mu.RLock()
	defer s.mu.RUnlock()

	v, exists := s.values[key]
	if exists == false {
		return "", ErrKeyNotFound
	}

	str, ok := v.(string)
	if ok == false {
		return "", ErrTypeAssertionFailed
	}

	return str, nil
}

// PutString adds a string value and corresponding key to the the session. Any
// existing values for the key will be replaced.
func PutString(r *http.Request, key string, val string) error {
	s, err := sessionFromContext(r)
	if err != nil {
		return err
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	if s.written == true {
		return ErrAlreadyWritten
	}
	s.values[key] = val
	s.modified = true
	return nil
}

// PopString returns the string value for a given key, and removes both the key
// and value from the session. If the key is not present a ErrKeyNotFound error
// is returned. If the value could not be type asserted to a string then a
// ErrTypeAssertionFailed error is returned and the data is not removed.
func PopString(r *http.Request, key string) (string, error) {
	s, err := sessionFromContext(r)
	if err != nil {
		return "", err
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	if s.written == true {
		return "", ErrAlreadyWritten
	}
	v, exists := s.values[key]
	if exists == false {
		return "", ErrKeyNotFound
	}

	str, ok := v.(string)
	if ok == false {
		return "", ErrTypeAssertionFailed
	}

	delete(s.values, key)
	s.modified = true
	return str, nil
}

// GetBool returns the bool value for a given key from the session. If the key
// is not present a ErrKeyNotFound error is returned. If the value could not be
// type asserted to a bool then a ErrTypeAssertionFailed error is returned.
func GetBool(r *http.Request, key string) (bool, error) {
	s, err := sessionFromContext(r)
	if err != nil {
		return false, err
	}

	s.mu.RLock()
	defer s.mu.RUnlock()

	v, exists := s.values[key]
	if exists == false {
		return false, ErrKeyNotFound
	}

	b, ok := v.(bool)
	if ok == false {
		return false, ErrTypeAssertionFailed
	}

	return b, nil
}

// PutBool adds a bool value and corresponding key to the session. Any existing
// values for the key will be replaced.
func PutBool(r *http.Request, key string, val bool) error {
	s, err := sessionFromContext(r)
	if err != nil {
		return err
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	if s.written == true {
		return ErrAlreadyWritten
	}
	s.values[key] = val
	s.modified = true
	return nil
}

// PopBool returns the bool value for a given key from the session values,
// and removes both the key and value from the session. If the key is not present
// a ErrKeyNotFound error is returned. If the value could not be type asserted
// to a bool then a ErrTypeAssertionFailed error is returned and the data is
// not removed.
func PopBool(r *http.Request, key string) (bool, error) {
	s, err := sessionFromContext(r)
	if err != nil {
		return false, err
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	if s.written == true {
		return false, ErrAlreadyWritten
	}
	v, exists := s.values[key]
	if exists == false {
		return false, ErrKeyNotFound
	}

	b, ok := v.(bool)
	if ok == false {
		return false, ErrTypeAssertionFailed
	}

	delete(s.values, key)
	s.modified = true
	return b, nil
}

// GetInt returns the int value for a given key from the session. If the key is
// not present a ErrKeyNotFound error is returned. If the value could not be type
// asserted or converted to a int then a ErrTypeAssertionFailed error is returned.
func GetInt(r *http.Request, key string) (int, error) {
	s, err := sessionFromContext(r)
	if err != nil {
		return 0, err
	}

	s.mu.RLock()
	defer s.mu.RUnlock()

	v, exists := s.values[key]
	if exists == false {
		return 0, ErrKeyNotFound
	}

	switch v := v.(type) {
	case int:
		return v, nil
	case json.Number:
		return strconv.Atoi(v.String())
	}

	return 0, ErrTypeAssertionFailed
}

// PutInt adds an int value and corresponding key to the session. Any existing
// values for the key will be replaced.
func PutInt(r *http.Request, key string, val int) error {
	s, err := sessionFromContext(r)
	if err != nil {
		return err
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	if s.written == true {
		return ErrAlreadyWritten
	}
	s.values[key] = val
	s.modified = true
	return nil
}

// PopInt returns the int value for a given key, and removes both the key
// and value from the session. If the key is not present a ErrKeyNotFound error
// is returned. If the value could not be converted to an int then a
// ErrTypeAssertionFailed error is returned and the data is not removed.
func PopInt(r *http.Request, key string) (int, error) {
	s, err := sessionFromContext(r)
	if err != nil {
		return 0, err
	}

	s.mu.RLock()
	defer s.mu.RUnlock()

	if s.written == true {
		return 0, ErrAlreadyWritten
	}
	v, exists := s.values[key]
	if exists == false {
		return 0, ErrKeyNotFound
	}

	var i int
	switch v := v.(type) {
	case int:
		i = v
	case json.Number:
		i, err = strconv.Atoi(v.String())
		if err != nil {
			return 0, err
		}
	default:
		return 0, ErrTypeAssertionFailed
	}

	delete(s.values, key)
	s.modified = true
	return i, nil
}

// Remove deletes the given key and corresponding value from the session.
// If the key is not present this operation is a no-op.
func Remove(r *http.Request, key string) error {
	s, err := sessionFromContext(r)
	if err != nil {
		return err
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	if s.written == true {
		return ErrAlreadyWritten
	}
	delete(s.values, key)
	s.modified = true
	return nil
}

// Clear deletes all keys and values stored in the session.
func Clear(r *http.Request) error {
	s, err := sessionFromContext(r)
	if err != nil {
		return err
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	if s.written == true {
		return ErrAlreadyWritten
	}
	for key := range s.values {
		delete(s.values, key)
	}
	s.modified = true
	return nil
}
