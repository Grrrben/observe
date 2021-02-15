package observe

import (
	"crypto/sha256"
	"encoding"
	"fmt"
	"log"
)

func HashString(s string) string {
	h := sha256.New()
	h.Write([]byte(s))

	marshaler, ok := h.(encoding.BinaryMarshaler)
	if !ok {
		log.Fatal("h does not implement encoding.BinaryMarshaler")
	}
	_, err := marshaler.MarshalBinary()
	if err != nil {
		log.Fatal("unable to marshal hash:", err)
	}

	return fmt.Sprintf("%x", h.Sum(nil))
}
