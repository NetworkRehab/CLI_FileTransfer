package main

import (
    "github.com/spf13/viper"
    "log"
)

// initConfig initializes the configuration by reading the config file.
func initConfig() {
    viper.SetConfigName("config")
    viper.SetConfigType("yaml")
    viper.AddConfigPath(".")

    // Default configurations
    viper.SetDefault("azure.accountName", "")
    viper.SetDefault("azure.accountKey", "")
    viper.SetDefault("azure.containerName", "")
    viper.SetDefault("aws.bucket", "")
    viper.SetDefault("aws.region", "")
    viper.SetDefault("cifs.mountPoint", "")
    viper.SetDefault("sftp.host", "")
    viper.SetDefault("sftp.port", "22")
    viper.SetDefault("sftp.username", "")
    viper.SetDefault("sftp.password", "")

    if err := viper.ReadInConfig(); err != nil {
        log.Fatalf("Error reading config file: %v", err)
    }
}
