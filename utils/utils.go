package utils

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/fxamacker/cbor/v2"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

var Validate = validator.New()

func WriteJSON(w http.ResponseWriter, status int, v any) error {
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(status)
	return json.NewEncoder(w).Encode(v)
}

func WriteError(w http.ResponseWriter, status int, err error) {
	WriteJSON(w, status, map[string]string{"error": err.Error()})
}

func ParseJSON(r *http.Request, v any) error {
	if r.Body == nil {
		return fmt.Errorf("missing request body")
	}

	return json.NewDecoder(r.Body).Decode(v)
}

func GetTokenFromRequest(r *http.Request) string {
	tokenAuth := r.Header.Get("Authorization")
	tokenQuery := r.URL.Query().Get("token")

	if tokenAuth != "" {
		return tokenAuth
	}

	if tokenQuery != "" {
		return tokenQuery
	}

	return ""
}

func StringToBytes32(hexStr string) ([32]byte, error) {
	var b32 [32]byte

	hexStr = strings.TrimPrefix(hexStr, "0x")

	bytes, err := hex.DecodeString(hexStr)
	if err != nil {
		return b32, err
	}
	if len(bytes) != 32 {
		return b32, fmt.Errorf("invalid length: got %d, expected 32", len(bytes))
	}

	copy(b32[:], bytes)
	return b32, nil
}

func CalculateHash(data any) ([32]byte, error) {
	dataBytes, err := cbor.Marshal(data)
	if err != nil {
		return [32]byte{}, fmt.Errorf("failed to marshal data: %w", err)
	}

	hash := sha256.Sum256(dataBytes)
	return hash, nil
}

func ToBinaryUUID(u uuid.UUID) primitive.Binary {
	return primitive.Binary{Subtype: 4, Data: u[:]}
}
