package beater

import (
	"fmt"
	"time"

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

  var one_second time.Duration = 1 * time.Second
  ticker := time.NewTicker(one_second)
	counter := 0

	for {
	 	select {
	 	  case <-bt.done:
	 		  return nil
	 	  case <-ticker.C:
	  }

//    for _, input := range bt.config.Inputs {
//      if input.Enabled {
//        fmt.Println("%s is enabled: %s", input.Name, input.Enabled)
//      }
//    }


	  event := beat.Event{
		  Timestamp: time.Now(),
		  Fields: common.MapStr{
			  "type":    b.Info.Name,
			  "counter": counter,
		  },
	  }
	  bt.client.Publish(event)
	  logp.Info("Event sent")
	  counter++
	}
  return nil
}

// Stop stops fileagebeat.
func (bt *Fileagebeat) Stop() {
	bt.client.Close()
	close(bt.done)
}
