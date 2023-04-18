// SPDX-License-Identifier: Apache-2.0
//

package encoding

import (
	"crypto/sha256"
	"encoding/hex"
	"sort"
	"strings"
)

// UniqueID returns a unique and deterministic ID based on the given strings. Uses sha256 hash.
func UniqueID(input []string) string {
	// Sort the input to ensure deterministic output
	sort.Strings(input)

	// Concatenate the sorted strings
	concatenated := strings.Join(input, "")

	// Generate the SHA-256 hash of the concatenated strings
	hash := sha256.Sum256([]byte(concatenated))

	// Convert the hash to a hexadecimal string
	return hex.EncodeToString(hash[:])
}
