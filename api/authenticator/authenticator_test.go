package authenticator_test

import (
	"crypto/rsa"
	"fmt"
	"time"

	"code.cloudfoundry.org/clock/fakeclock"
	"github.com/ankeesler/anwork/api/authenticator"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"gopkg.in/square/go-jose.v2/jwt"
)

var _ = Describe("Authenticator", func() {
	var (
		a *authenticator.Authenticator

		clock *fakeclock.FakeClock

		privateKey *rsa.PrivateKey
		secret     []byte
	)

	BeforeEach(func() {
		clock = fakeclock.NewFakeClock(time.Now())

		privateKey = getPrivateKey()
		secret = getSecret()

		a = authenticator.New(clock, privateKey, secret)
	})

	Describe("Authenticate", func() {
		It("validates the 'Authorization' header holds a real token", func() {
			validToken := generateValidToken(privateKey, secret)
			Expect(a.Authenticate(validToken)).To(Succeed())
		})

		Context("on unencrypted token", func() {
			It("returns an error", func() {
				unencryptedToken := generateUnencryptedToken(secret)

				err := a.Authenticate(unencryptedToken)
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(HavePrefix("could not parse token:"))
			})
		})

		Context("on token encrypted with the wrong key", func() {
			It("returns an error", func() {
				wrongPrivateKey := getWrongPrivateKey()
				wrongKeyToken := generateValidToken(wrongPrivateKey, secret)

				err := a.Authenticate(wrongKeyToken)
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(HavePrefix("could not decrypt token:"))
			})
		})

		Context("on token signed with the wrong secret", func() {
			It("returns an error", func() {
				wrongSecret := getWrongSecret()
				wrongKeyToken := generateValidToken(privateKey, wrongSecret)

				err := a.Authenticate(wrongKeyToken)
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(HavePrefix("could not get claims"))
			})
		})

		Context("on invalid claims", func() {
			testInvalidClaim := func(which string, invalidateClaimsFunc func(*jwt.Claims)) {
				Context(fmt.Sprintf("%s is wrong", which), func() {
					It("returns an error", func() {
						claims := generateValidClaims()
						invalidateClaimsFunc(&claims)
						token := generateValidTokenWithClaims(privateKey, secret, claims)

						err := a.Authenticate(token)
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

	Describe("Token", func() {
		It("returns an encrypted and signed token with the correct claims", func() {
			token, err := a.Token()
			Expect(err).NotTo(HaveOccurred())

			claims := parseClaims(token, privateKey, secret)
			Expect(claims).To(Equal(jwt.Claims{
				Issuer:    "anwork",
				Subject:   "andrew",
				Expiry:    jwt.NewNumericDate(clock.Now().Add(time.Second * 1)),
				NotBefore: jwt.NewNumericDate(clock.Now()),
				IssuedAt:  jwt.NewNumericDate(clock.Now()),
			}))
		})
	})
})
