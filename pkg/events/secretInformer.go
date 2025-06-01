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
	updateRedisKey(annotations["cdn.ik8s.ir/hostname"], secret.Data["tls.crt"], secret.Data["tls.key"])
}

func OnUpdateSecret(prev interface{}, obj interface{}) {
	secret := obj.(*corev1.Secret)
	annotations := secret.GetAnnotations()
	if annotations == nil || annotations["cdn.ik8s.ir/hostname"] == "" || annotations["cert-manager.io/certificate-name"] == "" {
		return
	}
	updateRedisKey(annotations["cdn.ik8s.ir/hostname"], secret.Data["tls.crt"], secret.Data["tls.key"])
}

func OnDeleteSecret(obj interface{}) {
	secret := obj.(*corev1.Secret)
	annotations := secret.GetAnnotations()
	if annotations == nil || annotations["cdn.ik8s.ir/hostname"] == "" || annotations["cert-manager.io/certificate-name"] == "" {
		return
	}
	invalidateCertCache(annotations["cdn.ik8s.ir/hostname"])
}

func updateRedisKey(hostname string, crt []byte, key []byte) {
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
	if err == nil {
		redisClient.Publish(context.Background(), "new_cert", stringify)
		
	}
}

// func deleteRedisKey(hostname string) {
// 	redisClient := redis.CreateClient()
// 	err := redisClient.Del(context.Background(), "cdngateway:"+hostname+":tls:crt", "cdngateway:"+hostname+":tls:key").Err()
// 	if err != nil {
// 		logger.Errorf("error on cert, key deletion in redis for host %s : %v", hostname, err)
// 	}
// }

func invalidateCertCache(hostname string) {
	redisClient := redis.CreateClient()
	err := redisClient.Publish(context.Background(), "invalidate_cert_cache", hostname).Err()
	if err != nil {
		logger.Errorf("[cert] cache invalidation for %s was unsuccessful: %v", hostname, err)
	}
}
