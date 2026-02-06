package security

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"fmt"
	"math/big"
	"os"
	"time"
)

// JWTConfig JWT配置结构体
type JWTConfig struct {
	Secret          string
	ExpireHours     int
	RefreshHours    int
	PrivateKeyPath  string
	PublicKeyPath   string
	PrivateKey      *rsa.PrivateKey // 运行时加载的私钥
}

// LoadJWTSecret 从环境变量加载JWT密钥
// 支持通过环境变量 GOMALL_JWT_SECRET 或 GOMALL_JWT_PRIVATE_KEY 覆盖配置
func LoadJWTSecret(config *JWTConfig) error {
	// 优先使用环境变量
	if secret := os.Getenv("GOMALL_JWT_SECRET"); secret != "" {
		config.Secret = secret
	}

	// 如果配置了私钥路径，尝试加载
	if config.PrivateKeyPath != "" {
		if _, err := os.Stat(config.PrivateKeyPath); err == nil {
			privateKey, err := LoadPrivateKey(config.PrivateKeyPath)
			if err != nil {
				return fmt.Errorf("加载JWT私钥失败: %w", err)
			}
			// 使用私钥替代对称密钥，保存到配置中
			config.Secret = "" // 私钥模式不使用secret
			config.PrivateKey = privateKey
		}
	}

	// 如果没有配置密钥，生成一个
	if config.Secret == "" && config.PrivateKeyPath == "" {
		secret := generateRandomString(64)
		config.Secret = secret
		fmt.Printf("[安全警告] JWT密钥未配置，已自动生成随机密钥。生产环境请设置 GOMALL_JWT_SECRET 环境变量。\n")
	}

	return nil
}

// LoadPrivateKey 加载RSA私钥
func LoadPrivateKey(path string) (*rsa.PrivateKey, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	block, _ := pem.Decode(data)
	if block == nil {
		return nil, fmt.Errorf("无效的PEM数据")
	}

	if x509.IsEncryptedPEMBlock(block) {
		// 需要解密（这里简化处理，实际应支持密码）
		return nil, fmt.Errorf("加密的PEM块暂不支持")
	}

	key, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		return nil, err
	}

	return key, nil
}

// GenerateRSAKeyPair 生成RSA密钥对并保存到文件
func GenerateRSAKeyPair(privateKeyPath, publicKeyPath string, bits int) error {
	privateKey, err := rsa.GenerateKey(rand.Reader, bits)
	if err != nil {
		return fmt.Errorf("生成RSA密钥对失败: %w", err)
	}

	// 保存私钥
	privateKeyBytes := x509.MarshalPKCS1PrivateKey(privateKey)
	privateKeyFile, err := os.OpenFile(privateKeyPath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		return fmt.Errorf("创建私钥文件失败: %w", err)
	}
	defer privateKeyFile.Close()

	if err := pem.Encode(privateKeyFile, &pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: privateKeyBytes,
	}); err != nil {
		return fmt.Errorf("写入私钥失败: %w", err)
	}

	// 生成自签名证书用于公钥
	template := x509.Certificate{
		SerialNumber: big.NewInt(1),
		Subject: pkix.Name{
			Organization: []string{"Gomall"},
		},
		NotBefore:             time.Now(),
		NotAfter:              time.Now().AddDate(10, 0, 0),
		KeyUsage:              x509.KeyUsageKeyEncipherment | x509.KeyUsageDigitalSignature,
		ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth},
		BasicConstraintsValid: true,
	}

	certBytes, err := x509.CreateCertificate(rand.Reader, &template, &template, &privateKey.PublicKey, privateKey)
	if err != nil {
		return fmt.Errorf("生成证书失败: %w", err)
	}

	// 保存公钥证书
	certFile, err := os.OpenFile(publicKeyPath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		return fmt.Errorf("创建公钥文件失败: %w", err)
	}
	defer certFile.Close()

	if err := pem.Encode(certFile, &pem.Block{
		Type:  "CERTIFICATE",
		Bytes: certBytes,
	}); err != nil {
		return fmt.Errorf("写入公钥失败: %w", err)
	}

	return nil
}

// generateRandomString 生成随机字符串
func generateRandomString(length int) string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	result := make([]byte, length)
	for i := range result {
		b := make([]byte, 1)
		rand.Read(b)
		result[i] = charset[int(b[0])%len(charset)]
	}
	return string(result)
}

// SensitiveDataMask 敏感数据脱敏工具
var SensitiveDataMask = map[string]string{
	"password":     "***MASKED***",
	"Password":     "***MASKED***",
	"new_password": "***MASKED***",
	"token":        "***MASKED***",
	"Token":        "***MASKED***",
	"access_token": "***MASKED***",
	"refresh_token": "***MASKED***",
	"secret":       "***MASKED***",
	"Secret":       "***MASKED***",
	"private_key":  "***MASKED***",
	"credit_card":  "****-****-****-****",
	"cvv":          "***",
}

// MaskValue 根据键名判断是否需要脱敏
func MaskValue(key, value string) string {
	if _, ok := SensitiveDataMask[key]; ok {
		if len(value) <= 4 {
			return "***"
		}
		return value[:2] + "***" + value[len(value)-2:]
	}
	return value
}
