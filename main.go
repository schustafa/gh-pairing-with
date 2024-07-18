package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"

	"github.com/cli/go-gh/v2/pkg/auth"
	mapset "github.com/deckarep/golang-set/v2"
	fgql "github.com/mergestat/fluentgraphql"
)

type User struct {
	DatabaseID int
	Email      string
	Login      string
	Name       string
}

func (user User) coAuthoredBy() string {
	coauthoredName := user.Name
	if coauthoredName == "" {
		coauthoredName = user.Login
	}

	coauthoredEmail := user.Email
	if coauthoredEmail == "" {
		coauthoredEmail = fmt.Sprintf("%d+%s@users.noreply.github.com", user.DatabaseID, user.Login)
	}

	return fmt.Sprintf("Co-authored-by: %s <%s>\n", coauthoredName, coauthoredEmail)
}

func main() {
	flag.Parse()

	if len(flag.Args()) < 1 {
		fmt.Printf(`
Usage:
  pairing-with <github_login>...
`)
		return
	}

	if err := cli(); err != nil {
		fmt.Fprintf(os.Stderr, "gh-pairing-with failed: %s\n", err.Error())
		os.Exit(1)
	}
}

// generateQuery accepts an array of usernames and returns a
// map[string]interface{} suitable for marshalling to JSON for a GraphQL query.
func generateQuery(usernames []string) map[string]interface{} {
	userQuery := fgql.NewQuery()

	for i, login := range usernames {
		userQuery.Selection(
			"user",
			fgql.WithAlias(fmt.Sprintf("user_%d", i)),
			fgql.WithArguments(
				fgql.NewArgument("login", fgql.NewStringValue(login)),
			),
		).
			Scalar("name").
			Scalar("email").
			Scalar("login").
			Scalar("databaseId")
	}

	body := map[string]interface{}{
		"query": userQuery.Root().String(),
	}

	return body
}

// missingTokenScopes returns a set of scopes that are required but not present
// in the passed string value of the X-OAuth-Scopes header.
func missingTokenScopes(scopesHeader string) mapset.Set[string] {
	requiredScopes := mapset.NewSet[string]("user:email", "read:user")
	scopesHeader = strings.ReplaceAll(scopesHeader, " ", "")

	tokenScopes := mapset.NewSet[string](strings.Split(scopesHeader, ",")...)

	return requiredScopes.Difference(tokenScopes)
}

func cli() error {
	// Parse handles from command-line arguments and generate a request body
	handles := flag.Args()
	body := generateQuery(handles)

	// Marshal the request body to JSON; return and print error if that fails
	b, err := json.Marshal(body)
	if err != nil {
		return fmt.Errorf("could not marshal body: %w", err)
	}

	// Build the request; return and print error if that fails
	req, err := http.NewRequest(http.MethodPost, "https://api.github.com/graphql", bytes.NewReader(b))
	if err != nil {
		return fmt.Errorf("could not build request: %w", err)
	}

	// Update the authorization header with the GitHub token
	githubToken, _ := auth.TokenForHost("github.com")
	req.Header.Set("Authorization", fmt.Sprintf("bearer %s", githubToken))

	// Make the request; return and print error if that fails
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return fmt.Errorf("API call failed: %w", err)
	}

	// Read the response body; return and print error if that fails
	resBody, err := io.ReadAll(res.Body)
	if err != nil {
		return fmt.Errorf("could not read: %w", err)
	}

	var graphqlResponse map[string]interface{}

	// Unmarshal the response body; return and print error if that fails
	err = json.Unmarshal(resBody, &graphqlResponse)
	if err != nil {
		return fmt.Errorf("could not unmarshal: %w", err)
	}

	data, ok := graphqlResponse["data"].(map[string]interface{})

	// If the response body does not contain a "data" key, the token may be
	// missing required scopes
	if !ok {
		missingScopes := missingTokenScopes(res.Header.Get("X-OAuth-Scopes"))

		if missingScopes.Cardinality() > 0 {
			return fmt.Errorf("your token is missing required scopes. try running the following:\n\tgh auth refresh --scopes %s", strings.Join(missingScopes.ToSlice(), ","))
		}

		return fmt.Errorf("could not parse response")
	}

	// Print a co-authored-by line for each user in the returned data set
	for _, user := range data {
		if user == nil {
			continue
		}

		var userData User
		userJson, _ := json.Marshal(user)
		json.Unmarshal(userJson, &userData)

		fmt.Print(userData.coAuthoredBy())
	}

	return nil
}
