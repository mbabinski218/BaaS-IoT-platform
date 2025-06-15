package utils

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/cbergoon/merkletree"
	"github.com/fxamacker/cbor/v2"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"github.com/mbabinski218/BaaS-IoT-platform/types"
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

func CalculateHash(data map[string]any) ([32]byte, error) {
	opts := cbor.CanonicalEncOptions()
	encMode, err := opts.EncMode()
	if err != nil {
		return [32]byte{}, fmt.Errorf("failed to get CBOR encoding mode: %w", err)
	}

	b, err := encMode.Marshal(data)
	if err != nil {
		return [32]byte{}, fmt.Errorf("failed to marshal data: %w", err)
	}

	return sha256.Sum256(b), nil
}

func ToBinaryUUID(u uuid.UUID) primitive.Binary {
	return primitive.Binary{Subtype: 4, Data: u[:]}
}

func ToUUID(insertedID interface{}) (uuid.UUID, error) {
	bin, ok := insertedID.(primitive.Binary)
	if !ok || bin.Subtype != 0x04 || len(bin.Data) != 16 {
		return uuid.Nil, fmt.Errorf("invalid UUID format: %+v", insertedID)
	}

	return uuid.FromBytes(bin.Data)
}

func BytesTo32(b []byte) ([32]byte, error) {
	var out [32]byte
	if len(b) != 32 {
		return out, fmt.Errorf("invalid Merkle root length: got %d, expected 32", len(b))
	}
	copy(out[:], b)
	return out, nil
}

func CreateMerkleRoot(data []types.DocData) ([32]byte, map[uuid.UUID][][]byte, error) {
	var contents []merkletree.Content
	var hashMap = make(map[string][32]byte)
	var audit = make(map[uuid.UUID][][]byte)

	for _, doc := range data {
		hash, err := CalculateHash(doc.Data)
		if err != nil {
			return [32]byte{}, nil, fmt.Errorf("failed to calculate hash for document with id: %s, err: %w", doc.Id, err)
		}

		hashMap[string(doc.Id[:])] = hash
		contents = append(contents, types.MerkleData{Hash: hash})
	}

	tree, err := merkletree.NewTree(contents)
	if err != nil {
		return [32]byte{}, nil, fmt.Errorf("failed to create Merkle tree: %w", err)
	}

	rootBytes := tree.MerkleRoot()
	root, err := BytesTo32(rootBytes)
	if err != nil {
		return [32]byte{}, nil, fmt.Errorf("failed to convert Merkle root to 32 bytes: %w", err)
	}

	for _, doc := range data {
		hash := hashMap[string(doc.Id[:])]
		leaf := types.MerkleData{Hash: hash}

		proof, _, err := tree.GetMerklePath(leaf)
		if err != nil {
			return root, nil, fmt.Errorf("proof not found for %s: %w", doc.Id, err)
		}

		audit[doc.Id] = proof
	}

	return root, audit, nil
}
