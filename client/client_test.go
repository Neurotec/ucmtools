package client

import (
	"testing"
)

func TestPing(t *testing.T) {
	client, err := New("https://localhost:8089")
	if err != nil {
		t.Fatal(err.Error())
	}
	if client.Login("admin", "oo") == true {
		t.Fatal("Contrasena invalida")
	}

	if client.Login("admin", "test") != true {
		t.Fatal("Contrasena valida")
	}

	t.Log(client.ListTrunks())
	t.Log(client.ListAccounts())
}
