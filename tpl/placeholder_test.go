package tpl

import (
	log "github.com/Sirupsen/logrus"
	"testing"
)

func init() {
	log.SetLevel(log.DebugLevel)
}

func TestPlaceholderTemplateRender(t *testing.T) {
	actual, err := PlaceholderTemplate{
		[]byte(`apiVersion: apps/v1beta1
kind: Deployment
metadata:
  name: $NAME-deployment
  annotations:
    replicas-as-string: "$REPLICAS"
    key: "${NAME}$$VALUE" # $${...} and $$$$ test
spec:
  replicas: $REPLICAS
`),
	}.Render(map[string]interface{}{
		"NAME":     "app",
		"NOT_USED": "value",
		"REPLICAS": 1,
	})
	if err != nil {
		t.Fatal(err)
	}
	expected := `apiVersion: apps/v1beta1
kind: Deployment
metadata:
  name: app-deployment
  annotations:
    replicas-as-string: "1"
    key: "app$VALUE" # ${...} and $$ test
spec:
  replicas: 1
`
	if string(actual) != expected {
		t.Fatalf("actual: \n%s != expected: \n%s", actual, expected)
	}
}

func TestPlaceholderTemplateRenderIncomplete(t *testing.T) {
	_, err := PlaceholderTemplate{
		[]byte(`apiVersion: apps/v1beta1
kind: Deployment
metadata:
  name: $NAME-deployment
`),
	}.Render(map[string]interface{}{
		"NOT_USED": "value",
	})
	if err == nil {
		t.Fatal()
	}
	expected := `4:9: "NAME" isn't set`
	if err.Error() != expected {
		t.Fatalf("actual: \n%s != expected: \n%s", err.Error(), expected)
	}
}
