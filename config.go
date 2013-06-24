package main

import (
    "code.google.com/p/goconf/conf"
    "log"
)

func getInt(c *conf.ConfigFile, section, option string) (value int) {
    value, err := c.GetInt(section, option)
    if err != nil {
        log.Fatal("missing config value: ", option)
    }
    return
}

func getString(c *conf.ConfigFile, section, option string) (value string) {
    value, err := c.GetString(section, option)
    if err != nil {
        log.Fatal("missing config value: ", option)
    }
    return
}

func getBool(c *conf.ConfigFile, section, option string) (value bool) {
    value, err := c.GetBool(section, option)
    if err != nil {
        log.Fatal("missing config value: ", option)
    }
    return
}
