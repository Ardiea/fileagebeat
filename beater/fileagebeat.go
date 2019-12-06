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

package beater

import (
	"fmt"
	"time"
  "os"
  "strings"
  "path/filepath"
  "regexp"
  "runtime"

	"github.com/elastic/beats/libbeat/beat"
	"github.com/elastic/beats/libbeat/common"
	"github.com/elastic/beats/libbeat/logp"
	"github.com/Ardiea/fileagebeat/config"

  "gopkg.in/djherbis/times.v1"
)

// Fileagebeat configuration.
type Fileagebeat struct {
	done   chan struct{}
	config config.Config
	client beat.Client
  inputs []config.Input
}

// New creates an instance of fileagebeat.
func New(b *beat.Beat, cfg *common.Config) (beat.Beater, error) {
	c := config.DefaultConfig
	if err := cfg.Unpack(&c); err != nil {
		return nil, fmt.Errorf("Error reading config file: %v", err)
	}

	bt := &Fileagebeat{
		done:   make(chan struct{}),
		config: c,
	}

  //str := common.DebugString(bt.config, false)
  //fmt.Println(str)

  bt.inputs = make([]config.Input, 0)


  e := config.Validate(bt.config.Inputs, bt.inputs)
  if e != nil {
    return nil, e
  }

	return bt, nil
}

// Run starts fileagebeat.
func (bt *Fileagebeat) Run(b *beat.Beat) error {
	logp.Info("fileagebeat is running! Hit CTRL-C to stop it.")
	var err error
	bt.client, err = b.Publisher.Connect()
	if err != nil {
		return err
	}

  for _, input := range bt.config.Inputs {
    if ! input.Disabled {
      go SpawnCrawler(input, bt, b)
    }
  }

  <-bt.done
  return nil
}

func SpawnCrawler(input config.Input, bt *Fileagebeat, b *beat.Beat) {
  ticker := time.NewTicker(input.Period)
  //counter := 1
  logp.Info("Crawler: %s started.",input.Name)
  for {
    select {
    case <- bt.done:
      return
    case <-ticker.C:
    }

    // Build a fresh list of files every period
    files := BuildFileList(input)

    for _, f := range files {
      t := GetAge(f, input.Attribute)
      age := time.Now().Sub(t)
      if age > input.Threshold {
        fi, _ := os.Stat(f)
        t, err := times.Stat(f)

        if err != nil {
          logp.Err("Encountered an error collecting the times.")
          return
        } 

        atime := t.AccessTime()
        mtime := t.ModTime()
        ctime := t.ChangeTime()

        event := beat.Event{
          Timestamp: time.Now(),
          Fields: common.MapStr{
            "event": common.MapStr{
              "action": "aging_file_found",
            },
            "file": common.MapStr{
              "mtime": mtime,
              "atime": atime,
              "ctime": ctime,
              "size": fi.Size(),
              "mode": fi.Mode().String(),
              "path": f,
            },
            "agent": common.MapStr{
              "name": input.Name,
            },
          },
        }
        bt.client.Publish(event)
      }
    }

    if input.Heartbeat == true {
      event := beat.Event{
        Timestamp: time.Now(),
        Fields: common.MapStr{
          "event": common.MapStr{
             "action": "fileagebeat_heartbeat",
          },
          "agent": common.MapStr{
            "name": input.Name,
          },
        },
      }
      bt.client.Publish(event)
    }
  }
}

func GetAge(path string, attribute string) (val time.Time) {
  t, err := times.Stat(path)

  if err != nil {
    logp.Err("Encountered an error collecting the times.")
    return
  } 

  if attribute == "mtime" {
    val = t.ModTime()
  } else if attribute == "atime" {
    val = t.AccessTime()
  } else if attribute == "ctime" {
    val = t.ChangeTime()
  }
  return
}

// This returns a list of absolute paths that should have their age checked.
func BuildFileList(input config.Input) []string {
  var working_list []string
  for _, path := range input.Paths {
    err := filepath.Walk(path, func(p string, info os.FileInfo, err error) error{
      // Convert all the paths to relative paths
      // and remove the empty root path
      p = strings.Replace(p, path, "", -1)
      if len(p) > 1 {
        // input.Max_depth is set, remove violating paths from the list as
        // well as directories
        var delimiter string
        switch runtime.GOOS {
          case "windows":
            delimiter = "\\"
          default:
            delimiter = "/"
        }
        if input.Max_depth > 0 &&
           strings.Count(p, delimiter) <= input.Max_depth &&
           ! info.IsDir()  {
            working_list = append(working_list, path + p)
        } else if input.Max_depth == 0 &&
          ! info.IsDir() {
          working_list = append(working_list, path + p)
        }
      }
      return nil
    })
    if err != nil {
      panic(err)
    }
  }

  var files []string

  // The config parsing made sure whitelist and black list were mutually
  // exclusive. These won't both run.
  if len(input.Whitelist) > 0 {
    for _, r := range input.Whitelist {
      rx, _ := regexp.Compile(r)
      for _, f := range working_list {
        if rx.MatchString(f) {
          files = append(files, f)
        }
      }
    }
  }

  if len(input.Blacklist) > 0 {
    for _, r := range input.Blacklist {
      rx, _ := regexp.Compile(r)
      for _, f := range working_list {
        if ! rx.MatchString(f) {
          files = append(files, f)
        }
      }
    }
  }

  if len(input.Blacklist) == 0 && len(input.Whitelist) == 0 {
    files = append(files, working_list...)
  }

  return files
}

// Stop stops fileagebeat.
func (bt *Fileagebeat) Stop() {
	bt.client.Close()
	close(bt.done)
}
