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
  "io/ioutil"
  "math/rand"

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
  
  inputs, e := config.Validate(bt.config.Inputs)
  bt.inputs = inputs
  if e != nil {
    return nil, e
  }
  bt.inputs = inputs

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

  for _, input := range bt.inputs {
    if ! input.Disabled {
      go SpawnCrawler(input, bt, b)
    }
  }

  <-bt.done
  return nil
}

func SpawnCrawler(input config.Input, bt *Fileagebeat, b *beat.Beat) {
  rand.Seed(time.Now().UnixNano())
  delay_seconds := rand.Float64()*input.Period.Seconds()
  logp.Info("Crawler: %s started. Delaying collection for %f seconds.",input.Name, delay_seconds)
  time.Sleep(time.Duration(delay_seconds)*time.Second)

  ticker := time.NewTicker(input.Period)
  for {
    // Build a fresh list of files every period
    
    start := time.Now()
    files := BuildFileList2(input)
    elapsed := time.Since(start)
    logp.Debug("Time to build filelist for Crawler %s = %s", input.Name, elapsed)
    logp.Debug("Length of filelist for Crawler %s = %d", input.Name, len(files))

    for _, f := range files {
      logp.Err("Input is %s", input.Attribute)
      t := GetAge(f, input.Attribute)
      age := time.Now().Sub(t)
      if age > input.Threshold {
        fi, _ := os.Stat(f)
        t, err := times.Stat(f)

        if err != nil {
          logp.Err("Encountered an error collecting the times: %s", err.Error())
          return
        } 

        atime := t.AccessTime()
        mtime := t.ModTime()
        ctime := t.ChangeTime()
    
        var delimiter string
        switch runtime.GOOS {
          case "windows":
            delimiter = "\\"
          default:
            delimiter = "/"
        }

        directory := f[:strings.LastIndex(f, delimiter)]
        
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
              "directory": directory,
              "age": age.Round(time.Second).String(),
            },
            "agent": common.MapStr{
              "name": input.Name,
            },
          },
        }
        if len(bt.config.Fields) != 0 { 
          event.PutValue("fields", bt.config.Fields)
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
      if len(bt.config.Fields) != 0 { 
        event.PutValue("fields", bt.config.Fields)
      }
      bt.client.Publish(event)
    }
    select {
    case <- bt.done:
      return
    case <-ticker.C:
    }
  }
}

func GetAge(path string, attribute string) (val time.Time) {
  t, err := times.Stat(path)

  if err != nil {
    logp.Err("Encountered an error collecting the times: %s", err.Error())
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

// Performance enhanced version of the subroutine to build the filelist.
// No more costly filewalk. 
func BuildFileList2(input config.Input) (list []string) {
  var working_list []string
  max_depth := 128

  var delimiter string
  switch runtime.GOOS {
    case "windows":
      delimiter = "\\"
    default:
      delimiter = "/"
  }

  if input.Max_depth != 0 {
    max_depth = input.Max_depth
  }

  for _, path := range input.Paths {
    working_list = append(working_list, getFiles(path, delimiter, max_depth)...)
  }

  // The config parsing made sure whitelist and black list were mutually
  // exclusive. These won't both run. That was determined at parse time. 
  if len(input.Whitelist) > 0 {
    for _, r := range input.Whitelist {
      rx, _ := regexp.Compile(r)
      for _, f := range working_list {
        if rx.MatchString(f) {
          list = append(list, f)
        }
      }
    }
  }

  if len(input.Blacklist) > 0 {
    for _, r := range input.Blacklist {
      rx, _ := regexp.Compile(r)
      for _, f := range working_list {
        if ! rx.MatchString(f) {
          list = append(list, f)
        }
      }
    }
  }

  if len(input.Blacklist) == 0 && len(input.Whitelist) == 0 {
    list = append(list, working_list...)
  }
  return
}

// Recursive function to retrieve all the files down to a specified depth
func getFiles(path string, delimiter string, depth int) (list []string) {
  path = filepath.Clean(path)
  if depth == 0 {
    return
  } else {
    files, err := ioutil.ReadDir(path)
    if err != nil {
      logp.Err("Error producing filelist.")
    }
    for _, f := range(files) {
      if f.IsDir() {
        list = append(list, getFiles(path + delimiter + f.Name(), delimiter, depth-1)...)
      } else {
        list = append(list, path + delimiter + f.Name())
      }
    }
  }
  return
}

// Stop stops fileagebeat.
func (bt *Fileagebeat) Stop() {
	bt.client.Close()
	close(bt.done)
}
