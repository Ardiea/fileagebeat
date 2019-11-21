package beater

import (
	"fmt"
	"time"
  "os"
  "strings"
  "path/filepath"

	"github.com/elastic/beats/libbeat/beat"
	"github.com/elastic/beats/libbeat/common"
	"github.com/elastic/beats/libbeat/logp"
	"github.com/Ardiea/fileagebeat/config"
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

  for _, input := range bt.inputs {
    fmt.Println(input)
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
    go SpawnCrawler(input, bt, b)
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

    files := BuildFileList(input)

    for _, f := range files {
      fmt.Println(f)
    }

    event := beat.Event{
      Timestamp: time.Now(),
      Fields: common.MapStr{
        "event": "eventname",
      },
    }
    bt.client.Publish(event)
  }
}

// This returns a list of files within the paths for an input.
func BuildFileList(input config.Input) []string {
  var files []string
  for _, path := range input.Paths {
    err := filepath.Walk(path, func(p string, info os.FileInfo, err error) error{
      // Convert all the paths to relative paths
      // and remove the empty root path
      p = strings.Replace(p, path, "", -1)
      if len(p) > 1 {
        // input.Max_depth is set, remove violating paths from the list as
        // well as directories
        if input.Max_depth > 0 &&
           strings.Count(p, "/") <= input.Max_depth &&
           ! info.IsDir() {
            files = append(files, p)
        }
      }
      return nil
    })
    if err != nil {
      panic(err)
    }
  }
  return files
}

// Stop stops fileagebeat.
func (bt *Fileagebeat) Stop() {
	bt.client.Close()
	close(bt.done)
}
