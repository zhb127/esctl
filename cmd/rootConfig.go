package cmd

import "errors"

type config struct {
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
	CertVerify bool   `mapstructure:"cretVerify"`
	CertData   string `mapstructure:"cretData"`
}

type configContextsItem struct {
	Name    string              `mapstructure:"name"`
	Context ContextsItemContext `mapstructure:"context"`
}

type ContextsItemContext struct {
	Cluster string `mapstructure:"cluster"`
	User    string `mapstructure:"user"`
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
			cfg.CurrentClusterName = v.Context.Cluster
			cfg.CurrentUserName = v.Context.User
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

	for _, v := range cfg.Users {
		if v.Name == cfg.CurrentUserName {
			cfg.CurrentUserItem = &v
			break
		}
	}

	return nil
}
