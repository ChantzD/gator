package config

import "os"
import "encoding/json"
import "fmt"

type Config struct {
	Db_url            string `json:"db_url"`
	Current_user_name string `json:"current_user_name"`
}

func Read() Config {
	home_dir, err := os.UserHomeDir()
	if err != nil {
		fmt.Println(err)
	}

	configFile, err := os.Open(home_dir + "/.gatorconfig.json")
	if err != nil {
		fmt.Println(err)
	}

	var config Config
	err = json.NewDecoder(configFile).Decode(&config)
	if err != nil {
		fmt.Println("Error unmarshaling json: ", err)
	}

	return config
}

func (c *Config) SetUser(username string) error {
	c.Current_user_name = username
	jsonData, err := json.Marshal(c)
	if err != nil {
		return err
	}

	home_dir, err := os.UserHomeDir()
	if err != nil {
		return err
	}

	err = os.WriteFile(home_dir+"/.gatorconfig.json", jsonData, 0644)
	if err != nil {
		return err
	}
	return nil
}
