package client

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	piondtls "github.com/pion/dtls/v2"
	log "github.com/sirupsen/logrus"
	"github.com/zourva/lwm2m/coap"
	. "github.com/zourva/lwm2m/core"
	"github.com/zourva/pareto/cipher"
)

type ServerInfo struct {
	address string

	// securityMode
	// Determines which security mode is used
	// - 0: PreShared Key mode
	// - 1: Raw Public Key mode
	// - 2: Certificate mode
	// - 3: NoSec mode
	// - 4: Certificate mode with EST
	securityMode int

	// publicKeyOrIdentity
	// Stores the LwM2M Client's certificate, public key (RPK mode) or PSK Identity (PSK mode).
	// securityMode is 2 : client certificate file
	publicKeyOrIdentity []byte

	// serverPublicKey
	// Stores the LwM2M Server's, respectively LwM2M
	// Bootstrap-Server's, certificate, public key (RPK mode) or trust anchor. The Certificate Mode
	// Resource determines the content of this resource.
	serverPublicKey []byte

	// secretKey
	// Stores the secret key (PSK mode) or private key(RPK or certificate mode).
	// securityMode is 2 : client private key
	secretKey []byte
}

func loadDTLSConfig(server *ServerInfo) (*piondtls.Config, error) {
	var dtlsConf *piondtls.Config

	switch server.securityMode {
	case SecurityModeCertificate:
		cert, err := cipher.LoadKeyAndCertificate(server.secretKey, server.publicKeyOrIdentity)
		if err != nil {
			log.Errorf("load client key and certificate failed, err:%v", err)
			return nil, err
		}

		log.Debugf("load client key and certificate certificate file successfully")

		var rootCertPool *x509.CertPool
		if len(server.serverPublicKey) != 0 {
			rootCertPool, err = cipher.LoadAllCertPool([]string{string(server.serverPublicKey)})
			if err != nil {
				log.Errorf("load root certificate failed, err:%v", err)
				return nil, err
			}

			log.Debugf("load root certificate file successfully")
		}

		dtlsConf = &piondtls.Config{
			Certificates:         []tls.Certificate{*cert},
			ExtendedMasterSecret: piondtls.RequireExtendedMasterSecret,
			RootCAs:              rootCertPool,
			//InsecureSkipVerify:   dtls.InsecureSkipVerify,
		}
	//case SecurityModePreSharedKey:
	//case SecurityModeRawPublicKey:
	case SecurityModeNoSec:
		// nothing todo
		break
	default:
		return nil, fmt.Errorf("unsupported security mode:%d", server.securityMode)
	}

	return dtlsConf, nil
}

func dial(client *LwM2MClient, server *ServerInfo) (*MessagerClient, error) {
	var dtlsConf *piondtls.Config
	var err error

	if dtlsConf, err = loadDTLSConfig(server); err != nil {
		log.Errorf("load dtls config failed: %v", err)
		return nil, err
	}

	messager := NewMessager(client)
	if err = messager.Dial(server.address, coap.WithDTLSConfig(dtlsConf)); err != nil {
		return nil, err
	}

	messager.Start()

	return messager, nil
}
