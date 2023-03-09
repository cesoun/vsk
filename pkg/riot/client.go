package riot

import (
	"github.com/fsnotify/fsnotify"
	"log"
	"strings"
	"sync"
)

type channels struct {
	fsw chan struct{}
	lf  chan bool
	err chan error
}

type paths struct {
	lockfile string
	config   string
}

type Client struct {
	Active chan bool
	fsw    *fsnotify.Watcher
	lf     *Lockfile
	wg     sync.WaitGroup
	ch     channels
	paths  paths
}

// NewClient creates a new Client with default values
func NewClient() *Client {
	// get lockfile path
	lockfile, err := GetLockfilePath()
	if err != nil {
		log.Fatal(err)
	}

	// get config path
	config, err := GetConfigPath()
	if err != nil {
		log.Fatal(err)
	}

	return &Client{
		Active: make(chan bool, 1),
		fsw:    nil,
		lf:     nil,
		wg:     sync.WaitGroup{},
		ch: channels{
			fsw: make(chan struct{}),
			lf:  make(chan bool, 1),
		},
		paths: paths{
			lockfile: lockfile,
			config:   config,
		},
	}
}

// Start starts the client
func (c *Client) Start() {
	// setup error log loop
	go func() {
		c.ch.err = make(chan error)
		for {
			select {
			case err := <-c.ch.err:
				log.Println("[error]", err)
			}
		}
	}()

	go setupLockfileWatcher(c)
}

// Stop stops the client
func (c *Client) Stop() {
	log.Println("[stop] closing subroutines")

	c.ch.fsw <- struct{}{}

	c.wg.Wait()
	log.Println("[stop] all subroutines finished")

	close(c.ch.err)
}

// setupLockfileWatcher watches the config directory and updates the Client.lf when it changes.
func setupLockfileWatcher(c *Client) {
	// create watcher
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		c.ch.err <- err
	}

	// defer handler
	defer func(watcher *fsnotify.Watcher) {
		err := watcher.Close()
		if err != nil {
			c.ch.err <- err
		}

		log.Println("[fsnotify] closed")

		c.fsw = nil
		c.wg.Done()
	}(watcher)

	// Start listening
	go func() {
		for {
			select {
			case ev, ok := <-watcher.Events:
				if !ok {
					return
				}

				if strings.EqualFold(ev.Name, c.paths.lockfile) {
					if ev.Has(fsnotify.Remove) {
						c.lf = nil

						c.ch.lf <- false

						log.Println("[lockfile] removed")
					} else if ev.Has(fsnotify.Write) {
						c.lf, err = NewLockfile()
						if err != nil {
							c.ch.err <- err
						} else {
							c.ch.lf <- true
						}

						log.Println("[lockfile] written")
					}
				}
			case err, ok := <-watcher.Errors:
				if !ok {
					return
				}

				c.ch.err <- err
			}
		}
	}()

	c.fsw = watcher

	if err := c.fsw.Add(c.paths.config); err != nil {
		c.ch.err <- err
		return
	}

	// block until interrupted
	c.wg.Add(1)
	<-c.ch.fsw

	log.Println("[file-watcher] closing")
}
