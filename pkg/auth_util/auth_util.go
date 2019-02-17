package auth_util

import (
	_ "github.com/lib/pq"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"os"
)

const callersFilename = "/etc/delinkcious/mutual-auth.yaml"

var callersByName = map[string]string{}
var callersByToken = map[string][]string{}

func init() {
	if os.Getenv("DELINKCIOUS_MUTUAL_AUTH") == "false" {
		return
	}

	data, err := ioutil.ReadFile(callersFilename)
	if err != nil {
		panic(err)
	}
	err = yaml.Unmarshal(data, callersByName)
	if err != nil {
		panic(err)
	}

	for caller, token := range callersByName {
		callersByToken[token] = append(callersByToken[token], caller)
	}
}

func GetToken(caller string) string {
	return callersByName[caller]
}

func HasCaller(caller string, token string) bool {
	for _, c := range callersByToken[token] {
		if c == caller {
			return true
		}
	}

	return false
}
