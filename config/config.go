// Licensed to Elasticsearch B.V. under one or more contributor
// license agreements. See the NOTICE file distributed with
// this work for additional information regarding copyright
// ownership. Elasticsearch B.V. licenses this file to you under
// the Apache License, Version 2.0 (the "License"); you may
// not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing,
// software distributed under the License is distributed on an
// "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
// KIND, either express or implied.  See the License for the
// specific language governing permissions and limitations
// under the License.

// Config is put into a different package to prevent cyclic imports in case
// it is needed in several locations

package config

import (
  "fmt"
  //"os"
  "time"

  logp "github.com/elastic/beats/libbeat/logp"
  //"github.com/elastic/beats/libbeat/common"
)

type Config struct {
	Inputs []Input `config:"inputs"`
}

type Input struct {
  Disabled bool `config:"disabled"`
  Name string `config:"name"`
  Period time.Duration `config:"period"`
  Threshold time.Duration `config:"threshold"`
  Paths []string `config:"paths"`
  Whitelist []string `config:"whitelist"`
  Blacklist []string `config:"blacklist"`
  Max_depth int `config:"max_depth"`
  Attribute string `config:"attribute"`
  Heartbeat bool `config:"heartbeat"`
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
    
    // Period can be undefined, if so 60 is default
    if si.Period != 0 {
      v_input.Period = si.Period
    }

    // Period can be undefined, if so 60 is the default
    if si.Threshold != 0 {
      v_input.Threshold = si.Threshold
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

    if len(v_input.Whitelist) > 0 && len(v_input.Blacklist) > 0 {
      return fmt.Errorf("It seems like an input config has both whitelist" +
                        "and blacklist are specified.")
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
