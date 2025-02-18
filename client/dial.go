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
	network string //network namely tcp, udp, http, mqtt
	address string //address with schema stripped already

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

func checkCommonName(name string, cert *tls.Certificate) error {
	var commonName string
	if cert.Leaf == nil {
		if leaf, err := x509.ParseCertificate(cert.Certificate[0]); err != nil {
			log.Errorf("x509 parser certificate failed, err:%v", err)
			return err
		} else {
			commonName = leaf.Subject.CommonName
		}
	} else {
		commonName = cert.Leaf.Subject.CommonName
	}

	if commonName != name {
		err := fmt.Errorf("the Common Name(%s) in the client certificate does not match the device name(%s)",
			commonName, name)
		log.Errorf("%v", err)
		return err
	}
	return nil
}

func makeSecurityLayerOption(client *LwM2MClient, server *ServerInfo) (coap.PeerOption, error) {
	switch server.securityMode {
	case SecurityModeCertificate:
		cert, err := cipher.LoadKeyAndCertificate(server.secretKey, server.publicKeyOrIdentity)
		if err != nil {
			log.Errorf("load client key and certificate failed, err:%v", err)
			return nil, err
		}

		if err = checkCommonName(client.name, cert); err != nil {
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

		if server.network == coap.UDPBearer {
			dtlsConf := &piondtls.Config{
				Certificates:         []tls.Certificate{*cert},
				ExtendedMasterSecret: piondtls.RequireExtendedMasterSecret,
				RootCAs:              rootCertPool,
				//InsecureSkipVerify:   dtls.InsecureSkipVerify,
			}

			return coap.WithSecurityLayerConfig(coap.SecurityLayerDTLS, dtlsConf), nil
		} else {
			tlsConf := &tls.Config{
				Certificates:       []tls.Certificate{*cert},
				RootCAs:            rootCertPool,
				InsecureSkipVerify: true,
			}

			return coap.WithSecurityLayerConfig(coap.SecurityLayerTLS, tlsConf), nil
		}

	//case SecurityModePreSharedKey:
	//case SecurityModeRawPublicKey:
	case SecurityModeNoSec:
		// nothing todo
		return nil, nil
	default:
		return nil, fmt.Errorf("unsupported security mode:%d", server.securityMode)
	}
}

func dial(client *LwM2MClient, server *ServerInfo) (*MessagerClient, error) {
	var options []coap.PeerOption

	// loads security layer config, which may be nil if security mode is NoSec
	option, err := makeSecurityLayerOption(client, server)
	if err != nil {
		log.Errorf("load security layer config failed: %v", err)
		return nil, err
	}

	if option != nil {
		options = append(options, option)
	}

	messager := NewMessager(client)
	if err = messager.Dial(server.network, server.address, options...); err != nil {
		return nil, err
	}

	messager.Start()

	return messager, nil
}
