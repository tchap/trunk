// Copyright (c) 2014 The AUTHORS
//
// This file is part of trunk.
//
// trunk is free software: you can redistribute it and/or modify
// it under the terms of the GNU General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// trunk is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU General Public License for more details.
//
// You should have received a copy of the GNU General Public License
// along with trunk.  If not, see <http://www.gnu.org/licenses/>.

package shift

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"path/filepath"
	"regexp"

	"github.com/tchap/trunk/config"
	_ "github.com/tchap/trunk/config/autoload"
	"github.com/tchap/trunk/git"
)

const FullVersionPattern = "[0-9]+([.][0-9]+){2}"

type Versions struct {
	TrunkCurrent      string
	TrunkNext         string
	ReleaseCurrent    string
	ReleaseNext       string
	ProductionCurrent string
	ProductionNext    string
}

func ComputeVersions(next string) (vs *Versions, stderr *bytes.Buffer, err error) {
	if next == "auto" {
		var current string
		current, stderr, err = git.CurrentBranch()
		if err != nil {
			return
		}
		defer func() {
			errStream, ex := git.Checkout(current)
			if ex != nil {
				if err == nil {
					err = ex
				}
				if stderr == nil {
					stderr = errStream
				} else {
					if _, err := io.Copy(stderr, errStream); err != nil {
						panic(err)
					}
				}
			}
		}()

		stderr, err = git.Checkout(config.Local.TrunkBranch)
		if err != nil {
			return
		}

		next, err = readVersion()
		if err != nil {
			return
		}
		next = regexp.MustCompile(FullVersionPattern).FindString(next)
		if next == "" {
			err = fmt.Errorf("branch %v: package.json: malformed version string: %v",
				config.Local.TrunkBranch, next)
		}
	}
	return computeVersions(next)
}

func computeVersions(next string) (versions *Versions, stderr *bytes.Buffer, err error) {

}

func readVersion() (version string, err error) {
	root, _, err := git.RepositoryRootAbsolutePath()
	if err != nil {
		return
	}
	path := filepath.Join(root, "package.json")

	content, err := ioutil.ReadFile(path)
	if err != nil {
		return
	}

	var packageJson struct {
		Version string
	}
	err = json.Unmarshal(content, &packageJson)
	if err != nil {
		return
	}

	return packageJson.Version, nil
}
