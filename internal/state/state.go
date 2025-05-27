package state

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"encoding/binary"
	"encoding/json"
	"errors"
	"hash/crc32"
	"os"
	"path/filepath"
)

// State holds persistent game data such as high scores.
type State struct {
	HighScore int `json:"high_score"`
	// Future fields can be added here
}

var encryptionKey = generateKey()

// generateKey creates a 32-byte AES key from system-specific data.
func generateKey() []byte {
	user, _ := os.UserHomeDir()
	sum := sha256.Sum256([]byte(user + "_pacmanai_secret"))
	return sum[:]
}

// Save persists the state to an encrypted file with integrity check.
func Save(s State) error {
	path, err := getSavePath()
	if err != nil {
		return err
	}

	// Serialize to JSON
	raw, err := json.Marshal(s)
	if err != nil {
		return err
	}

	// Prepend CRC32 checksum
	crc := crc32.ChecksumIEEE(raw)
	data := make([]byte, 4+len(raw))
	binary.LittleEndian.PutUint32(data[:4], crc)
	copy(data[4:], raw)

	// Encrypt
	encrypted, err := encrypt(data)
	if err != nil {
		return err
	}

	// Save to file
	return os.WriteFile(path, encrypted, 0644)
}

// Load reads the state from disk, decrypts and verifies it.
func Load() State {
	var s State

	path, err := getSavePath()
	if err != nil {
		return s
	}

	encrypted, err := os.ReadFile(path)
	if err != nil {
		return s
	}

	decrypted, err := decrypt(encrypted)
	if err != nil || len(decrypted) < 5 {
		return s
	}

	crcStored := binary.LittleEndian.Uint32(decrypted[:4])
	payload := decrypted[4:]
	if crc32.ChecksumIEEE(payload) != crcStored {
		return s
	}

	_ = json.Unmarshal(payload, &s)
	return s
}

// ======================
// 🔐 AES Encryption
// ======================

func encrypt(plain []byte) ([]byte, error) {
	block, err := aes.NewCipher(encryptionKey)
	if err != nil {
		return nil, err
	}
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}
	nonce := make([]byte, gcm.NonceSize())
	if _, err := rand.Read(nonce); err != nil {
		return nil, err
	}
	return gcm.Seal(nonce, nonce, plain, nil), nil
}

func decrypt(ciphertext []byte) ([]byte, error) {
	block, err := aes.NewCipher(encryptionKey)
	if err != nil {
		return nil, err
	}
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}
	if len(ciphertext) < gcm.NonceSize() {
		return nil, errors.New("ciphertext too short")
	}
	return gcm.Open(nil, ciphertext[:gcm.NonceSize()], ciphertext[gcm.NonceSize():], nil)
}

// getSavePath returns the path to the save file inside the user config directory.
func getSavePath() (string, error) {
	configDir, err := os.UserConfigDir()
	if err != nil {
		return "", err
	}
	saveDir := filepath.Join(configDir, "pacmanai")
	if err := os.MkdirAll(saveDir, 0755); err != nil {
		return "", err
	}
	return filepath.Join(saveDir, "state.dat"), nil
}
