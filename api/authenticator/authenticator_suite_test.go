package authenticator_test

import (
	"crypto/rsa"
	"crypto/x509"
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

type dumbRandReader struct{}

func (drr dumbRandReader) Read(data []byte) (int, error) {
	as := "aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa"
	copy(data, []byte(as))
	return len(as), nil
}

func TestAuthenticator(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Authenticator Suite")
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
		ID:        "aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa", // 32 a's, see dumbRandReader
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
