package auth_test

import (
	"crypto/rsa"
	"fmt"
	"time"

	"code.cloudfoundry.org/clock/fakeclock"
	"github.com/ankeesler/anwork/api/auth"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"gopkg.in/square/go-jose.v2/jwt"
)

var _ = Describe("Client", func() {
	var (
		c *auth.Client

		clock *fakeclock.FakeClock

		privateKey *rsa.PrivateKey
		secret     []byte
	)

	BeforeEach(func() {
		clock = fakeclock.NewFakeClock(time.Now())

		privateKey = getPrivateKey()
		secret = getSecret()

		c = auth.NewClient(clock, privateKey, secret)
	})

	It("decrypts a token, verifies it, and then signs it again", func() {
		encryptedToken := generateEncryptedToken(&privateKey.PublicKey, secret)
		decryptedToken, err := c.Validate(encryptedToken)
		Expect(err).NotTo(HaveOccurred())
		Expect(decryptedToken).To(Equal(generateValidToken(secret)))
	})

	Context("when the token is bogus", func() {
		It("returns an error", func() {
			_, err := c.Validate("marshmallow is on the counter")
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(HavePrefix("could not parse token"))
		})
	})

	Context("when the token is encrypted with the wrong key", func() {
		It("returns an error", func() {
			privateKey2 := getPrivateKey2()
			encryptedToken := generateEncryptedToken(&privateKey2.PublicKey, secret)
			_, err := c.Validate(encryptedToken)
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(HavePrefix("could not decrypt token"))
		})
	})

	Context("when the token is signed with the wrong secret", func() {
		It("returns an error", func() {
			secret2 := []byte("she crossed the street")
			encryptedToken := generateEncryptedToken(&privateKey.PublicKey, secret2)
			_, err := c.Validate(encryptedToken)
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(HavePrefix("could not verify claims"))
		})
	})

	Context("when the token has invalid claims", func() {
		testInvalidClaim := func(which string, invalidateClaimsFunc func(*jwt.Claims)) {
			Context(fmt.Sprintf("%s is wrong", which), func() {
				It("returns an error", func() {
					claims := generateValidClaims()
					invalidateClaimsFunc(&claims)
					token := generateEncryptedValidTokenWithClaims(
						&privateKey.PublicKey, secret, claims)

					_, err := c.Validate(token)
					Expect(err).To(HaveOccurred())
					Expect(err.Error()).To(HavePrefix("invalid claims"))
					Expect(err.Error()).To(ContainSubstring(which))
				})
			})
		}

		testInvalidClaim("issuer claim (iss)", func(claims *jwt.Claims) {
			claims.Issuer = "wrong-issuer"
		})

		testInvalidClaim("subject claim (sub)", func(claims *jwt.Claims) {
			claims.Subject = "wrong-subject"
		})

		testInvalidClaim("not valid yet (nbf)", func(claims *jwt.Claims) {
			claims.NotBefore = jwt.NewNumericDate(time.Now().Add(time.Hour * 24))
		})

		testInvalidClaim("expired (exp)", func(claims *jwt.Claims) {
			claims.Expiry = jwt.NewNumericDate(time.Now().Add(time.Hour * -24))
		})
	})
})
