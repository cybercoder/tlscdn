package events

import (
	"context"
	"encoding/base64"

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
	updateRedisKey(annotations["cdn.ik8s.ir/hostname"], secret.Data["tls.crt"], secret.Data["tls.key"])
	invalidateCertCache(annotations["cdn.ik8s.ir/hostname"])
}

func OnUpdateSecret(prev interface{}, obj interface{}) {
	secret := obj.(*corev1.Secret)
	annotations := secret.GetAnnotations()
	if annotations == nil || annotations["cdn.ik8s.ir/hostname"] == "" || annotations["cert-manager.io/certificate-name"] == "" {
		return
	}
	updateRedisKey(annotations["cdn.ik8s.ir/hostname"], secret.Data["tls.crt"], secret.Data["tls.key"])
	invalidateCertCache(annotations["cdn.ik8s.ir/hostname"])
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

func updateRedisKey(hostname string, crt []byte, key []byte) {
	redisClient := redis.CreateClient()
	pipeline := redisClient.Pipeline()
	tlsCrt, err := base64.StdEncoding.DecodeString(string(crt))
	if err != nil {
		logger.Errorf("error decoding tls.crt for host %s: %v", hostname, err)
		return
	}
	tlsKey, err := base64.StdEncoding.DecodeString(string(key))
	if err != nil {
		logger.Errorf("error decoding tls.key for host %s: %v", hostname, err)
		return
	}
	pipeline.Set(context.Background(), "cdngateway:"+hostname+":tls:crt", tlsCrt, 0)
	pipeline.Set(context.Background(), "cdngateway:"+hostname+":tls:key", tlsKey, 0)

	_, err = pipeline.Exec(context.Background())
	if err != nil {
		logger.Errorf("error on redis cert set for host %s : %v", hostname, err)
	}
}

func deleteRedisKey(hostname string) {
	redisClient := redis.CreateClient()
	err := redisClient.Del(context.Background(), "cdngateway:"+hostname+":tls:crt", "cdngateway:"+hostname+":tls:key").Err()
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
