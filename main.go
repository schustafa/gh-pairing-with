package main

import (
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/cli/go-gh/v2/pkg/api"
	"github.com/cli/go-gh/v2/pkg/term"
	graphql "github.com/cli/shurcooL-graphql"
)

type userQuery struct {
	User struct {
		Name       graphql.String
		Email      graphql.String
		DatabaseID graphql.Int
	} `graphql:"user(login: $login)"`
}

func main() {
	flag.Parse()

	if len(flag.Args()) < 1 {
		fmt.Printf(`
Usage:
  pairing-with <github-login>
`)
		return
	}

	if err := cli(); err != nil {
		fmt.Fprintf(os.Stderr, "gh-pairing-with failed: %s\n", err.Error())
		os.Exit(1)
	}
}

func cli() error {
	login := strings.ToLower(strings.Join(flag.Args(), " "))

	terminal := term.FromEnv()
	if terminal.IsTerminalOutput() {
		fmt.Printf("looking for user %s\n", login)
	}

	client, err := api.DefaultGraphQLClient()
	if err != nil {
		return fmt.Errorf("could not make graphql client: %w", err)
	}

	variables := map[string]interface{}{
		"login": graphql.String(login),
	}

	var query userQuery

	err = client.Query("UserSearch", &query, variables)
	if err != nil {
		return fmt.Errorf("API call failed: %w", err)
	}

	coauthoredName := fmt.Sprintf("%s", query.User.Name)
	if coauthoredName == "" {
		coauthoredName = login
	}

	coauthoredEmail := fmt.Sprintf("%s", query.User.Email)
	if coauthoredEmail == "" {
		coauthoredEmail = fmt.Sprintf("%d+%s@users.noreply.github.com", query.User.DatabaseID, login)
	}

	fmt.Printf("Co-authored-by: %s <%s>\n", coauthoredName, coauthoredEmail)

	return nil
}
