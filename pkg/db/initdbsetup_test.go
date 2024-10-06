package db_test

import (
	"database/sql"
	"fmt"
	"strings"
	"testing"
)

func TestStartYourEngine(t *testing.T) {
	tests := []struct {
		urlString string
		dbType    string
	}{
		{"localhost,port=5432,user=myuser,password=mypassword,database=mydb", ""},
		{"host=localhost port=5432 user=myuser password=mypassword dbname=mydb sslmode=require", "postgresql"},
	}

	for _, test := range tests {
		engine, err := StartYourEngine(test.urlString, test.dbType)
		if err != nil {
			t.Errorf("StartYourEngine() returned an error: %v", err)
		} else if engine == nil {
			t.Error("StartYourEngine() did not return a valid database connection")
		}
	}

}

func StartYourEngine(urlString string, dbType string) (*sql.DB, error) {
	urlList := strings.Split(urlString, ",")
	if len(urlList) == 6 || (dbType != "" && len(urlList) == 7) {
		pgURL := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s",
			urlList[0], urlList[4], urlList[1], urlList[2], urlList[3])
		if dbType != "" && len(urlList) == 7 {
			pgURL += fmt.Sprintf(" sslmode=require")
		}
		db, err := sql.Open("postgres", pgURL)
		return &db, err
	} else if len(urlList) == 1 {
		fmt.Println("Complete url received.")
		newUrl := urlList[0]
		dialect := ""
		if dbType != "" && len(urlList) > 6 {
			urlDialect := fmt.Sprintf("%s+%s", newUrl, urlList[5])
			pgURL = urlDialect
			dialect = urlList[5]
		} else if len(urlList) == 1 {
			fmt.Println("Url shall be created")
			pgURL = "host=localhost port=5432 user=myuser password=mypassword dbname=mydb sslmode=require"
			dialect := ""
		}
		db, err := sql.Open("postgres", pgURL)
		return &db, err
	} else {
		return nil, fmt.Errorf("%d is not the right amount of values for the url. [1, 6 or 7]\n", len(urlList))
	}
}
