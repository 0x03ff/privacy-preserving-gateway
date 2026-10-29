package configutils

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/hex"
	"encoding/pem"
	"log/slog"
	"math/big"
	"os"
	"path/filepath"
	"time"
)

// SetupAppKeys ensures ECDSA and AES signing keys exist in var/key/
func SetupAppKeys() error {
	dir := "var/key"
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}

	privPath := filepath.Join(dir, "ECDSA_384_private.pem")
	pubPath := filepath.Join(dir, "ECDSA_384_public.pem")
	aesPath := filepath.Join(dir, "secret.key")

	// Check if keys already exist
	if _, err := os.Stat(privPath); err == nil {
		slog.Info("App keys setup: Keys already exist - no action needed")
		return nil
	}

	slog.Info("App keys setup: Generating new ECDSA and AES keys...")

	// 1. Generate ECDSA P-384 Key Pair
	privateKey, err := ecdsa.GenerateKey(elliptic.P384(), rand.Reader)
	if err != nil {
		return err
	}

	// Marshal private key to PKCS#8 / EC PEM
	derBytes, err := x509.MarshalECPrivateKey(privateKey)
	if err != nil {
		return err
	}
	privPem := pem.EncodeToMemory(&pem.Block{Type: "EC PRIVATE KEY", Bytes: derBytes})
	if err := os.WriteFile(privPath, privPem, 0600); err != nil {
		return err
	}

	// Marshal public key to PEM
	pubDer, err := x509.MarshalPKIXPublicKey(&privateKey.PublicKey)
	if err != nil {
		return err
	}
	pubPem := pem.EncodeToMemory(&pem.Block{Type: "PUBLIC KEY", Bytes: pubDer})
	if err := os.WriteFile(pubPath, pubPem, 0644); err != nil {
		return err
	}

	// 2. Generate AES Symmetric Secret Key (32 bytes hex)
	aesKey := make([]byte, 32)
	if _, err := rand.Read(aesKey); err != nil {
		return err
	}
	if err := os.WriteFile(aesPath, []byte(hex.EncodeToString(aesKey)), 0600); err != nil {
		return err
	}

	slog.Info("App keys setup: ECDSA and AES keys generated successfully")
	return nil
}

// SetupTLSCertificates ensures self-signed TLS certs exist in var/ssl/
func SetupTLSCertificates() error {
	dir := "var/ssl"
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}

	certPath := filepath.Join(dir, "go_cert.pem")
	keyPath := filepath.Join(dir, "go_key.pem")
	fullchainPath := filepath.Join(dir, "fullchain.pem")

	// Check if certificates already exist
	if _, err := os.Stat(fullchainPath); err == nil {
		slog.Info("TLS setup: Certificates already exist - no action needed")
		return nil
	}

	slog.Info("TLS setup: Generating self-signed TLS certificates...")

	// 1. Generate Private Key for Server
	priv, err := ecdsa.GenerateKey(elliptic.P384(), rand.Reader)
	if err != nil {
		return err
	}

	// 2. Define Certificate Template
	serialNumberLimit := new(big.Int).Lsh(big.NewInt(1), 128)
	serialNumber, err := rand.Int(rand.Reader, serialNumberLimit)
	if err != nil {
		return err
	}

	template := x509.Certificate{
		SerialNumber: serialNumber,
		Subject: pkix.Name{
			Organization: []string{"capstone"},
			CommonName:   "privary",
		},
		NotBefore:             time.Now(),
		NotAfter:              time.Now().AddDate(1, 0, 0), // Valid for 1 year
		KeyUsage:              x509.KeyUsageKeyEncipherment | x509.KeyUsageDigitalSignature,
		ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth},
		BasicConstraintsValid: true,
	}

	// 3. Create Self-Signed Certificate
	certBytes, err := x509.CreateCertificate(rand.Reader, &template, &template, &priv.PublicKey, priv)
	if err != nil {
		return err
	}

	// 4. Write Server Certificate (go_cert.pem)
	certOut, err := os.Create(certPath)
	if err != nil {
		return err
	}
	pem.Encode(certOut, &pem.Block{Type: "CERTIFICATE", Bytes: certBytes})
	certOut.Close()

	// 5. Write Server Private Key (go_key.pem)
	keyOut, err := os.Create(keyPath)
	if err != nil {
		return err
	}
	privBytes, err := x509.MarshalECPrivateKey(priv)
	if err != nil {
		return err
	}
	pem.Encode(keyOut, &pem.Block{Type: "EC PRIVATE KEY", Bytes: privBytes})
	keyOut.Close()

	// 6. Write Fullchain
	fullchainData, err := os.ReadFile(certPath)
	if err != nil {
		return err
	}
	if err := os.WriteFile(fullchainPath, fullchainData, 0644); err != nil {
		return err
	}

	slog.Info("TLS setup: Certificates generated successfully")
	return nil
}