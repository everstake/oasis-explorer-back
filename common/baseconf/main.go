package baseconf

import (
	"oasisTracker/common/helpers"
	"bytes"
	"encoding/json"
	"fmt"
	uuid "github.com/satori/go.uuid"
	"io"
	"os"
	"reflect"

	"github.com/fatih/structs"
	"github.com/hashicorp/vault/api"
	"github.com/mitchellh/mapstructure"
	"github.com/wedancedalot/decimal"
)

// BaseConfig is interface for all validatable structures
type BaseConfig interface {
	Validate() error
}

const (
	DEPLOYMENT_NAME = "DEPLOYMENT_NAME"
	KUBE_NAMESPACE   = "KUBE_NAMESPACE"
)

func decoderHook(in, out reflect.Type, input interface{}) (interface{}, error) {
	if out == reflect.TypeOf(decimal.Decimal{}) && in.Kind() == reflect.String {
		var num string
		switch n := input.(type) {
		case json.Number:
			num = string(n)
		case string:
			num = n
		default:
			num = fmt.Sprintf("%v", n)
		}
		return decimal.NewFromString(num)
	}

	if out == reflect.TypeOf(uuid.UUID{}) && in.Kind() == reflect.String {
		var uid string
		switch u := input.(type) {
		case string:
			uid = u
		default:
			uid = fmt.Sprintf("%v", u)
		}
		return uuid.FromString(uid)
	}

	return input, nil
}

// New creates a new config from vault. It heavily depends on the env variables. Namely:
// "DEPLOYMENT_NAME":"evertstake-service" - deployment name;
// "NS":"dev" - environment;
// "VAULT_ADDR":"https://vt.cryexch.xyz" - address of the vault;
// "VAULT_TOKEN":"<vault token>" - the token.
// The config is loaded using the "$DEPLOYMENT_NAME/$NS/data/config_json" key.
func New(cfg interface{}, cleanup bool) error {
	// load from vault
	vault, err := api.NewClient(nil)
	if err != nil {
		return err
	}
	if cleanup {
		// cleanup the env
		err = cleanEnv()
		if err != nil {
			return err
		}
	}
	err = loadToken(vault)
	if err != nil {
		return err
	}
	err = InitFromVault(cfg, vault)
	if err != nil {
		return err
	}

	return nil
}

func InitFromVault(cfg interface{}, vault *api.Client) error {
	if vault == nil {
		return fmt.Errorf("vault is nil")
	}
	basePath := getVaultPath()

	sec, err := vault.Logical().Read(fmt.Sprintf("%s/config_json", basePath))
	if err != nil {
		return err
	}
	if sec == nil {
		return fmt.Errorf("secret is nil")
	}
	d, ok := sec.Data["data"]
	if !ok {
		return fmt.Errorf("no data in vault")
	}
	config := &mapstructure.DecoderConfig{
		DecodeHook:       decoderHook,
		Result:           cfg,
		WeaklyTypedInput: true,
	}

	decoder, err := mapstructure.NewDecoder(config)
	if err != nil {
		return err
	}
	err = decoder.Decode(d)
	if err != nil {
		return err
	}
	return nil
}

func cleanEnv() error {
	err := os.Unsetenv(api.EnvVaultToken)
	if err != nil {
		return err
	}

	err = os.Unsetenv(api.EnvVaultAddress)
	if err != nil {
		return err
	}
	return nil
}

func loadToken(vault *api.Client) error {
	if vault == nil {
		return fmt.Errorf("the vault is nil")
	}
	basePath := getVaultPath()

	tokenSec, err := vault.Logical().Read(fmt.Sprintf("%s/token", basePath))
	if err != nil {
		return err
	}

	d, ok := tokenSec.Data["data"].(map[string]interface{})
	if !ok {
		return fmt.Errorf("no data in vault")
	}

	token, ok := d["token"].(string)
	if !ok {
		return fmt.Errorf("can't get the access token")
	}

	vault.SetToken(token)
	return nil
}

func getVaultPath() string {
	deploy := os.Getenv(DEPLOYMENT_NAME)
	ns := os.Getenv(KUBE_NAMESPACE)
	return fmt.Sprintf("%s/%s/data", deploy, ns)
}

// Init loads config data from files or from ENV
func Init(cfg interface{}, filename *string) error {

	//check if file with config exists
	if _, err := os.Stat(*filename); os.IsNotExist(err) {
		//expected file with config not exists
		//check special /.secrets directory (DevOps special)
		developmentConfigPath := "/.secrets/config.json"
		if _, err := os.Stat(developmentConfigPath); os.IsNotExist(err) {
			return err
		}

		filename = &developmentConfigPath
	}

	file, err := os.Open(*filename)
	if err != nil {
		return err
	}

	return Load(cfg, file)

}

// Load loads config data from any reader or from ENV
func Load(cfg interface{}, source io.Reader) error {
	decoder := json.NewDecoder(source)
	err := decoder.Decode(&cfg)
	if err != nil {
		return err
	}

	return overrideConfigWithEnvVars(cfg)
}

//Init loads configuration from the provided byte array.
func InitFromByteArray(cfg interface{}, buf []byte) error {
	return Load(cfg, bytes.NewReader(buf))
}

//overrideConfigWithEnvVars Overrides configuration with env parameters if any provided in UNDERSCORED_UPPERCASED format
func overrideConfigWithEnvVars(cfg interface{}) error {
	configsMap := structs.Map(cfg)

	for key := range configsMap {
		envVal, isPresent := os.LookupEnv(helpers.UndescoreUppercased(key))
		if isPresent {
			configsMap[key] = envVal
		}
	}

	err := mapstructure.Decode(configsMap, &cfg)
	if err != nil {
		return err
	}

	return err
}

// ValidateBaseConfigStructs validates additional structures (which implements BaseConfig)
func ValidateBaseConfigStructs(cfg interface{}) (err error) {
	v := reflect.ValueOf(cfg).Elem()
	baseConfigType := reflect.TypeOf((*BaseConfig)(nil)).Elem()

	for i := 0; i < v.NumField(); i++ {
		if v.Type().Field(i).Type.Implements(baseConfigType) {
			err = v.Field(i).Interface().(BaseConfig).Validate()
			if err != nil {
				return
			}
		}
	}

	return
}
