package auth_test

import (
	"crypto/rsa"
	"crypto/x509"
	"encoding/hex"
	"encoding/pem"
	"fmt"
	"testing"
	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	jose "gopkg.in/square/go-jose.v2"
	"gopkg.in/square/go-jose.v2/jwt"
)

const privateKeyPEM = `-----BEGIN RSA PRIVATE KEY-----
MIIEpQIBAAKCAQEA6EvnyRuKHSCl48a1E1lEUAzZuemqzWR58LKXW6gsbnnRhlIO
s7U1I14OqJcnZ/2WLbcZL1+imOoKFLJ+1kcUn1PEcnvN+VzcEtcetOT8u61D/w1R
7TwmlIGHjVUiakUCutGyW4hssVkKX/F0VxuoTcRzJdOS+HBCdsxSpk4Do2Vzdire
6XTEyTSCDY8avdHEQl+h7hBW0wPqi5uOqeyToxt4036pJJWKzEcUjjFRlwIsBjXC
mspmogZ5dSMAggXp2g/KzwmiDjFf6UGGEq/W7QPVWQZSAaO749spSYDSUaPXhCnp
IfQjx4Tq31DNvDxfbVbXb2UsRsjzCJ4/achYpwIDAQABAoIBAF7Jilz/hc344ngd
PwoUJNHdTIKWHIIO+8sBtM1LxEWYI60Bkso1qOcznBujAgtm6H7i9b3k8j4fUdm8
NBFPk9Sno8NrTVLxV0PAz+DUh2Q1gUdzsfhx0PRMHnnHQXHHkyOUKYk8E84SHS+l
IOnieTyyaqIGwTPq6tP678a4krd7CSwNhlnssg4Vvk+ARla8njU7NSwrQgTHfe5e
MkwoYpC1aoWXAdXsjJxLaVI9eyfZBPkcH2I1tNk32bCMExEP2CMmp19MpWBpAqk2
gus4EtrXAHkVrtQV6ZR40FDCBQXyXNI/mALjF7lVIAeVzDTiIxYTE7LfYuOF+IlQ
r/T0KdECgYEA9MesFOrn87ew0Vj9fsSau+7gF2dQXWqgmDL5Nt9aAZ+2WMiQUllE
PQNHdwUzf9ocY6qNod878zcAFUITKUozxsv+nSwr18bUGm/K+4Hh55cp3O+YUnrO
isWmp1GNblbrkT5IbgrNLu/kObSfUJqJy5ebvWfxwsrFi8cW93px+XMCgYEA8vG/
iWHUq+YTA0UOsOjRIQG2Y+o/QSWrTbHy2FpLXBtnd+THGG1TnYn4TijVQUAGuZN+
sP5xG/NgHp6RhQhtDd67F94kVhxDYDgbWRhd/PB5p2KwVu4UYGebVTWiuxzc5K4v
CdZ3enA7TSkutH6+UHByinRSwFwgivhNgW77Zv0CgYEA9A4tfgGk2TQaIPEIdp+R
47VNSyhgUXPhwT55ioNxG1NhnO4EL4b/aZHebEYMTfpq+dhwNKf9/wakl410y2NH
cnusVotW/2In0mAKU2/xjFYEkFt3VS1Kx8Q/4G2IhS42227tCoLMh3L566syeUxL
/WkhB90eiGhiRHZxaLh18Q8CgYEAikNPUK6ezJ4KIAhDTieSLYi99qswCLGZhoRJ
wxvQW8E056UKMjq2JaiJ6mGOzK3VpfPtXGnSkae3AnYYN3AOMKCcNf81CtTW+4Gy
/sfBZdyuP7cIyNCCRENywdHepULN7E+9cYnfQY1yEn3nmM2xHjKA3Y9KMTO9SZn6
8jjpVzkCgYEAyrTJiBJ/6UqiJzjP4mVHD0bziWy1nvr8gzE4jJ4NmxU9fCMcFrBk
ZViCq9g/2eWGFWsaQUQDRLVdg250sQN8dwE9wDpFK6+Qpo9RswR3HC2GSArrjvZC
ARrcccXc/A1zWGGb0iPcMZhENRBN9+TWEGbuFT7Wu0/nd55bLOMLJAk=
-----END RSA PRIVATE KEY-----`

