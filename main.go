package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"

	"github.com/urfave/cli"
)

func main() {
	app := cli.NewApp()
	app.Name = "tree"
	app.Version = "1.0.0"
	app.Usage = "tree"
	app.Action = func(c *cli.Context) error {
		level := c.Int("level")
		return tree(c.Args().Get(0), level)
	}

	app.Flags = []cli.Flag{
		cli.IntFlag{
			Name:  "level, L",
			Usage: "Descend only level directories deep.",
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}

func tree(path string, level int) error {
	if path == "" {
		path = "."
	}

	file, err := os.Stat(path)
	if err != nil {
		return err
	}

	if !file.IsDir() {
		return fmt.Errorf("%v is not directory", path)
	}
	fmt.Println(file.Name())
	return walk(file.Name(), level, 0, "")
}

func walk(path string, level int, depth int, prefixBase string) error {
	files, err := ioutil.ReadDir(path)
	if err != nil {
		return err
	}
	if level > 0 && depth > level {
		return nil
	}
	for i, f := range files {
		if f.Name() == "." || f.Name() == ".." {
			continue
		}
		isLast := i == len(files)-1
		if f.IsDir() {
			print(f.Name(), isLast, prefixBase)
			prefixChild := prefixBase
			if !isLast {
				prefixChild += "│   "
			} else {
				prefixChild += "    "
			}
			walk(filepath.Join(path, f.Name()), level, depth+1, prefixChild)
			continue
		}
		print(f.Name(), isLast, prefixBase)
	}

	return nil
}

func print(filename string, isLast bool, prefixBase string) {
	prefix := prefixBase
	if isLast {
		prefix += "└── "
	} else {
		prefix += "├── "
	}

	fmt.Println(prefix + filename)
}
