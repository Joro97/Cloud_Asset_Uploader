package data

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"io/ioutil"
	"os"
	"time"

	"CloudAssetUploader/constants"

	"github.com/gocql/gocql"
	"github.com/rs/zerolog/log"
)

type AstraDB struct {
	Sess *gocql.Session
}

// TODO: Add proper validation and documentation.
func ConnectToAstra() (*gocql.Session, error) {
	astraHost := os.Getenv("ASTRA_HOST")
	astraPort := os.Getenv("ASTRA_PORT")
	astraUsername := os.Getenv("ASTRA_USERNAME")
	astraPassword := os.Getenv("ASTRA_PASSWORD")
	astraKeySpace := os.Getenv("ASTRA_KEYSPACE")

	workDir, err := os.Getwd()
	if err != nil {
		return nil, err
	}

	certPath := fmt.Sprintf("%s/%s/cert", workDir, constants.AstraConfigDirectory)
	keyPath := fmt.Sprintf("%s/%s/key", workDir, constants.AstraConfigDirectory)
	caPath := fmt.Sprintf("%s/%s/ca.crt", workDir, constants.AstraConfigDirectory)
	cert, err := tls.LoadX509KeyPair(certPath, keyPath)
	if err != nil {
		log.Error().Msgf("Error loading Key Pair for Astra: %s", err)
		return nil, err
	}

	caCert, err := ioutil.ReadFile(caPath)
	if err != nil {
		log.Error().Msgf("Error reading CA file: %s", err)
		return nil, err
	}

	caCertPool := x509.NewCertPool()
	caCertPool.AppendCertsFromPEM(caCert)
	tlsConfig := &tls.Config{
		Certificates: []tls.Certificate{cert},
		RootCAs:      caCertPool,
	}

	cluster := gocql.NewCluster(astraHost)
	cluster.SslOpts = &gocql.SslOptions{
		Config:                 tlsConfig,
		EnableHostVerification: false,
	}
	cluster.Authenticator = gocql.PasswordAuthenticator{
		Username: astraUsername,
		Password: astraPassword,
	}
	cluster.Hosts = []string{astraHost + ":" + astraPort}
	cluster.ConnectTimeout = 7 * time.Second
	cluster.Timeout = 7 * time.Second
	cluster.DisableInitialHostLookup = true
	cluster.Keyspace = astraKeySpace
	cluster.Consistency = gocql.Quorum

	session, err := cluster.CreateSession()
	if err != nil {
		log.Error().Msgf("Could not create session for Astra: %s", err)
		return nil, err
	}
	return session, nil
}

func (adb *AstraDB) InitializeTables() error {
	q := fmt.Sprintf(`CREATE TABLE IF NOT EXISTS %s (
		aws_name text,
		url text,
		upload_status text,
		PRIMARY KEY (aws_name)
	)`, constants.AssetUploaderDatabaseName)

	astraQuery := adb.Sess.Query(q)
	return astraQuery.Exec()
}

func (adb *AstraDB) AddNewAsset(assetName, url string) (id string, err error) {
	qStr := fmt.Sprintf(`INSERT INTO %s 
						(aws_name, url, upload_status) VALUES ('%s', '%s', '%s')`,
		constants.AssetUploaderDatabaseName, assetName, url, constants.AssetStatusCreated)

	astraQuery := adb.Sess.Query(qStr)
	err = astraQuery.Exec()
	if err != nil {
		log.Error().Msgf("Could not insert new asset into Astra: %s", err)
		return "", err
	}

	return assetName, err
}

func (adb *AstraDB) SetAssetStatus(assetID, status string) (*AssetInfo, error) {
	qStr := fmt.Sprintf(`UPDATE %s SET upload_status = '%s' WHERE aws_name = '%s'`,
		constants.AssetUploaderDatabaseName, status, assetID)

	astraQuery := adb.Sess.Query(qStr)
	err := astraQuery.Exec()
	if err != nil {
		log.Error().Msgf("Could not set asset status for asset with id: %s. Err: %s", assetID, err)
		return nil, err
	}

	return &AssetInfo{
		Name:         assetID,
		UploadStatus: status,
	}, nil
}

func (adb *AstraDB) GetAsset(assetID string) (*AssetInfo, error) {
	qStr := fmt.Sprintf(`SELECT aws_name, url, upload_status FROM %s WHERE aws_name = '%s'`,
		constants.AssetUploaderDatabaseName, assetID)

	asset := AssetInfo{}
	err := adb.Sess.Query(qStr).Consistency(gocql.One).Scan(&asset.Name, &asset.URL, &asset.UploadStatus)
	if err != nil {
		log.Error().Msgf("Could not set asset of status with id %s. Err: %s", assetID, err)
		return nil, err
	}

	return &asset, nil
}
