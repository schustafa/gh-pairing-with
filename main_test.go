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

func TestCoAuthoredByWithNoNameOrEmail(t *testing.T) {
	user := User{
		DatabaseID: 100101,
		Login:      "monalisa",
		Email:      "",
		Name:       "",
	}

	expectedCoAuthoredBy := "Co-authored-by: monalisa <100101+monalisa@users.noreply.github.com>\n"
	assert.Equal(t, expectedCoAuthoredBy, user.coAuthoredBy())
}

func TestCoAuthoredByWithNameAndNoEmail(t *testing.T) {
	user := User{
		DatabaseID: 100101,
		Login:      "monalisa",
		Email:      "",
		Name:       "Miss Mona Lisa Octocat",
	}

	expectedCoAuthoredBy := "Co-authored-by: Miss Mona Lisa Octocat <100101+monalisa@users.noreply.github.com>\n"
	assert.Equal(t, expectedCoAuthoredBy, user.coAuthoredBy())
}

func TestCoAuthoredByWithEmailAndNoName(t *testing.T) {
	user := User{
		DatabaseID: 100101,
		Login:      "monalisa",
		Email:      "monalisa@github.com",
		Name:       "",
	}

	expectedCoAuthoredBy := "Co-authored-by: monalisa <monalisa@github.com>\n"
	assert.Equal(t, expectedCoAuthoredBy, user.coAuthoredBy())
}

func TestCoAuthoredByWithEmailAndName(t *testing.T) {
	user := User{
		DatabaseID: 100101,
		Login:      "monalisa",
		Email:      "monalisa@github.com",
		Name:       "Miss Mona Lisa Octocat",
	}

	expectedCoAuthoredBy := "Co-authored-by: Miss Mona Lisa Octocat <monalisa@github.com>\n"
	assert.Equal(t, expectedCoAuthoredBy, user.coAuthoredBy())
}
