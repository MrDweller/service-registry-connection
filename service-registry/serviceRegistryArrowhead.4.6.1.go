package serviceregistry

import (
	"bytes"
	"crypto/tls"
	"crypto/x509"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"time"

	"github.com/MrDweller/service-registry-connection/models"
)

const SERVICE_REGISTRY_ARROWHEAD_4_6_1 ServiceRegistryImplementationType = "serviceregistry-arrowhead-4.6.1"

type ServiceRegistryArrowhead_4_6_1 struct {
	ServiceRegistry
	models.CertificateInfo
}

type RegisterServiceDTO struct {
	models.ServiceDefinition
	Interfaces     []string                `json:"interfaces"`
	ProviderSystem models.SystemDefinition `json:"providerSystem"`
}

func (serviceRegistry ServiceRegistryArrowhead_4_6_1) Connect() error {

	result, err := serviceRegistry.echoServiceRegistry()
	if err != nil {
		return err
	}

	if string(result) != "Got it!" {
		return errors.New("can't establish a connection with the service registry")
	}

	return nil

}

func (serviceRegistry ServiceRegistryArrowhead_4_6_1) RegisterService(serviceDefinition models.ServiceDefinition, systemDefinition models.SystemDefinition) ([]byte, error) {
	reqisterServiceDTO := RegisterServiceDTO{
		ServiceDefinition: serviceDefinition,
		Interfaces: []string{
			"HTTP-SECURE-JSON",
		},
		ProviderSystem: systemDefinition,
	}
	payload, err := json.Marshal(reqisterServiceDTO)

	if err != nil {
		return nil, err
	}
	req, err := http.NewRequest("POST", "https://"+serviceRegistry.Address+":"+strconv.Itoa(serviceRegistry.Port)+"/serviceregistry/register", bytes.NewBuffer(payload))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")

	client, err := serviceRegistry.getClient()
	if err != nil {
		return nil, err
	}

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusCreated {
		errorString := fmt.Sprintf("status: %s, body: %s", resp.Status, string(body))
		return nil, errors.New(errorString)
	}

	return body, nil
}

func (serviceRegistry ServiceRegistryArrowhead_4_6_1) UnRegisterService(serviceDefinition models.ServiceDefinition, systemDefinition models.SystemDefinition) error {
	url := fmt.Sprintf("https://"+serviceRegistry.Address+":"+strconv.Itoa(serviceRegistry.Port)+"/serviceregistry/unregister?address=%s&port=%s&service_definition=%s&service_uri=%s&system_name=%s", url.QueryEscape(systemDefinition.Address), url.QueryEscape(strconv.Itoa(systemDefinition.Port)), url.QueryEscape(serviceDefinition.ServiceDefinition), url.QueryEscape(serviceDefinition.ServiceUri), url.QueryEscape(systemDefinition.SystemName))
	req, err := http.NewRequest("DELETE", url, nil)
	if err != nil {
		return err
	}

	client, err := serviceRegistry.getClient()
	if err != nil {
		return err
	}

	resp, err := client.Do(req)
	if err != nil {
		return err
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	if resp.StatusCode != http.StatusOK {
		errorString := fmt.Sprintf("status: %s, body: %s", resp.Status, string(body))
		return errors.New(errorString)
	}

	return nil
}

func (serviceRegistry ServiceRegistryArrowhead_4_6_1) RegisterSystem(systemDefinition models.SystemDefinition) ([]byte, error) {
	payload, err := json.Marshal(systemDefinition)

	if err != nil {
		return nil, err
	}
	req, err := http.NewRequest("POST", "https://"+serviceRegistry.Address+":"+strconv.Itoa(serviceRegistry.Port)+"/serviceregistry/register-system", bytes.NewBuffer(payload))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")

	client, err := serviceRegistry.getClient()
	if err != nil {
		return nil, err
	}

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusCreated {
		errorString := fmt.Sprintf("status: %s, body: %s", resp.Status, string(body))
		return nil, errors.New(errorString)
	}

	return body, nil
}

func (serviceRegistry ServiceRegistryArrowhead_4_6_1) UnRegisterSystem(systemDefinition models.SystemDefinition) error {
	url := fmt.Sprintf("https://"+serviceRegistry.Address+":"+strconv.Itoa(serviceRegistry.Port)+"/serviceregistry/unregister-system?address=%s&port=%s&system_name=%s", url.QueryEscape(systemDefinition.Address), url.QueryEscape(strconv.Itoa(systemDefinition.Port)), url.QueryEscape(systemDefinition.SystemName))
	req, err := http.NewRequest("DELETE", url, nil)
	if err != nil {
		return err
	}

	client, err := serviceRegistry.getClient()
	if err != nil {
		return err
	}

	resp, err := client.Do(req)
	if err != nil {
		return err
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	if resp.StatusCode != http.StatusOK {
		errorString := fmt.Sprintf("status: %s, body: %s", resp.Status, string(body))
		return errors.New(errorString)
	}

	return nil
}

func (serviceRegistry ServiceRegistryArrowhead_4_6_1) echoServiceRegistry() ([]byte, error) {
	req, err := http.NewRequest("GET", "https://"+serviceRegistry.Address+":"+strconv.Itoa(serviceRegistry.Port)+"/serviceregistry/echo", nil)
	if err != nil {
		return nil, err
	}

	client, err := serviceRegistry.getClient()
	if err != nil {
		return nil, err
	}

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return body, nil
}

func (serviceRegistry ServiceRegistryArrowhead_4_6_1) getClient() (*http.Client, error) {
	cert, err := tls.LoadX509KeyPair(serviceRegistry.CertFilePath, serviceRegistry.KeyFilePath)
	if err != nil {
		return nil, err
	}

	// Load truststore.p12
	truststoreData, err := os.ReadFile(serviceRegistry.Truststore)
	if err != nil {
		return nil, err

	}

	// Extract the root certificate(s) from the truststore
	pool := x509.NewCertPool()
	if ok := pool.AppendCertsFromPEM(truststoreData); !ok {
		return nil, err
	}

	client := &http.Client{
		Timeout: 5 * time.Second,
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				Certificates:       []tls.Certificate{cert},
				RootCAs:            pool,
				InsecureSkipVerify: false,
			},
		},
	}
	return client, nil
}
