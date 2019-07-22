package main

import (
	"log"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
)

func outputToRegexp(in string) *regexp.Regexp {
	in = regexp.QuoteMeta(in)
	in = strings.Replace(in, `\{year\}`, `\d{2,}`, -1)
	in = strings.Replace(in, `\{month\}`, `[01]?\d`, -1)
	in = strings.Replace(in, `\{day\}`, `[0123]?\d`, -1)
	in = strings.Replace(in, `\{md5\}`, `[0-9a-fA-F]{8,}`, -1)
	return regexp.MustCompile(`^` + in + `$`)
}

func eliminateAllExcept(output string, limit int) error {
	path := filepath.Dir(output)
	f, err := os.Open(path)
	if err != nil {
		return err
	}
	list, err := f.Readdir(-1)
	f.Close()
	if err != nil {
		return err
	}

	filter := outputToRegexp(filepath.Base(output))
	total := 0
	for _, v := range list { // clean out the list
		if !filter.MatchString(v.Name()) {
			continue
		}
		list[total] = v
		total++
	}
	list = list[:total]

	if len(list) > limit {
		sort.Slice(list, func(i, j int) bool {
			return list[i].ModTime().Unix() > list[j].ModTime().Unix()
		})
		for _, fi := range list[limit:] {
			target := filepath.Join(path, fi.Name())
			err = os.Remove(target)
			if err != nil {
				return err
			}
			log.Printf(`There are more than %d matching files. Eliminated "%s".`, limit, target)
		}
	}
	return nil
}
