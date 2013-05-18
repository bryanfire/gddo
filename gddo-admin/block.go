// Copyright 2013 Gary Burd
//
// Licensed under the Apache License, Version 2.0 (the "License"): you may
// not use this file except in compliance with the License. You may obtain
// a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS, WITHOUT
// WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied. See the
// License for the specific language governing permissions and limitations
// under the License.

package main

import (
	"github.com/garyburd/gddo/database"
	"log"
	"os"
)

var blockCommand = &command{
	name:  "block",
	run:   block,
	usage: "block path",
}

func block(c *command) {
	if len(c.flag.Args()) != 1 {
		c.printUsage()
		os.Exit(1)
	}
	db, err := database.New()
	if err != nil {
		log.Fatal(err)
	}
	if err := db.Block(c.flag.Args()[0]); err != nil {
		log.Fatal(err)
	}
}
