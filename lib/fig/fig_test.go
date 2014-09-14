package fig

import (
	"testing"
)

const FigFile = `web:
    build: .
    command: python app.py
    links:
    - db
    ports:
    - "8000:8000"
db:
  image: postgres`

func TestReadFig(t *testing.T) {
	figFile, err := ReadFig([]byte(FigFile))
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("%+v", figFile)
	if len(figFile.Services) != 1 {
		t.Fatalf("expected 1 service, got %d", len(figFile.Services))
	}
	dbSvc, ok := figFile.Services["db"]
	if !ok {
		t.Fatal("didn't find the db service")
	}
	if dbSvc.Image != "postgres" {
		t.Fatalf("found %s image in the db service, expected postgres", dbSvc.Image)
	}
}
