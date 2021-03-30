package utils

import (
	"crypto/tls"
	"crypto/x509"
	"io/ioutil"

	"github.com/blackdreamers/core/config"
)

func GetTLSConfig() (*tls.Config, error) {
	cert, err := tls.LoadX509KeyPair(config.Conf.EtcdCertPath, config.Conf.EtcdCertKeyPath)
	if err != nil {
		return nil, err
	}

	ca, err := ioutil.ReadFile(config.Conf.EtcdCaPath)
	if err != nil {
		return nil, err
	}

	pool := x509.NewCertPool()
	pool.AppendCertsFromPEM(ca)

	return &tls.Config{
		Certificates: []tls.Certificate{cert},
		RootCAs:      pool,
	}, nil
}
