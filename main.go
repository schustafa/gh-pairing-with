package main

import (
	"bufio"
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

	"github.com/schustafa/gh-pairing-with/config"
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

// gh pairing-with schustafa
// gh pairing-with schustafa stephanieg0
// gh pairing-with --alias buddies schustafa stephanieg0
// gh pairing-with --alias kiran krhkt
// gh pairing-with buddies
// gh pairing-with --alias buddies
// gh pairing-with --list-aliases
// gh pairing-with --delete-alias buddies

// LATER:
// gh pairing-with --start buddies
// gh pairing-with --stop

func main() {
	if err := cli(); err != nil {
		fmt.Fprintf(os.Stderr, "gh-pairing-with failed: %s\n", err.Error())
		os.Exit(1)
	}
}

func cli() error {
	cfg, err := config.LoadConfig()
	if err != nil {
		return err
	}

	var aliasFlag string
	var deleteAliasFlag string
	flag.StringVar(&aliasFlag, "alias", "", "alias for a handle or set of handles")
	flag.StringVar(&deleteAliasFlag, "delete-alias", "", "delete a specified alias")
	listAliasesFlag := flag.Bool("list-aliases", false, "list all aliases")

	flag.Parse()

	rawHandles := flag.Args()

	if aliasFlag != "" {
		if err := cfg.AddAliasForHandles(aliasFlag, rawHandles); err != nil {
			return err
		}

		return nil
	}

	if deleteAliasFlag != "" {
		aliasExists := cfg.AliasExists(deleteAliasFlag)

		if !aliasExists {
			// Alias is already gone, just be quiet
			return nil
		}

		reader := bufio.NewReader(os.Stdin)

		fmt.Printf("delete alias %s? [y/n] ", deleteAliasFlag)

		response, err := reader.ReadString('\n')
		if err != nil {
			return fmt.Errorf("could not read: %w", err)
		}

		response = strings.ToLower(strings.TrimSpace(response))

		if response == "y" || response == "yes" {
			if err := cfg.DeleteAlias(deleteAliasFlag); err != nil {
				return err
			}
		}

		return nil
	}

	if *listAliasesFlag {
		aliases := cfg.GetAllAliases()
		for alias, handles := range aliases {
			fmt.Printf("%s: %v\n", alias, strings.Join(handles, " "))
		}

		return nil
	}

	if len(rawHandles) < 1 {
		fmt.Printf(`
	Usage:
	  pairing-with <github_login>...
`)
		return nil
	}

	expandedHandles := cfg.ExpandHandles(rawHandles)

	if err := lookupAndPrintForHandles(expandedHandles); err != nil {
		return err
	}

	return nil
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

func lookupAndPrintForHandles(handles []string) error {
	// generate a request body for the handles
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
