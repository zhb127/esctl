package cmd

import (
	"os"
	"strconv"

	defaults "github.com/mcuadros/go-defaults"
	homedir "github.com/mitchellh/go-homedir"
	flag "github.com/spf13/pflag"

	"github.com/pkg/errors"
	"github.com/spf13/viper"
)

type config struct {
	Log                Log                  `mapstructure:"log"`
	Clusters           []configClustersItem `mapstructure:"clusters"`
	Users              []configUsersItem    `mapstructure:"users"`
	Contexts           []configContextsItem `mapstructure:"contexts"`
	CurrentContextName string               `mapstructure:"current-context"`
	CurrentClusterName string
	CurrentUserName    string
	CurrentUserItem    *configUsersItem
	CurrentClusterItem *configClustersItem
	CurrentContextItem *configContextsItem
}

type Log struct {
	Level  string `mapstructure:"level" default:"debug"`
	Format string `mapstructure:"format" default:"pretty"`
}

type configClustersItem struct {
	Name    string                    `mapstructure:"name"`
	Cluster configClustersItemCluster `mapstructure:"cluster"`
}

type configClustersItemCluster struct {
	Addresses string `mapstructure:"addresses"`
}

type configUsersItem struct {
	Name string              `mapstructure:"name"`
	User configUsersItemUser `mapstructure:"user"`
}

type configUsersItemUser struct {
	Username   string `mapstructure:"username"`
	Password   string `mapstructure:"password"`
	CertVerify bool   `mapstructure:"cret-verify"`
	CertData   string `mapstructure:"cret-data"`
}

type configContextsItem struct {
	Name    string              `mapstructure:"name"`
	Context ContextsItemContext `mapstructure:"context"`
}

type ContextsItemContext struct {
	Cluster string `mapstructure:"cluster"`
	User    string `mapstructure:"user"`
}

func initConfig(cfgFilePath string, pFlags *flag.FlagSet) (*config, error) {
	if cfgFilePath == "" {
		home, err := homedir.Dir()
		if err != nil {
			return nil, errors.Wrap(err, "failed to return home dir")
		}
		cfgFileDir := home + "/.esctl"
		cfgFilePath = cfgFileDir + "/config"

		if _, err := os.Stat(cfgFilePath); err != nil {
			if !os.IsNotExist(err) {
				return nil, errors.Wrap(err, "failed to check config file state")
			}
			// 生成示例配置
			if err := genExampleConfig(cfgFileDir, cfgFilePath); err != nil {
				return nil, errors.Wrap(err, "failed to generate example config")
			}
		}
	}

	viper.SetConfigFile(cfgFilePath)
	viper.SetConfigType("yaml")

	if err := viper.ReadInConfig(); err != nil {
		return nil, errors.Wrap(err, "failed to read config file err")
	}

	// 设置默认配置值
	cfg := &config{}
	defaults.SetDefaults(cfg)

	if err := viper.Unmarshal(cfg); err != nil {
		return nil, errors.Wrap(err, "failed to parse config file err")
	}

	if err := mergeFlagsToConfig(pFlags, cfg); err != nil {
		return nil, errors.Wrap(err, "failed to convert rootCmd pFlags to config")
	}

	if err := validateConfig(cfg); err != nil {
		return nil, errors.Wrap(err, "failed to validate config")
	}

	if err := injectConfigToEnvVars(cfg); err != nil {
		return nil, errors.Wrap(err, "failed to inject config to env vars")
	}

	return cfg, nil
}

func genExampleConfig(dir string, path string) error {
	if err := os.MkdirAll(dir, os.ModePerm); err != nil {
		return err
	}

	f, err := os.Create(path)
	if err != nil {
		return err
	}

	defer f.Close()

	_, err = f.Write([]byte(`# example config
log:
  level: debug
  format: pretty
clusters:
- name: localhost
  cluster:
    addresses: http://localhost:9200
users:
- name: localhost
  user:
    username:
    password:
    certVerify:
    certContent:
contexts:
- name: localhost
  context:
    cluster: localhost
    user: localhost
current-context: localhost
`))

	return err
}

func mergeFlagsToConfig(pFlags *flag.FlagSet, cfg *config) error {
	flagCluster, err := pFlags.GetString("cluster")
	if err != nil {
		return err
	}
	if flagCluster != "" {
		cfg.CurrentClusterName = flagCluster
	}

	flagUser, err := pFlags.GetString("user")
	if err != nil {
		return err
	}
	if flagUser != "" {
		cfg.CurrentUserName = flagUser
	}

	flagContext, err := pFlags.GetString("context")
	if err != nil {
		return err
	}
	if flagContext != "" {
		cfg.CurrentContextName = flagContext
	}

	return nil
}

func validateConfig(cfg *config) error {
	if cfg == nil {
		return errors.New("config is nil")
	}

	if cfg.CurrentContextName == "" {
		return errors.New("config.current-context is empty")
	}
	if len(cfg.Clusters) == 0 {
		return errors.New("config.clusters is empty")
	}
	if len(cfg.Contexts) == 0 {
		return errors.New("config.contexts is empty")
	}
	if len(cfg.Users) == 0 {
		return errors.New("config.users is empty")
	}

	for _, v := range cfg.Contexts {
		if v.Name == cfg.CurrentContextName {
			cfg.CurrentContextItem = &v
			if cfg.CurrentClusterName == "" {
				cfg.CurrentClusterName = v.Context.Cluster
			}
			if cfg.CurrentUserName == "" {
				cfg.CurrentUserName = v.Context.User
			}
			break
		}
	}
	if cfg.CurrentContextItem == nil {
		return errors.New("config.contexts not contains config.current-context")
	}

	for _, v := range cfg.Clusters {
		if v.Name == cfg.CurrentClusterName {
			cfg.CurrentClusterItem = &v
			break
		}
	}
	if cfg.CurrentClusterItem == nil {
		return errors.New("config.clusters not contains config.current-cluster")
	}

	for _, v := range cfg.Users {
		if v.Name == cfg.CurrentUserName {
			cfg.CurrentUserItem = &v
			break
		}
	}
	if cfg.CurrentUserItem == nil {
		return errors.New("config.users not contains config.current-user")
	}

	return nil
}

func injectConfigToEnvVars(cfg *config) error {
	if err := os.Setenv("LOG_LEVEL", cfg.Log.Level); err != nil {
		return err
	}

	if err := os.Setenv("LOG_FORMAT", cfg.Log.Format); err != nil {
		return err
	}

	if err := os.Setenv("ES_ADDRESSES", cfg.CurrentClusterItem.Cluster.Addresses); err != nil {
		return err
	}

	if err := os.Setenv("ES_USERNAME", cfg.CurrentUserItem.User.Username); err != nil {
		return err
	}

	if err := os.Setenv("ES_PASSWORD", cfg.CurrentUserItem.User.Password); err != nil {
		return err
	}

	if err := os.Setenv("ES_CERT_DATA", cfg.CurrentUserItem.User.CertData); err != nil {
		return err
	}

	if err := os.Setenv("ES_CERT_VERIFY", strconv.FormatBool(cfg.CurrentUserItem.User.CertVerify)); err != nil {
		return err
	}

	return nil
}
