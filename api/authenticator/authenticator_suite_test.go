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
MIIEogIBAAKCAQEAmSDDkhD63nQLhEQ+n+OAdYuWzeJrs5qCn3S7UzK+9KHUV0NG
FSLE5irhVw80ArEV3Csc2wLckA8NUN4+86W/xFPK4zipppizdkDD/GYOOgvrsp7f
IMqGzALsOQnptUe6CZPydCSlkVItiHueQPOPjXNKqfdNsvE9maCCggVcloMQBan1
uuNs3Wl1J2I83edqrU4iXSnpgL9pbq/0fhmBbk50sjxm+Mwt5ZOq3Wv0d8MKZobS
/4NXkg9sCT3AzJZij57ft7abq/6z2QA1gwCLFzPbVJvc8KWiKCi96Fp9PVvR5R4W
QztZIQRv986wjt7T+BEIRa6UsES4CBYdVYVedQIDAQABAoIBAFENuLeeaUxK3LW8
cfGsOJj/tGlyilPdW3sQDP/zAoT3DDDMmVIiv1qeI/0zPPKXzdxmrbV4BEv0y7Wc
jnHlsGY5fFFd8t8OQSA9FACL/MfY+3/m/HCdA1EF0wg2KREd0Gm1eEbmBwWvHA78
cD4tLjVPa5wgHW60p0ikX7B3KhU4pS1akHImeFrA04Q2E/WnpaQ220tPdiJiErLa
GKnowzhqqgdIOmtBPHeSCu8xus6Deduxx9QPsMJTX6hLOri6FTh3q0i8d8nINK1e
7UXXvtVfJVAenDVDxw2Xa94yyIaEqyZg2xBXUnViqw0dpqN0WZMkSXbri3cTz3ZI
IiNySsECgYEAxpdTpunsGPfdxQLkP3S1KlyrHQmL0WZ0aO1JCkDr7U72JRax1sEF
slWmVKEt0HNSyt/OsEV+AIZgCS37HXBxMDTKRCJfPd44x+y9iEVb92GxmSm+Y5D2
N9gVqLKilnqcIIbufHdWvjGFzDBbaa/dwL3ip5QjTEG2v5k9rYaGr80CgYEAxWTz
vyxsN7ZOpEDvsborL8zkQrHOSY/wE6Q2nIQ4AU2Saj+eVAiEwvY0i/b8ujxUcFsq
0aPmdR/+L6iwO7DjImzrOIfeLIkRkwfgxgfeA+LdQl1fPEVttF1mEEh1Nqw3oRw2
O1fHXHRGmHIQBbrs2Ob76URkQarhk++urfpTMUkCgYBdLdfEM3Hh7TsTG145H+t5
Ku3mu/nskKQCL4Pb2cZZHHmFwXZEC2E+4c0fXFkAu3uXURfLwW9zk5kv9XEjyQRy
1/InsfD6OHBv7faoH9Mc6avI77szQGsmnK7c7qQ28uSapnTz1ZLPvrDGs9HQbwIf
U3kro+hD+XljwOUdrEc56QKBgGkUvM1gkJt36ZV3HCK4wJTJthnrHa37egp6uLfJ
iybmbLMy+s4xWruO5Zo8+X0K8Mh/P+QqJFzlkyM74dVk7QU+hlmpupRqw3hKEVZa
ic33z2Gs4y3qp+QEdHjmmb9pHpQduKEOsYp+O7abwfvK7prpG+GqDtUGat+eEJd9
UxxBAoGALitSRa78uPo6PmAiVNqbGYGiAaZypSydo5Wyk+Q+o9c7ze3tWfl3W+yw
Jc4igRpxAsTTip9H+NEro7C+ENq+VDAaHx1HaKv9X0xGR0b3xa6uFZQEAHlfeWGu
OpxD9rtecQfP7BmOgKmxs8JYnaM4PtFmg2z+JcNBblRGbk2HF40=
-----END RSA PRIVATE KEY-----`

const wrongPrivateKeyPEM = `-----BEGIN RSA PRIVATE KEY-----
MIIEpAIBAAKCAQEA3SgB319uzDgiVIjjubMwSwM5f/QZqha95NnBYsATlC5Mr5u+
nMiPYoGn+1purhbiEv+1VQA1dHq2JHeHQFzEnvX/NeoYFcWSO1suUojaGOuBqJ0P
Xway8Roq+3ALNrr/q1oakyL9HlFoG5Gzts6efZVLO4CDL6I79Zay5JLRJrlwbIIx
QPVjTgacD6xKR7dqtt3ICUnZXr9tei/XDMoSJktzM1VfpZrdyFmO1oyOeMilQaUj
3mnhGpwDiIToshaAB8zHSTljRAOPJWKjNgQ7xUy6selXpM9k9XjCRWzBjJQlKF0i
xN3PFPzUblEiDhGywb6YWxN23OrDtle/UvCd9QIDAQABAoIBAHV77lq6eMKWQ26C
9e7NiScgRGh7xTJ9OE2dfEU4Ym5ClkxnwT/Kr5kV5aX6aXekAl8zZiIMpVkGY0JG
GRyQ64bCidwbSsYuGbvceWQw/SNl4JAxN1w009SPISxHXj2p424kxJ9LZAm6F0Mm
CMFhTALMz5WBdK9WUOHrN4eDkwYtWV27BIIyo1MnH/iWqJaIzZ4fqajmhLa3E3Hk
6YbAXpcS9w8OHwxWFeGN6+nsUoDOVshZGZ4ZVoxCKmA3ZnKFjF6zIXsdzQa/MjJr
Kw0PaxN902NROLR3I3Mu7QrsLJiarRAA49Rqfdix5Fu2JbN4XB6i4aGmz03PI5Yl
En+aLOECgYEA/6iiAU2chmlL4N+ydma5s+4G3hhMUkYVnSEbEjokm/Dnenkachdp
Wu2L+xBEElPnEd2urT46wJvv6HPpTSR7zLoyn1Xvtv/CvpKO7p8vLlpup2WROusN
1z6CbkYoqtrf1R0NMwkZhaTaQy0g68eQEd32dLYsrLS2lhSItLgznF0CgYEA3XOV
dnr9WVq+hGTRgG0w1P3nt6YmDcSCD8EiCcacKFIZudBRVu/9RdOAt/mqTQCHQK1a
/ROfHjfitoX5GRvZ30QL94GGbHHPyjOBYEa5OiH0T+yFUbjaT/BBVgFywVrn3joj
MZnK9itoM4yKKt27GKwzBc75F41BQ5Jggu8PLnkCgYEAqDXb4d+ezREay6pjUWPl
a22BNz/ld3yFXA0cMrHuxGuM4hgsPkUJHLqPD3F0WFq7/hVNiM8Y+QGgp+Eb75XB
nsIj7JIuVsmQ6LKlOHukH2uAwsMg+xMM2EJYrxWaTFAWVbH3rUyfbj85HFnk/z0e
naLdNY1nd3qvZ6+7AqzvyEECgYAQGQLYZgBcqngG77069LUEBqD9fJpvjcVWl9d9
lm5rj+xG0ZnYFAH5PXKx7PgwOMWcMf3XP8HlVHKqifqdlKS10iB8kXHQGEXekPfq
o7l7PFSiKrNWSXW1MeXN9rT80TrhsKA2TtOuKWGdva2diBi9pmbfGTiKOb5wxwc0
/WPBIQKBgQCJrNmZF2r+00EUDyYaqHsv5Y6kbcggtN+hVTs8NsOoiuL9YCzaJBgZ
AvUksHSWUHMF9g5YeIgGSCAzbw46qn8Lqm6yY5GCIMHZwDVW4AHq8YxslRJxHdWh
53DhVUYtn2gPSNKPG63rTV7L+KTEH5Yiw6YY5HjOZ80tqcUPaOSr9A==
-----END RSA PRIVATE KEY-----`

func TestAuthenticator(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Authenticator Suite")
}

func getPrivateKey() *rsa.PrivateKey {
	block, _ := pem.Decode([]byte(privateKeyPEM))
	if expected := "RSA PRIVATE KEY"; block.Type != expected {
		message := fmt.Sprintf("unexpected PEM type: got %s, expected %s", block.Type, expected)
		ExpectWithOffset(1, true).To(BeFalse(), message)
	}

	privateKey, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	ExpectWithOffset(1, err).NotTo(HaveOccurred())

	return privateKey
}

func getWrongPrivateKey() *rsa.PrivateKey {
	block, _ := pem.Decode([]byte(wrongPrivateKeyPEM))
	if expected := "RSA PRIVATE KEY"; block.Type != expected {
		message := fmt.Sprintf("unexpected PEM type: got %s, expected %s", block.Type, expected)
		ExpectWithOffset(1, true).To(BeFalse(), message)
	}

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
	}
	return claims
}

func generateValidToken(privateKey *rsa.PrivateKey, secret []byte) string {
	return generateValidTokenWithClaims(privateKey, secret, generateValidClaims())
}

func generateValidTokenWithClaims(privateKey *rsa.PrivateKey, secret []byte, claims jwt.Claims) string {
	signingKey := jose.SigningKey{Algorithm: jose.HS512, Key: secret}
	signerOptions := (&jose.SignerOptions{}).WithType("JWT")
	signer, err := jose.NewSigner(signingKey, signerOptions)
	ExpectWithOffset(1, err).NotTo(HaveOccurred())

	encrypterOptions := (&jose.EncrypterOptions{}).WithType("JWT").WithContentType("JWT")
	enc, err := jose.NewEncrypter(
		jose.A256GCM,
		jose.Recipient{
			Algorithm: jose.RSA_OAEP_256,
			Key:       &privateKey.PublicKey,
		},
		encrypterOptions,
	)
	ExpectWithOffset(1, err).NotTo(HaveOccurred())

	token, err := jwt.SignedAndEncrypted(signer, enc).Claims(claims).CompactSerialize()
	ExpectWithOffset(1, err).NotTo(HaveOccurred())

	return token
}

func generateUnencryptedToken(secret []byte) string {
	signingKey := jose.SigningKey{Algorithm: jose.HS512, Key: secret}
	signerOptions := (&jose.SignerOptions{}).WithType("JWT")
	signer, err := jose.NewSigner(signingKey, signerOptions)
	ExpectWithOffset(1, err).NotTo(HaveOccurred())

	now := time.Now()
	cl := jwt.Claims{
		Issuer:    "anwork",
		Subject:   "andrew",
		Expiry:    jwt.NewNumericDate(now.Add(time.Second * 1)),
		NotBefore: jwt.NewNumericDate(now),
		IssuedAt:  jwt.NewNumericDate(now),
	}
	token, err := jwt.Signed(signer).Claims(cl).CompactSerialize()
	ExpectWithOffset(1, err).NotTo(HaveOccurred())

	return token
}

func parseClaims(token string, privateKey *rsa.PrivateKey, secret []byte) jwt.Claims {
	parsed, err := jwt.ParseSignedAndEncrypted(token)
	ExpectWithOffset(1, err).NotTo(HaveOccurred())

	nested, err := parsed.Decrypt(privateKey)
	ExpectWithOffset(1, err).NotTo(HaveOccurred())

	out := jwt.Claims{}
	err = nested.Claims(secret, &out)
	ExpectWithOffset(1, err).NotTo(HaveOccurred())

	return out
}
