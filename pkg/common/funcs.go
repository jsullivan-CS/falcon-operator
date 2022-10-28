package common

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"regexp"
	"strings"

	"k8s.io/apimachinery/pkg/version"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
)

func InitContainerArgs() []string {
	return []string{
		"-c",
		"mkdir -p " + FalconDataDir +
			" && " +
			"touch " + FalconStoreFile,
	}
}

func InitCleanupArgs() []string {
	return []string{
		"-c",
		"rm -rf " + FalconDataDir,
	}
}

func CleanupSleep() []string {
	return []string{
		"-c",
		"sleep 10",
	}
}

func FCAdmissionReviewVersions() []string {
	kubeVersion := GetKubernetesVersion()
	fcArv := []string{"v1"}

	if strings.Compare(kubeVersion.Minor, "22") < 0 {
		fcArv = []string{"v1", "v1beta"}
	}

	return fcArv
}

func GetKubernetesVersion() *version.Info {
	// creates the in-cluster config
	config, err := rest.InClusterConfig()
	if err != nil {
		panic(err.Error())
	}
	// creates the clientset
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		panic(err.Error())
	}

	version, err := clientset.DiscoveryClient.ServerVersion()
	if err != nil {
		panic(err.Error())
	}

	return version
}

func EncodedBase64String(data string) []byte {
	base64EncodedData := make([]byte, base64.StdEncoding.EncodedLen(len(data)))
	base64.StdEncoding.Encode(base64EncodedData, []byte(data))
	return base64EncodedData
}

func EncodeBase64Interface(i interface{}) (string, error) {
	buf := bytes.Buffer{}
	b64enc := base64.NewEncoder(base64.StdEncoding, &buf)
	if err := json.NewEncoder(b64enc).Encode(i); err != nil {
		return "", fmt.Errorf("failed to convert to base64 encoding: %v", err)
	}
	if err := b64enc.Close(); err != nil {
		return "", fmt.Errorf("failed to close base64 encoder: %v", err)
	}
	return buf.String(), nil
}

func CleanDecodedBase64(s []byte) []byte {
	re := regexp.MustCompile(`[\t|\n]*`)
	b64byte, err := base64.StdEncoding.DecodeString(string(s))
	if err != nil {
		return []byte(re.ReplaceAllString(string(s), ""))
	}
	return []byte(re.ReplaceAllString(string(b64byte), ""))
}

func MapCopy(src map[string]string, dst map[string]string) map[string]string {
	for k, v := range src {
		dst[k] = v
	}
	return dst
}
