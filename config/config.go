// Config is put into a different package to prevent cyclic imports in case
// it is needed in several locations

package config

import (
  "fmt"
  //"os"
  //"time"

  logp "github.com/elastic/beats/libbeat/logp"
  //"github.com/elastic/beats/libbeat/common"
)

type Config struct {
	Inputs []Input `config:"inputs"`
}

type Input struct {
  Type string `config:"type"`
  Name string `config:"name"`
  Period int `config:"period"`
  Paths []string `config:"paths"`
  Whitelist []string `config:"whitelist"`
  Blacklist []string `config:"blacklist"`
  Max_depth int `config:"max_depth"`
  Attribute string `config:"attribute"`

}

var DefaultConfig Config

func Validate(src []Input, dest []Input) (error) {
  var err error

  for _, si := range src {
    v_input := Input{
      Period: 60,
      Attribute: "mtime",
      Max_depth: 0,
    }

    // If there is no name that is bad, error
    if si.Name == "" {
      return fmt.Errorf("import specified without a name.")
    } else {
      v_input.Name = si.Name
    }

    //Type will default to monitor so that is okay. Error
    // on any invalid specification though
    valid_types := []string{"verbose_monitor", "monitor", ""}
    if Contains(v_input.Type, valid_types) {
      if si.Type != "" {
        v_input.Type = si.Type
      } else {
        v_input.Type = "monitor"
      }
    } else {
      return fmt.Errorf("Invalid type specified: %s", si.Type)
    }

    // Period can be undefined, if so 60 is default
    if si.Period != 0 {
      v_input.Period = si.Period
    }

    // There is a set of valid attributes and a default.
    // if it is set to anything else we're going to error
    valid_attributes := []string{"mtime", "ctime", "atime", ""}
    if Contains(v_input.Attribute, valid_attributes) {
      if si.Attribute != "" {
        v_input.Attribute = si.Attribute
      } else {
        v_input.Attribute = "mtime"
      }
    } else {
      return fmt.Errorf("Invalid attribute: %s", si.Attribute)
    }

    // Paths must not be empty
    if len(si.Paths) == 0 {
      logp.Err("No paths specified in input: %s", si.Name)
      return err
    } else {
      for _, path := range si.Paths{
        v_input.Paths = append(v_input.Paths, path)
      }
    }

    for _, rx := range si.Blacklist{
      v_input.Blacklist = append(v_input.Blacklist, rx)
    }

    for _, rx := range si.Whitelist{
      v_input.Whitelist = append(v_input.Whitelist, rx)
    }
    dest = append(dest, v_input)
  }
  return nil
}

func Contains(str string, list []string) bool {
  for _, v := range list {
    if v == str {
      return true
    }
  }
  return false
}
