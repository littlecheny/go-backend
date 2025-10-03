package services

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"io"

	"golang.org/x/crypto/pbkdf2"
	"golang.org/x/crypto/bcrypt"
	"github.com/littlecheny/go-backend/domain"
)

type cryptoService struct {
	defaultSaltSize int
	defaultKeySize  int
	iterations      int
}

func NewCryptoService() domain.CryptoService {
	return &cryptoService{
		defaultSaltSize: 32,
		defaultKeySize:  32,
		iterations:      10000,
	}
}

// Encrypt 使用AES-GCM加密数据
func (c *cryptoService) Encrypt(data string, key string) (string, error) {
	// 将密钥转换为32字节
	keyBytes := c.deriveKeyFromString(key)
	
	// 创建AES cipher
	block, err := aes.NewCipher(keyBytes)
	if err != nil {
		return "", fmt.Errorf("failed to create cipher: %v", err)
	}

	// 创建GCM
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", fmt.Errorf("failed to create GCM: %v", err)
	}

	// 生成随机nonce
	nonce := make([]byte, gcm.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return "", fmt.Errorf("failed to generate nonce: %v", err)
	}

	// 加密数据
	ciphertext := gcm.Seal(nonce, nonce, []byte(data), nil)

	// 返回base64编码的结果
	return base64.StdEncoding.EncodeToString(ciphertext), nil
}

// Decrypt 使用AES-GCM解密数据
func (c *cryptoService) Decrypt(encryptedData string, key string) (string, error) {
	// 解码base64
	ciphertext, err := base64.StdEncoding.DecodeString(encryptedData)
	if err != nil {
		return "", fmt.Errorf("failed to decode base64: %v", err)
	}

	// 将密钥转换为32字节
	keyBytes := c.deriveKeyFromString(key)

	// 创建AES cipher
	block, err := aes.NewCipher(keyBytes)
	if err != nil {
		return "", fmt.Errorf("failed to create cipher: %v", err)
	}

	// 创建GCM
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", fmt.Errorf("failed to create GCM: %v", err)
	}

	// 检查密文长度
	nonceSize := gcm.NonceSize()
	if len(ciphertext) < nonceSize {
		return "", fmt.Errorf("ciphertext too short")
	}

	// 分离nonce和密文
	nonce, ciphertext := ciphertext[:nonceSize], ciphertext[nonceSize:]

	// 解密
	plaintext, err := gcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return "", fmt.Errorf("failed to decrypt: %v", err)
	}

	return string(plaintext), nil
}

// DeriveKey 使用PBKDF2派生密钥
func (c *cryptoService) DeriveKey(password string, salt string) (string, error) {
	saltBytes, err := base64.StdEncoding.DecodeString(salt)
	if err != nil {
		return "", fmt.Errorf("failed to decode salt: %v", err)
	}

	key := pbkdf2.Key([]byte(password), saltBytes, c.iterations, c.defaultKeySize, sha256.New)
	return base64.StdEncoding.EncodeToString(key), nil
}

// GenerateSalt 生成随机盐值
func (c *cryptoService) GenerateSalt() (string, error) {
	salt := make([]byte, c.defaultSaltSize)
	if _, err := rand.Read(salt); err != nil {
		return "", fmt.Errorf("failed to generate salt: %v", err)
	}
	return base64.StdEncoding.EncodeToString(salt), nil
}

// Hash 使用bcrypt哈希数据
func (c *cryptoService) Hash(data string) string {
	hash, err := bcrypt.GenerateFromPassword([]byte(data), bcrypt.DefaultCost)
	if err != nil {
		// 如果bcrypt失败，使用SHA256作为备选
		h := sha256.Sum256([]byte(data))
		return base64.StdEncoding.EncodeToString(h[:])
	}
	return string(hash)
}

// VerifyHash 验证哈希
func (c *cryptoService) VerifyHash(data string, hash string) bool {
	// 首先尝试bcrypt验证
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(data))
	if err == nil {
		return true
	}

	// 如果bcrypt失败，尝试SHA256验证（向后兼容）
	h := sha256.Sum256([]byte(data))
	expectedHash := base64.StdEncoding.EncodeToString(h[:])
	return expectedHash == hash
}

// deriveKeyFromString 从字符串派生32字节密钥
func (c *cryptoService) deriveKeyFromString(key string) []byte {
	hash := sha256.Sum256([]byte(key))
	return hash[:]
}

// EncryptWithPassword 使用密码加密（包含盐值派生）
func (c *cryptoService) EncryptWithPassword(data string, password string) (string, string, error) {
	// 生成盐值
	salt, err := c.GenerateSalt()
	if err != nil {
		return "", "", fmt.Errorf("failed to generate salt: %v", err)
	}

	// 派生密钥
	derivedKey, err := c.DeriveKey(password, salt)
	if err != nil {
		return "", "", fmt.Errorf("failed to derive key: %v", err)
	}

	// 加密数据
	encryptedData, err := c.Encrypt(data, derivedKey)
	if err != nil {
		return "", "", fmt.Errorf("failed to encrypt data: %v", err)
	}

	return encryptedData, salt, nil
}

// DecryptWithPassword 使用密码解密（包含盐值派生）
func (c *cryptoService) DecryptWithPassword(encryptedData string, password string, salt string) (string, error) {
	// 派生密钥
	derivedKey, err := c.DeriveKey(password, salt)
	if err != nil {
		return "", fmt.Errorf("failed to derive key: %v", err)
	}

	// 解密数据
	data, err := c.Decrypt(encryptedData, derivedKey)
	if err != nil {
		return "", fmt.Errorf("failed to decrypt data: %v", err)
	}

	return data, nil
}

// GenerateRandomKey 生成随机密钥
func (c *cryptoService) GenerateRandomKey(size int) (string, error) {
	if size <= 0 {
		size = c.defaultKeySize
	}

	key := make([]byte, size)
	if _, err := rand.Read(key); err != nil {
		return "", fmt.Errorf("failed to generate random key: %v", err)
	}

	return base64.StdEncoding.EncodeToString(key), nil
}

// ValidatePassword 验证密码强度
func (c *cryptoService) ValidatePassword(password string) error {
	if len(password) < 8 {
		return fmt.Errorf("password must be at least 8 characters long")
	}

	hasUpper := false
	hasLower := false
	hasDigit := false
	hasSpecial := false

	for _, char := range password {
		switch {
		case 'A' <= char && char <= 'Z':
			hasUpper = true
		case 'a' <= char && char <= 'z':
			hasLower = true
		case '0' <= char && char <= '9':
			hasDigit = true
		default:
			hasSpecial = true
		}
	}

	if !hasUpper {
		return fmt.Errorf("password must contain at least one uppercase letter")
	}
	if !hasLower {
		return fmt.Errorf("password must contain at least one lowercase letter")
	}
	if !hasDigit {
		return fmt.Errorf("password must contain at least one digit")
	}
	if !hasSpecial {
		return fmt.Errorf("password must contain at least one special character")
	}

	return nil
}