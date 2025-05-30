package k8s

import (
	"context"
	"os"

	cmapi "github.com/cert-manager/cert-manager/pkg/apis/certmanager/v1"
	cmmeta "github.com/cert-manager/cert-manager/pkg/apis/meta/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func CreateLetsEncryptWildCardCertificate(namespace string, gateway string, hostname string) (*cmapi.Certificate, error) {
	clusterIssuer := os.Getenv("CERT_ISSUER")
	if clusterIssuer == "" {
		clusterIssuer = "ik8s-letsencrypt-webhook"
	}
	cert := &cmapi.Certificate{
		ObjectMeta: metav1.ObjectMeta{
			Name:      gateway,
			Namespace: namespace,
		},
		Spec: cmapi.CertificateSpec{
			SecretName: gateway + "-cert-tls",
			IssuerRef: cmmeta.ObjectReference{
				Name: clusterIssuer,
				Kind: "ClusterIssuer",
			},
			DNSNames: []string{hostname, "*." + hostname},
			SecretTemplate: &cmapi.CertificateSecretTemplate{
				Annotations: map[string]string{
					"cdn.ik8s.ir/hostname": hostname,
				},
			},
		},
	}

	cmClient = CreateCertManagerClient()
	crt, err := cmClient.CertmanagerV1().Certificates(namespace).Create(context.Background(), cert, metav1.CreateOptions{})
	if err != nil {
		return nil, err
	}
	return crt, nil
}
