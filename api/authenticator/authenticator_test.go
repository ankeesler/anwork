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

		publicKey *rsa.PublicKey
		secret    []byte
	)

	BeforeEach(func() {
		rand := dumbRandReader{}

		clock = fakeclock.NewFakeClock(time.Now())

		publicKey = getPublicKey()
		secret = getSecret()

		a = authenticator.New(clock, rand, publicKey, secret)
	})

	Describe("Authenticate", func() {
		It("returns no error on a valid token", func() {
			validToken := generateValidToken(secret)
			Expect(a.Authenticate(validToken)).To(Succeed())
		})

		Context("on empty token", func() {
			It("fails with an error message of sorts", func() {
				err := a.Authenticate("")
				Expect(err).To(HaveOccurred())
				// TODO
			})
		})

		Context("the token is incorrectly formatted", func() {
			It("fails with an error message of sorts", func() {
				err := a.Authenticate("bearerasdfasdfasdf")
				Expect(err).To(HaveOccurred())
				// TODO
			})
		})

		Context("the token is not of type bearer", func() {
			It("fails with an error message of sorts", func() {
				err := a.Authenticate("tuna asdfasdfasdf")
				Expect(err).To(HaveOccurred())
				// TODO
			})
		})

		Context("on encrypted token", func() {
			It("returns an error", func() {
				unencryptedToken := generateEncryptedToken(publicKey, secret)

				err := a.Authenticate(unencryptedToken)
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(HavePrefix("could not parse token:"))
			})
		})

		Context("on token signed with the wrong secret", func() {
			It("returns an error", func() {
				wrongSecret := getWrongSecret()
				wrongKeyToken := generateValidToken(wrongSecret)

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
						token := generateValidTokenWithClaims(secret, claims)

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

			privateKey := getPrivateKey()
			claims := parseClaims(token, privateKey, secret)
			Expect(claims).To(Equal(jwt.Claims{
				Issuer:    "anwork",
				Subject:   "andrew",
				Expiry:    jwt.NewNumericDate(clock.Now().Add(time.Hour)),
				NotBefore: jwt.NewNumericDate(clock.Now()),
				IssuedAt:  jwt.NewNumericDate(clock.Now()),
			}))
		})
	})
})
