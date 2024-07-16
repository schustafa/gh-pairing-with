package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGenerateQuery(t *testing.T) {
	expectedQuery := `{
  user_0: user(login: "monalisa") {
    name
    email
    login
    databaseId
  }
  user_1: user(login: "hubot") {
    name
    email
    login
    databaseId
  }
}`
	usernames := []string{"monalisa", "hubot"}
	queryBody := generateQuery(usernames)
	assert.NotNil(t, queryBody["query"])
	assert.Equal(t, expectedQuery, queryBody["query"])
}
