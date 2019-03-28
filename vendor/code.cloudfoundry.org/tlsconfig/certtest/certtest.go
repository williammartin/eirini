// Package certtest can be used to build a PKI for test purposes. The
// certificates generated by this package should not be used for production or
// other sensitive traffic.
package certtest

import (
	"crypto/tls"
	"crypto/x509"
	"net"
	"time"

	"github.com/square/certstrap/pkix"
)

const (
	o        = "certtest Organization"
	ou       = "certtest Unit"
	country  = "AQ"
	province = "Ross Island"
	city     = "McMurdo Station"

	// This is nowhere near enough bits for a real certificate but creating
	// larger keys on each test run takes too long.
	//
	// Do not use these certificates to transport secrets.
	keySize = 1024
)

// Authority represents a Certificate Authority. It should not be used for
// anything except ephemeral test usage.
type Authority struct {
	cert *pkix.Certificate
	key  *pkix.Key
}

// BuildCA creates a new test Certificate Authority. The name argument can be
// used to distinguish between multiple authorities.
func BuildCA(name string) (*Authority, error) {
	key, err := pkix.CreateRSAKey(keySize)
	if err != nil {
		return nil, err
	}

	// XXX: Add a month so CA expires after its certificates.
	expiry := time.Now().AddDate(1, 1, 0)

	crt, err := pkix.CreateCertificateAuthority(key, ou, expiry, o, country, province, city, name)
	if err != nil {
		return nil, err
	}

	return &Authority{
		cert: crt,
		key:  key,
	}, nil
}

// SignOption is used to alter the signed certificate parameters.
type SignOption func(*signOptions)

// WithIPs adds the passed IPs to be valid for the requested certificate.
func WithIPs(ips ...net.IP) SignOption {
	return func(options *signOptions) {
		options.ips = ips
	}
}

// WithDomains adds the passed domains to be valid for the requested
// certificate.
func WithDomains(domains ...string) SignOption {
	return func(options *signOptions) {
		options.domains = domains
	}
}

// BuildSignedCertificate creates a new signed certificate which is valid for
// `localhost` and `127.0.0.1` by default. This can be changed by passing in
// the various options. The certificates it creates should only be used
// ephemerally in tests.
func (a *Authority) BuildSignedCertificate(name string, options ...SignOption) (*Certificate, error) {
	key, err := pkix.CreateRSAKey(keySize)
	if err != nil {
		return nil, err
	}

	opts := defaultSignOptions()
	for _, o := range options {
		opts.apply(o)
	}

	csr, err := pkix.CreateCertificateSigningRequest(key, ou, opts.ips, opts.domains, nil, o, country, province, city, name)
	if err != nil {
		return nil, err
	}

	expiry := time.Now().AddDate(1, 0, 0)

	crt, err := pkix.CreateCertificateHost(a.cert, a.key, csr, expiry)
	if err != nil {
		return nil, err
	}

	return &Certificate{
		cert: crt,
		key:  key,
	}, nil
}

// CertificatePEM returns the authorities certificate as a PEM encoded bytes.
func (a *Authority) CertificatePEM() ([]byte, error) {
	return a.cert.Export()
}

// CertPool returns a certificate pool which is pre-populated with the
// Certificate Authority.
func (a *Authority) CertPool() (*x509.CertPool, error) {
	cert, err := a.CertificatePEM()
	if err != nil {
		return nil, err
	}

	pool := x509.NewCertPool()
	pool.AppendCertsFromPEM(cert)

	return pool, nil
}

// Certificate represents a Certificate which has been signed by a Certificate
// Authority.
type Certificate struct {
	cert *pkix.Certificate
	key  *pkix.Key
}

// TLSCertificate returns the certificate as Go standard library
// tls.Certificate.
func (c *Certificate) TLSCertificate() (tls.Certificate, error) {
	certBytes, err := c.cert.Export()
	if err != nil {
		return tls.Certificate{}, nil
	}

	keyBytes, err := c.key.ExportPrivate()
	if err != nil {
		return tls.Certificate{}, nil
	}

	return tls.X509KeyPair(certBytes, keyBytes)
}

// CertificatePEMAndPrivateKey returns the certificate as a PEM encoded bytes and the private key bytes.
func (c *Certificate) CertificatePEMAndPrivateKey() ([]byte, []byte, error) {
	certBytes, err := c.cert.Export()
	if err != nil {
		return nil, nil, err
	}

	keyBytes, err := c.key.ExportPrivate()
	if err != nil {
		return nil, nil, err
	}

	return certBytes, keyBytes, nil
}

type signOptions struct {
	domains []string
	ips     []net.IP
}

func defaultSignOptions() *signOptions {
	return &signOptions{
		domains: []string{"localhost"},
		ips:     []net.IP{net.ParseIP("127.0.0.1")},
	}
}

func (s *signOptions) apply(option SignOption) {
	option(s)
}