package cmd

import (
	"crypto/tls"
	"crypto/x509"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"time"

	"github.com/cafebazaar/booker-resources/common"
)

var (
	_keyPair  *tls.Certificate
	_certPool *x509.CertPool
	startTime time.Time
)

func init() {
	startTime = time.Now()
}

func keyPair() (*tls.Certificate, error) {
	if _keyPair == nil {
		pair, err := tls.LoadX509KeyPair(common.ConfigString("CERT_FILE"), common.ConfigString("KEY_FILE"))
		if err != nil {
			return nil, fmt.Errorf("Failed to load tls key pair: %s", err)
		}

		_keyPair = &pair
	}

	return _keyPair, nil
}

func certPool() (*x509.CertPool, error) {
	if _certPool == nil {
		newCertPool := x509.NewCertPool()

		caContent, err := ioutil.ReadFile(common.ConfigString("CA_FILE"))
		if err != nil {
			return nil, fmt.Errorf("Failed to load tls ca file: %s", err)
		}

		ok := newCertPool.AppendCertsFromPEM(caContent)
		if !ok {
			return nil, errors.New("Failed to append tls ca certs")
		}

		_certPool = newCertPool
	}
	return _certPool, nil
}

// grpcHandlerFunc returns an http.Handler that returns healthz info on /healthz,
// or delegates to otherHandler otherwise. Copied from cockroachdb and modified.
func healthzAddedHandlerFunc(otherHandler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if strings.HasPrefix(r.URL.Path, "/healthz") {
			w.Header().Set("Content-Type", "text/plain")
			uptime := time.Now().Sub(startTime)
			w.Write([]byte(fmt.Sprintf("OK\nVersion: %s\nBuild Time: %s\nUptime: %s",
				common.Version, common.BuildTime, uptime)))
		} else {
			otherHandler.ServeHTTP(w, r)
		}
	})
}