const privateKey2PEM = `-----BEGIN RSA PRIVATE KEY-----
MIIEowIBAAKCAQEA5cuz4AmE812LSImRk+fI05BRb9LeiiHGS1ENeexXcyruQaZp
u0xSeQ+PX9XjMiyI5SzkECZkLmdFvkD1Cy+G8DGwA34JZotoKUnrKVrHAV73fe0U
NraK5Dj9iwVx4EIcvbQfidhq3tXp0EumQZLGjpyDXdAkfJSmnElYLgNDGw1GbFap
VXL6lM5KgyKzkEailNOIB25jgjeuLwVxboyXCerQW7IqL2hZnfDTqUvFo7fPDbRI
DO2AnFNjHuLH5nveSV4hxvByyRM8eewVNSaKcGlcJi+snz/Vi7XP4Y+wcrCjmJRY
BuWaHn54Y/S/j5vwMfX0eBm8F1QQ1YcntfziPwIDAQABAoIBAAxjT94bBUTNXS3a
5LMxgp0NIabCmmad7X+v7ecNu6UkwlVdsEN9mNCX6yXCdQ4GYpbgNac9OpdZz+Oy
wsMIm+Ck/RUjHUSe7U2Ug43mK+ZCBVuPhVBxxMkK3Xg6IepyBfSgGjcnKJO8um8V
NPBCBlw3yckr6Fui89xnA87vNBXoJiVeRe/8paEOfMZs9DAIJDDYQvzOSd0wiuP3
Ve9zPgIBxFuswJKFGMl/vY4byBzJuC+GKRmB79XOegerSoV5QTUYRuk5vOANHc7d
IKk8xTD5XnIYbsETiCgBm+ts6OB6NvgkFunedpty+UW3GSKkpurjL0QRuXmXmgzF
e8o9DjkCgYEA9kVOvP3FsSDS23OB5AGl4XDyocqW6VMWtUwHj9S4HF8qgbeKztq2
RJR3OWE/KfRzy/6vrdmtOA5gEEN7u9liv82YPtjFeGswyaOvW8BZQXVKk96txh8Z
pIq+EGQ0al1ZJCqXajW0scFYrJlld9VBZcwsYpP8eImweQIkRU1TyAsCgYEA7t/F
0V5fkMZcQIxZ1Az5SLaTmNiJm2iPUo4Nd5ug7VpDy4eYvPyb0JeHunbn////d45g
girqZNUVUa/nSXEJRgVpSFZj7IrPZuHoanLi/e5mMni3YvMMqzRJRT9qd+hWc3Kx
kSTtVzvGK2NZq72zxTIdwcS9bsnrf7341dIsSx0CgYByWvyVBcIm3fcLsDdAiQNe
C/Se7FPnRI3m4cchIsXbZtV2JqRuKWE5tzcljeKmuLyMnVc2gz3MKeCxrKRoNimE
pxNrG32WzS96cmebU1Ye7zgSMfS/avGdVk+rjNxKB8683Ioy531gjUd/3jsfygb0
Hjr+C3nQ/x7TEguFosKkwQKBgAtnobE6WUO3RMZMLSnDqM9A8FEW3ZMO7fDaGWiB
hLBwY9Y+1hsH0ISoB3HupWsClPbnVFJCrEg+KDNrO5a1D+VI8triTQkJI5fc51TV
wWKwVC7Ktq7BvfQanfjxayroa+A9NJ8ibTaCAxclOi3J8+BRYTxUIVs9xsGll1DW
JQk9AoGBAIjushDIoTELY/s0Q4ReZHuhXkVJLqVREgUBAunk6Ooh6YnqitUwYugv
kUJPbNNtqr5ED56nL/iSWhV9Xfnuf9J5h5663eNVhw6f/b8U6U/sGPRUx91PAIqb
7XsAfC/HqIRNC8Cj6isHjFE2tov4BcicrMACPtWAXDVM7jGQ5muB
-----END RSA PRIVATE KEY-----`

type dumbRandReader struct{}

func (drr dumbRandReader) Read(data []byte) (int, error) {
	for i := range data {
		data[i] = 'a'
	}
	return len(data), nil
}

func TestAuth(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Auth Suite")
}

func getPrivateKey() *rsa.PrivateKey {
	block, extra := pem.Decode([]byte(privateKeyPEM))
	if expected := "RSA PRIVATE KEY"; block.Type != expected {
		message := fmt.Sprintf("unexpected PEM type: got %s, expected %s", block.Type, expected)
		Fail(message)
	}
	ExpectWithOffset(1, extra).To(HaveLen(0))

	privateKey, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	ExpectWithOffset(1, err).NotTo(HaveOccurred())

	return privateKey
}

func getPublicKey() *rsa.PublicKey {
	return &getPrivateKey().PublicKey
}

