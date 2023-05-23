package db

import (
	"os"
	"fmt"

	"urlshortener/internal/utils"

	_ "github.com/joho/godotenv/autoload"
	log "github.com/sirupsen/logrus"
	"github.com/tarantool/go-tarantool"
)

var conn *tarantool.Connection

// getConnection takes host:port from env variable
// Then, trying to connect to tarantool and returns connection pointer.
func GetConnection() *tarantool.Connection {
	var err error
	connStr := fmt.Sprintf("%s:%s", os.Getenv("TARANTOOL_HOST"), os.Getenv("TARANTOOL_PORT"))
	conn, err = tarantool.Connect(connStr, tarantool.Opts{})
	if err != nil {
		log.Error(fmt.Errorf("error in connecting to db: %v", err))
	}
	return conn
}

// CloseConnection is used as defer in main.
func CloseConnection() {
	conn.Close()
}

// InitDB builds tarantool schema and indexes if them not exists yet.
func InitDB() {
	conn := GetConnection()
	defer conn.Close()

	_, _ = conn.Call("box.schema.space.create", []interface{}{
		"urls",
		map[string]bool{"if_not_exists": true}})

	_, err := conn.Call(fmt.Sprintf("box.space.%s:format", os.Getenv("TARANTOOL_COLL")), [][]map[string]string{
		{
			{"name": "id", "type": "string"},
			{"name": "original_url", "type": "string"},
		}})
	if err != nil {
		log.Error(fmt.Errorf("error in creating tarantool schema: %v", err))
	}

	_, err = conn.Call(fmt.Sprintf("box.space.%s:create_index", os.Getenv("TARANTOOL_COLL")), []interface{}{
		"primary",
		map[string]interface{}{
			"parts":         []string{"id"},
			"if_not_exists": true}})
	if err != nil {
		log.Error(fmt.Errorf("error in creating tarantool indexes: %v", err))
	}

}

// InsertUrl takes pair of hashed and original url and saves it in tarantool.
func InsertUrl(hashedUrl string, originalUrl string) (*tarantool.Response, error) {
	res, err := conn.Insert("urls", []interface{}{hashedUrl, originalUrl})
	if err != nil {
		return nil, fmt.Errorf("error in inserting record :%v", err)
	}
	return res, nil
}

// GetUrl selects record in tarantool for requested hashed url.
// If exists - maps it to UrlObject and returns original url.
func GetUrl(hashedUrl string) (string, error) {
	var url []utils.UrlObject
	err := conn.SelectTyped("urls", "primary", 0, 1, tarantool.IterEq, []interface{}{hashedUrl}, &url)
	if err != nil {
		return "", fmt.Errorf("error in getting record: %v", err)
	}
	if len(url) == 0 {
		return "", nil
	}
	return url[0].OriginalUrl, nil
}
