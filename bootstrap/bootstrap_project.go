package bootstrap

import (
	"bytes"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"text/template"

	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"

	"gopkg.in/src-d/go-billy.v4/osfs"
	"gopkg.in/src-d/go-git.v4"
	"gopkg.in/src-d/go-git.v4/storage/filesystem"
)

var gopath = os.Getenv("GOPATH")

// NewProject bootstraps a new project with the given name
func NewProject(fullName string) error {
	baseDir := filepath.Join(gopath, "src", fullName)
	path := strings.Split(baseDir, string(filepath.Separator))
	projectName := path[len(path)-1]
	log.Infof("bootstrapping new project named '%s' in '%s'...", projectName, baseDir)

	// first, ensure that the root directory of the project does not exist yet
	log.Debugf("checking that directory `%s` does not exist yet...", baseDir)
	if _, err := os.Stat(baseDir); err == nil {
		return errors.Errorf("'%s' already exists", baseDir)
	}
	if err := os.MkdirAll(baseDir, os.ModePerm); err != nil {
		return errors.Wrapf(err, "failed to make directory '%s'", baseDir)
	}
	for _, n := range AssetNames() {
		content, err := Asset(n)
		name, err := filepath.Rel("assets", n)
		if err != nil {
			return errors.Wrapf(err, "failed to obtain relative path of the '%s' asset", n)
		}
		filename := filepath.Join(baseDir, name)
		ext := filepath.Ext(filename)
		filecontent := content
		// process the asset as a golang text template instead of a static file if the file extension is `.tpl`
		if ext == ".tpl" {
			log.Infof("processing '%s' as a template...", filename)
			filename = strings.TrimSuffix(filename, ".tpl") // remove the `.tpl` extension
			log.Debugf(" -> '%s'", filename)
			tpl, err := template.New(n).Parse(string(content))
			if err != nil {
				return errors.Wrapf(err, "failed to parse the content of the '%s' template", n)
			}
			buf := bytes.NewBuffer(nil)
			err = tpl.Execute(buf, struct {
				ProjectName string
				MetricsName string
			}{
				ProjectName: projectName,
				MetricsName: strings.Replace(projectName, "-", "_", -1),
			})
			if err != nil {
				return errors.Wrapf(err, "failed to process '%v' template", n)
			}
			filecontent = buf.Bytes()
		}
		targetDir := filepath.Dir(filename)
		// make sure the target directory exists
		if _, err := os.Stat(targetDir); err != nil {
			log.Infof("directory '%s' does not seem to exists and thus it will be created", targetDir)
			if err := os.MkdirAll(targetDir, os.ModePerm); err != nil {
				return errors.Wrapf(err, "failed to make directory '%s'", targetDir)
			}
		}
		log.Infof("generating '%s'", filename)
		err = ioutil.WriteFile(filename, filecontent, 0644)
		if err != nil {
			return errors.Wrapf(err, "failed to write content of '%v'", n)
		}
	}
	// perform a `git init` at the root directory of the new project
	gitDir := osfs.New(filepath.Join(baseDir, ".git"))
	st, err := filesystem.NewStorage(gitDir)
	if err != nil {
		return errors.Wrapf(err, "failed to initialize a git repository at '%s'", baseDir)
	}
	_, err = git.Init(st, nil)
	if err != nil {
		return errors.Wrapf(err, "failed to initialize a git repository at '%s'", baseDir)
	}

	return nil
}