func getPrivateKey2() *rsa.PrivateKey {
	block, extra := pem.Decode([]byte(privateKey2PEM))
	if expected := "RSA PRIVATE KEY"; block.Type != expected {
		message := fmt.Sprintf("unexpected PEM type: got %s, expected %s", block.Type, expected)
		Fail(message)
	}
	ExpectWithOffset(1, extra).To(HaveLen(0))

	privateKey, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	ExpectWithOffset(1, err).NotTo(HaveOccurred())

	return privateKey
}

func getSecret() []byte {
	return []byte("tuna-fish-marlin")
}

func getWrongSecret() []byte {
	return []byte("logan sat on zelda's beanbag chair")
}

func generateValidClaims() jwt.Claims {
	now := time.Now()
	claims := jwt.Claims{
		Issuer:    "anwork",
		Subject:   "andrew",
		Expiry:    jwt.NewNumericDate(now.Add(time.Second * 1)),
		NotBefore: jwt.NewNumericDate(now),
		IssuedAt:  jwt.NewNumericDate(now),
		ID:        hex.EncodeToString([]byte("aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa")), // 32 a's, see dumbRandReader
	}
	return claims
}

func generateValidToken(secret []byte) string {
	return generateValidTokenWithClaims(secret, generateValidClaims())
}

func generateValidTokenWithClaims(secret []byte, claims jwt.Claims) string {
	signingKey := jose.SigningKey{Algorithm: jose.HS512, Key: secret}
	signerOptions := (&jose.SignerOptions{}).WithType("JWT")
	signer, err := jose.NewSigner(signingKey, signerOptions)
	ExpectWithOffset(1, err).NotTo(HaveOccurred())

	token, err := jwt.Signed(signer).Claims(claims).CompactSerialize()
	ExpectWithOffset(1, err).NotTo(HaveOccurred())

	return token
}

func generateEncryptedToken(publicKey *rsa.PublicKey, secret []byte) string {
	claims := generateValidClaims()

	signingKey := jose.SigningKey{Algorithm: jose.HS512, Key: secret}
	signerOptions := (&jose.SignerOptions{}).WithType("JWT")
	signer, err := jose.NewSigner(signingKey, signerOptions)
	ExpectWithOffset(1, err).NotTo(HaveOccurred())

	encrypterOptions := (&jose.EncrypterOptions{}).WithType("JWT").WithContentType("JWT")
	enc, err := jose.NewEncrypter(
		jose.A256GCM,
		jose.Recipient{
			Algorithm: jose.RSA_OAEP_256,
			Key:       publicKey,
		},
		encrypterOptions,
	)
	ExpectWithOffset(1, err).NotTo(HaveOccurred())

	token, err := jwt.SignedAndEncrypted(signer, enc).Claims(claims).CompactSerialize()
	ExpectWithOffset(1, err).NotTo(HaveOccurred())

	return token
}

func generateEncryptedValidTokenWithClaims(
	publicKey *rsa.PublicKey,
	secret []byte,
	claims jwt.Claims,
) string {
	signingKey := jose.SigningKey{Algorithm: jose.HS512, Key: secret}
	signerOptions := (&jose.SignerOptions{}).WithType("JWT")
	signer, err := jose.NewSigner(signingKey, signerOptions)
	ExpectWithOffset(1, err).NotTo(HaveOccurred())

	encrypterOptions := (&jose.EncrypterOptions{}).WithType("JWT").WithContentType("JWT")
	enc, err := jose.NewEncrypter(
		jose.A256GCM,
		jose.Recipient{
			Algorithm: jose.RSA_OAEP_256,
			Key:       publicKey,
		},
		encrypterOptions,
	)
	ExpectWithOffset(1, err).NotTo(HaveOccurred())

	token, err := jwt.SignedAndEncrypted(signer, enc).Claims(claims).CompactSerialize()
	ExpectWithOffset(1, err).NotTo(HaveOccurred())

	return token
}

func parseClaims(token string, privateKey *rsa.PrivateKey, secret []byte) jwt.Claims {
	parsed, err := jwt.ParseSignedAndEncrypted(token)
	ExpectWithOffset(1, err).NotTo(HaveOccurred())

	nested, err := parsed.Decrypt(privateKey)
	ExpectWithOffset(1, err).NotTo(HaveOccurred())

	claims := jwt.Claims{}
	err = nested.Claims(secret, &claims)
	ExpectWithOffset(1, err).NotTo(HaveOccurred())

	return claims
}
