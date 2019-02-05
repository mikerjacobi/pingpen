package main

import (
	"bytes"
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	_ "github.com/go-sql-driver/mysql"
	"github.com/google/uuid"
	"github.com/husobee/vestigo"
	"github.com/sirupsen/logrus"
)

type Controller struct {
	//config items
	Grid   string
	Table  string
	Region string
	dbhost string
	dbuser string
	dbpw   string
	dbname string
}

type Note struct {
	id        string
	accountID string
	content   string
}

func dosql(db *sql.DB) string {
	n := Note{}
	n.id = uuid.New().String()
	n.accountID = uuid.New().String()
	n.content = fmt.Sprintf("content is %s", uuid.New().String())
	stmt := `insert into notes(id,account_id,note) values (?,?,?)`
	if _, err := db.Exec(stmt, n.id, n.accountID, n.content); err != nil {
		return "DBERROR: " + err.Error()
	}
	return "success"
}

// Handler is our lambda handler invoked by the `lambda.Start` function call
func (c *Controller) Handler(ctx context.Context, req *events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	var buf bytes.Buffer
	resp := events.APIGatewayProxyResponse{}

	connStr := fmt.Sprintf("%s:%s@tcp(%s:3306)/%s?parseTime=true", c.dbuser, c.dbpw, c.dbhost, c.dbname)

	db, err := sql.Open("mysql", connStr)
	if err != nil {
		logrus.WithError(err).Errorf("failed to open mysql: %s", strings.Replace(connStr, c.dbpw, "XXX", -1))
		resp.StatusCode = 500
		return resp, err
	}
	defer db.Close()

	if err := db.Ping(); err != nil {
		logrus.WithError(err).Errorf("failed to open mysql: %s", strings.Replace(connStr, c.dbpw, "XXX", -1))
		resp.StatusCode = 500
		return resp, err
	}

	out := dosql(db)
	body, err := json.Marshal(map[string]interface{}{
		"message": fmt.Sprintf("config: %+v", c),
		"out":     out,
	})
	if err != nil {
		resp.StatusCode = 404
		return resp, err
	}
	logrus.Infof("here8")
	json.HTMLEscape(&buf, body)

	resp.StatusCode = 200
	resp.IsBase64Encoded = false
	resp.Body = buf.String()
	resp.Headers = map[string]string{
		"Content-Type":           "application/json",
		"X-MyCompany-Func-Reply": "hello-handler",
	}

	return resp, nil
}

func (c *Controller) httpHandler(rw http.ResponseWriter, r *http.Request) {
	req := &events.APIGatewayProxyRequest{}
	if err := json.NewDecoder(r.Body).Decode(req); err != nil {
		rw.WriteHeader(http.StatusBadRequest)
		logrus.WithError(err).Errorf("failed to unmarshal request into api gateway request")
		return
	}

	resp, err := c.Handler(context.Background(), req)
	rw.WriteHeader(resp.StatusCode)
	if err != nil {
		logrus.WithError(err).Errorf("failed to execute handler")
		return
	}

	respBody, err := json.Marshal(resp.Body)
	if err != nil {
		logrus.WithError(err).Errorf("failed to marshal response")
		return
	}
	rw.Write(respBody)
}

func main() {
	c := &Controller{
		Grid:   os.Getenv("GRID"),
		Table:  os.Getenv("TABLE"),
		Region: os.Getenv("REGION"),
		dbhost: os.Getenv("DBHOST"),
		dbuser: os.Getenv("DBUSER"),
		dbpw:   os.Getenv("DBPW"),
		dbname: os.Getenv("DBNAME"),
	}

	logrus.Info(c)
	if c.Grid == "sandbox" {
		port := os.Getenv("PORT")
		logrus.Info(port)
		router := vestigo.NewRouter()
		router.Post("/note", c.httpHandler)
		logrus.Fatal(http.ListenAndServe("0.0.0.0:"+port, router))
	} else {
		lambda.Start(c.Handler)
	}
}
