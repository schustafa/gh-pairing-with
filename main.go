package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/cli/go-gh/v2/pkg/auth"
	mapset "github.com/deckarep/golang-set/v2"
	fgql "github.com/mergestat/fluentgraphql"
	"gopkg.in/yaml.v3"
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

type Config struct {
	Aliases map[string][]string
}

func createConfigFileIfMissing(configFilePath string) error {
	if _, err := os.Stat(configFilePath); os.IsNotExist(err) {
		newConfigFile, err := os.OpenFile(
			configFilePath,
			os.O_RDWR|os.O_CREATE|os.O_EXCL,
			0666,
		)
		if err != nil {
			return err
		}

		var emptyConfig Config
		emptyConfig.Aliases = make(map[string][]string)

		blankConfigFile, err := yaml.Marshal(emptyConfig)
		if err != nil {
			return err
		}

		_, err = io.Writer.Write(newConfigFile, blankConfigFile)
		if err != nil {
			return err
		}

		defer newConfigFile.Close()
		return nil
	}

	return nil
}

func getConfigFilePath() (string, error) {
	const PairingWithDir = "gh-pairing-with"
	const ConfigYmlFileName = "config.yml"
	const DEFAULT_XDG_CONFIG_DIRNAME = ".config"

	configDir := os.Getenv("XDG_CONFIG_HOME")

	if configDir == "" {
		homeDir, err := os.UserHomeDir()
		if err != nil {
			return "", err
		}
		configDir = filepath.Join(homeDir, DEFAULT_XDG_CONFIG_DIRNAME)
	}

	pairingWithConfigDir := filepath.Join(configDir, PairingWithDir)
	return filepath.Join(pairingWithConfigDir, ConfigYmlFileName), nil
}

func loadConfig() (*Config, error) {
	var config Config

	configFilePath, err := getConfigFilePath()

	if err != nil {
		return nil, err
	}

	configDir := filepath.Dir(configFilePath)
	if _, err := os.Stat(configDir); os.IsNotExist(err) {
		if err = os.MkdirAll(configDir, os.ModePerm); err != nil {
			return &config, err
		}
	}

	if err := createConfigFileIfMissing(configFilePath); err != nil {
		return &config, err
	}

	existingFile, err := os.ReadFile(configFilePath)
	if err != nil {
		return &config, err
	}

	err = yaml.Unmarshal(existingFile, &config)
	if err != nil {
		return &config, err
	}

	return &config, nil
}

func (c *Config) save() error {
	configFilePath, err := getConfigFilePath()
	if err != nil {
		return err
	}

	updatedFile, err := yaml.Marshal(c)
	if err != nil {
		return err
	}

	f, err := os.Create(configFilePath)
	if err != nil {
		return err
	}
	defer f.Close()

	_, err = io.Writer.Write(f, updatedFile)
	if err != nil {
		return fmt.Errorf("could not write to config file: %w", err)
	}

	return nil
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
	cfg, err := loadConfig()
	if err != nil {
		return err
	}

	var aliasFlag string
	flag.StringVar(&aliasFlag, "alias", "", "alias for a handle or set of handles")

	flag.Parse()

	if len(flag.Args()) < 1 {
		fmt.Printf(`
	Usage:
	  pairing-with <github_login>...
`)
		return nil
	}

	handles := flag.Args()

	if aliasFlag != "" {
		if err := storeAliasForHandles(cfg, aliasFlag, handles); err != nil {
			return err
		}
	} else if err := lookupAndPrintForHandles(handles); err != nil {
		return err
	}

	return nil
}

func getAlias(alias string) ([]string, error) {
	fmt.Printf("getting alias %s\n", alias)

	var config Config

	existingFile, err := os.ReadFile("config.yml")
	if err != nil {
		return nil, fmt.Errorf("could not find file: %w", err)
	}

	err = yaml.Unmarshal(existingFile, &config)
	if err != nil {
		return nil, fmt.Errorf("could not unmarshal yaml: %w", err)
	}

	return config.Aliases[alias], nil
}

func storeAliasForHandles(config *Config, alias string, handles []string) error {
	fmt.Printf("storing alias %s for handles %v\n", alias, handles)

	config.Aliases[alias] = handles

	config.save()

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
	var aliasedHandles []string
	var handlesForQuery []string

	var err error
	// TODO: handle the case where multiple aliases/handles are passed
	aliasedHandles, err = getAlias(handles[0])
	if err != nil {
		return fmt.Errorf("error getting alias: %w", err)
	} else {
		if len(aliasedHandles) > 0 {
			handlesForQuery = aliasedHandles
		} else {
			handlesForQuery = handles
		}
	}

	// generate a request body for the handles
	body := generateQuery(handlesForQuery)

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
