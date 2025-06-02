package events

import (
	"context"
	"encoding/json"

	"github.com/cybercoder/tlscdn-controller/pkg/logger"
	"github.com/cybercoder/tlscdn-controller/pkg/redis"
	corev1 "k8s.io/api/core/v1"
)

func OnAddSecret(obj interface{}) {
	secret := obj.(*corev1.Secret)
	annotations := secret.GetAnnotations()
	if annotations == nil || annotations["cdn.ik8s.ir/hostname"] == "" || annotations["cert-manager.io/certificate-name"] == "" {
		return
	}
	logger.Info("cdn.ik8s.ir/hostname:", annotations["cdn.ik8s.ir/hostname"])
	upsertCertificateToRedis(annotations["cdn.ik8s.ir/hostname"], secret.Data["tls.crt"], secret.Data["tls.key"])
	publishCert(annotations["cdn.ik8s.ir/hostname"], secret.Data["tls.crt"], secret.Data["tls.key"])
}

func OnUpdateSecret(prev interface{}, obj interface{}) {
	secret := obj.(*corev1.Secret)
	annotations := secret.GetAnnotations()
	if annotations == nil || annotations["cdn.ik8s.ir/hostname"] == "" || annotations["cert-manager.io/certificate-name"] == "" {
		return
	}
	upsertCertificateToRedis(annotations["cdn.ik8s.ir/hostname"], secret.Data["tls.crt"], secret.Data["tls.key"])
	publishCert(annotations["cdn.ik8s.ir/hostname"], secret.Data["tls.crt"], secret.Data["tls.key"])
}

func OnDeleteSecret(obj interface{}) {
	secret := obj.(*corev1.Secret)
	annotations := secret.GetAnnotations()
	if annotations == nil || annotations["cdn.ik8s.ir/hostname"] == "" || annotations["cert-manager.io/certificate-name"] == "" {
		return
	}
	deleteRedisKey(annotations["cdn.ik8s.ir/hostname"])
	invalidateCertCache(annotations["cdn.ik8s.ir/hostname"])
}

func publishCert(hostname string, crt []byte, key []byte) {
	redisClient := redis.CreateClient()

	certificateData := struct {
		Hostname string `json:"hostname"`
		Crt      string `json:"crt"`
		Key      string `json:"key"`
	}{
		Hostname: hostname,
		Crt:      string(crt),
		Key:      string(key),
	}
	stringify, err := json.Marshal(certificateData)
	if err != nil {
		logger.Errorf("Failed on json marshal certificate: %v", err)
		return
	}
	status := redisClient.Publish(context.Background(), "new_cert", stringify)
	logger.Info("status is: %v", status)
}

func upsertCertificateToRedis(hostname string, crt []byte, key []byte) {
	redisClient := redis.CreateClient()
	err := redisClient.HSet(
		context.Background(),
		"cdngateway:"+hostname+":tls",
		"hostname", hostname,
		"crt", string(crt),
		"key", string(key),
	).Err()
	if err != nil {
		logger.Errorf("error on redis cert set for host %s : %v", hostname, err)
	}
}

func deleteRedisKey(hostname string) {
	redisClient := redis.CreateClient()
	err := redisClient.Del(context.Background(), hostname+":tls").Err()
	if err != nil {
		logger.Errorf("error on cert, key deletion in redis for host %s : %v", hostname, err)
	}
}

func invalidateCertCache(hostname string) {
	redisClient := redis.CreateClient()
	err := redisClient.Publish(context.Background(), "invalidate_cert_cache", hostname).Err()
	if err != nil {
		logger.Errorf("[cert] cache invalidation for %s was unsuccessful: %v", hostname, err)
	}
}
