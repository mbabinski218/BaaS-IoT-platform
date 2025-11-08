package utils

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/fxamacker/cbor/v2"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"github.com/mbabinski218/BaaS-IoT-platform/types"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

var Validate = validator.New()

func WriteJSON(w http.ResponseWriter, status int, v any) error {
	if w == nil {
		return nil
	}
	w.Header().Set("Content-Type", "application/json")
	if status != http.StatusOK {
		w.WriteHeader(status)
	}
	return json.NewEncoder(w).Encode(v)
}

func WriteError(w http.ResponseWriter, status int, err error) {
	if w == nil {
		return
	}
	WriteJSON(w, status, map[string]string{"error": err.Error()})
}

func ParseJSON(r *http.Request, v any) error {
	if r.Body == nil {
		return fmt.Errorf("missing request body")
	}

	return json.NewDecoder(r.Body).Decode(v)
}

func MapToBSON(payload types.NewDataPayload) bson.M {
	doc := bson.M{
		"_id":       ToBinaryUUID(payload.DataId),
		"device_id": ToBinaryUUID(payload.DeviceId),
		"data":      payload.Data,
	}
	return doc
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

func FixTimestamps(from, to time.Time, interval time.Duration) (time.Time, time.Time, error) {
	startTime, _ := time.ParseInLocation(types.TimeLayout, types.BlockchainBatchStartTime, time.UTC)
	interval = time.Duration(interval) * time.Second

	fromOffset := from.Sub(startTime)
	toOffset := to.Sub(startTime)

	fromRounded := startTime.Add((fromOffset / interval) * interval)
	toRounded := startTime.Add(((toOffset + interval - 1) / interval) * interval)

	return fromRounded, toRounded, nil
}

func CreateMerkleRoot(data []types.DocData) ([32]byte, map[uuid.UUID][][]byte, error) {
	var leaves [][32]byte
	indexToID := make(map[int]uuid.UUID)
	idToHash := make(map[uuid.UUID][32]byte)

	for i, doc := range data {
		hash, _ := CalculateHash(doc.Data)
		leaves = append(leaves, hash)
		indexToID[i] = doc.Id
		idToHash[doc.Id] = hash
	}

	root, tree := buildMerkleTreeSorted(leaves)

	audit := make(map[uuid.UUID][][]byte)
	for i, id := range indexToID {
		proof := getMerkleProofSorted(tree, i)
		audit[id] = proof
	}

	return root, audit, nil
}

func buildMerkleTreeSorted(leaves [][32]byte) ([32]byte, [][][32]byte) {
	var tree [][][32]byte
	tree = append(tree, leaves)

	level := leaves
	for len(level) > 1 {
		var nextLevel [][32]byte
		for i := 0; i < len(level); i += 2 {
			left := level[i]
			var right [32]byte
			if i+1 < len(level) {
				right = level[i+1]
			} else {
				right = left
			}

			if bytes.Compare(left[:], right[:]) > 0 {
				left, right = right, left
			}

			combined := append(left[:], right[:]...)
			parent := sha256.Sum256(combined)
			nextLevel = append(nextLevel, parent)
		}
		tree = append(tree, nextLevel)
		level = nextLevel
	}

	return tree[len(tree)-1][0], tree
}

func getMerkleProofSorted(tree [][][32]byte, index int) [][]byte {
	var proof [][]byte

	for level := 0; level < len(tree)-1; level++ {
		layer := tree[level]
		siblingIdx := index ^ 1

		var sibling [32]byte
		if siblingIdx < len(layer) {
			sibling = layer[siblingIdx]
		} else {
			sibling = layer[index]
		}

		proof = append(proof, sibling[:])
		index /= 2
	}

	return proof
}

func PauseSimulatorFile() error {
	return os.WriteFile(filepath.Join("D:\\Studia_mgr\\Praca_magisterska\\Code\\BaaS-IoT-platform\\simulators", "pause.flag"), []byte("1"), 0644)
}

func ResumeSimulatorFile() error {
	return os.WriteFile(filepath.Join("D:\\Studia_mgr\\Praca_magisterska\\Code\\BaaS-IoT-platform\\simulators", "pause.flag"), []byte("0"), 0644)
}
