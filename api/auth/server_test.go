package auth_test

import (
	"crypto/rsa"
	"encoding/hex"
	"fmt"
	"time"

	"code.cloudfoundry.org/clock/fakeclock"
	"github.com/ankeesler/anwork/api/auth"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"gopkg.in/square/go-jose.v2/jwt"
)

var _ = Describe("Server", func() {
	var (
		s *auth.Server

		clock *fakeclock.FakeClock

		publicKey *rsa.PublicKey
		secret    []byte
	)

	BeforeEach(func() {
		rand := dumbRandReader{}

		clock = fakeclock.NewFakeClock(time.Now())

		publicKey = getPublicKey()
		secret = getSecret()

		s = auth.NewServer(clock, rand, publicKey, secret)
	})

	Describe("Authenticate", func() {
		Context("when Token() has been called", func() {
			BeforeEach(func() {
				_, err := s.Token() // generate token to set currentJTI
				Expect(err).NotTo(HaveOccurred())
			})

			It("returns no error on a valid token", func() {
				validToken := generateValidToken(secret)
				Expect(s.Authenticate(validToken)).To(Succeed())
			})
		})

		Context("when Token() has not been called yet", func() {
			It("returns an error", func() {
				validToken := generateValidToken(secret)
				err := s.Authenticate(validToken)
				Expect(err).To(MatchError("null jti"))
			})
		})

		Context("on encrypted token", func() {
			It("returns an error", func() {
				unencryptedToken := generateEncryptedToken(publicKey, secret)

				err := s.Authenticate(unencryptedToken)
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(HavePrefix("could not parse token:"))
			})
		})

		Context("on token signed with the wrong secret", func() {
			It("returns an error", func() {
				wrongSecret := getWrongSecret()
				wrongKeyToken := generateValidToken(wrongSecret)

				err := s.Authenticate(wrongKeyToken)
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(HavePrefix("could not get claims"))
			})
		})

		Context("on invalid claims", func() {
			testInvalidClaim := func(which string, invalidateClaimsFunc func(*jwt.Claims)) {
				Context(fmt.Sprintf("%s is wrong", which), func() {
					It("returns an error", func() {
						_, err := s.Token() // generate token to set currentJTI
						Expect(err).NotTo(HaveOccurred())

						claims := generateValidClaims()
						invalidateClaimsFunc(&claims)
						token := generateValidTokenWithClaims(secret, claims)

						err = s.Authenticate(token)
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

			testInvalidClaim("ID claim (jti)", func(claims *jwt.Claims) {
				claims.ID = "bbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbb"
			})
		})
	})

	Describe("Token", func() {
		It("returns an encrypted and signed token with the correct claims", func() {
			token, err := s.Token()
			Expect(err).NotTo(HaveOccurred())

			privateKey := getPrivateKey()
			claims := parseClaims(token, privateKey, secret)
			Expect(claims).To(Equal(jwt.Claims{
				Issuer:    "anwork",
				Subject:   "andrew",
				Expiry:    jwt.NewNumericDate(clock.Now().Add(time.Hour)),
				NotBefore: jwt.NewNumericDate(clock.Now()),
				IssuedAt:  jwt.NewNumericDate(clock.Now()),
				ID:        hex.EncodeToString([]byte("aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa")), // see dumbRandReader
			}))
		})
	})
})
