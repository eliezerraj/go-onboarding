package configuration

import(
	"os"
	"encoding/base64"

	"github.com/joho/godotenv"
	"github.com/go-onboarding/internal/core/model"
)

// Load the certs TLS
func GetCertEnv() model.Cert {
	childLogger.Info().Str("func","GetCertEnv").Send()

	err := godotenv.Load(".env")
	if err != nil {
		childLogger.Info().Err(err).Send()
	}

	var certTls model.Cert

	if os.Getenv("SERVER_WITH_TLS") == "true" {	
		childLogger.Info().Msg("*** Loading cert.pem AND private_key.pem ***")

		certTls.IsTLS = true

		certTls.CertPEM, err = os.ReadFile("/var/pod/cert/tls.crt") // full_chain_b64.pem
		if err != nil {
			childLogger.Error().Err(err).Send()
			panic(err)
		} 

		certTls.CertPrivKeyPEM, err = os.ReadFile("/var/pod/cert/tls.key") //decrypted_private_key_b64.pem 
		if err != nil {
			childLogger.Error().Err(err).Send()
			panic(err)
		}

		// Just to show the cert in plain text 
		cert_str, err := base64.StdEncoding.DecodeString(string(certTls.CertPEM))
		if err != nil {
			childLogger.Error().Err(err).Send()
			panic(err)
		}
		certTls.CertPEMStr = string(cert_str)
		certTls.CertPEM = cert_str

		cert_str, err = base64.StdEncoding.DecodeString(string(certTls.CertPrivKeyPEM))
		if err != nil {
			childLogger.Error().Err(err).Send()
			panic(err)
		}
		certTls.CertPrivKeyPEMStr = string(cert_str)
		certTls.CertPrivKeyPEM = cert_str
	}

	return certTls
}