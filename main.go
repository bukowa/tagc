package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"github.com/pkg/errors"
	"io"
	"log"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

// input flags
var (
	InputSearch = flag.Bool("s", false, "perform search with tags")
	InputCmd = flag.String("c", "", "text to store")
	InputTags = flag.String("t", "", "tags to store separated by ','")
	InputInfo = flag.String("i", "", "info about command")
)

// file paths
var (
	CmdFileDir = homeDir()
	CmdFilePath = ".tagc.commands.txt"
)

// file content
type Data struct{
	Commands []*Command `json:"commands"`
}

// stored input
type Command struct {
	Command string `json:"command"`
	Tags []string `json:"tags,omitempty"`
	Info string `json:"info,omitempty"`
}

// output
type Results struct {
	Output []map[int][]*Command
}

func main() {
	var data Data

	f, err := openCmdFile()
	defer func() {
		if err = f.Close(); err != nil {
			log.Fatal(errors.Wrap(err, "while closing file"))
		}
	}()

	if err != nil {
		log.Fatal(errors.Wrap(err, "while opening file"))
	}

	b, err := io.ReadAll(f)
	if err != nil {
		log.Fatal(errors.Wrap(err, "while reading file"))
	}

	err = json.Unmarshal(b, &data)
	if err != nil {
		log.Fatal(errors.Wrap(err, "while unmarshaling file"))
	}

	flag.Parse()

	if *InputTags == "" {
		log.Fatal("tags cannot be empty")
	}

	if *InputSearch {
		results, err := search(&data, *InputTags)
		if err != nil {
			log.Fatal(errors.Wrap(err, "while searching"))
		}
		b, err := json.MarshalIndent(results, "", " ")
		if err != nil {
			log.Fatal(errors.Wrap(err, "while converting results json"))
		}
		fmt.Print(string(b))
		return
	}

	if *InputCmd == "" {
		log.Fatal("command cannot be empty")
	}

	if err = store(f, &data, *InputCmd, *InputTags, *InputInfo); err != nil {
		log.Fatal(err)
	}
}

func store(f *os.File, data *Data, cmd, tags, info string) error {
	t := strings.Split(tags, ",")

	newCmd := &Command{
		Command: cmd,
		Tags:    t,
		Info: info,
	}

	if len(data.Commands) == 0 {
		data.Commands = []*Command{newCmd}
		goto save
	}

	for i, c := range data.Commands {
		if c.Command == cmd {
			data.Commands[i] = &Command{
				Command: cmd,
				Tags:    t,
			}
			goto save
		} else {
			data.Commands = append(data.Commands, newCmd)
			goto save
		}
	}

	save:
	b, err := json.MarshalIndent(data, "", "\t")
	if err != nil {
		return errors.Wrap(err, "4455")
	}

	err = f.Truncate(0)
	if err != nil {
		return errors.Wrap(err, "while truncate file")
	}

	_, err = f.Seek(0,0)
	if err != nil {
		return errors.Wrap(err, "while seek file")
	}

	_, err = f.Write(b)
	if err != nil {
		log.Print("you may want to save that")
		var b []byte
		err = json.Unmarshal(b, data)
		if err != nil {
			log.Fatal(errors.Wrap(err, "ups something went really wrong"))
		}
		log.Fatal(string(b))
	}

	return nil
}

func search(data *Data, tags string) (Results, error) {
	var tagSlice = strings.Split(tags, ",")

	var matchSlice = map[*Command]int{}

	for _, cmd := range data.Commands {
		for _, tag := range cmd.Tags {
			for _, searchTag := range tagSlice {
				if strings.ToLower(tag) == strings.ToLower(searchTag) {
					matchSlice[cmd] += 1
				}
			}
		}
	}

	var nN []int
	for _, i := range matchSlice {
		nN = append(nN, i)
	}

	var unique = make(map[int]struct{})
	for _, i := range nN {
		if _, ok := unique[i]; !ok {
			unique[i] = struct{}{}
		}
	}

	nN = []int{}
	for i, _ := range unique {
		nN = append(nN, i)
	}

	sort.Ints(nN)
	//sort.Sort(sort.Reverse(sort.IntSlice(nN)))

	var output = make([]map[int][]*Command, len(nN))
	for cmd, c := range matchSlice {
		for index, number := range nN {
			if number == c {
				if v, ok := output[index][number]; ok {
					output[index][number] = append(v, cmd)
					continue
				}
				output[index] = map[int][]*Command{number: {cmd}}
				continue
			}

		}
	}

	var result = Results{
		Output: output,
	}

	for _, v := range result.Output {
		for _, commands := range v {
			for _, c := range commands {
				c.Tags = nil
			}
		}
	}
	return result, nil
}

func openCmdFile() (*os.File, error) {
	path := filepath.Join(CmdFileDir, CmdFilePath)

	isNewFile := false
	f, err := os.Open(path)

	if err != nil && os.IsNotExist(err) {
		isNewFile = true
	}

	f, err = os.OpenFile(path, os.O_CREATE, 0600)
	if err != nil {
		return nil, errors.Wrap(err, "6765546")
	}

	if isNewFile {
		b, err := json.MarshalIndent(&Data{
			Commands: []*Command{},
		}, "", "\t")
		if err != nil {
			return nil, errors.Wrap(err, "43543545x")
		}
		if _, err = f.Write(b); err != nil {
			return nil, errors.Wrap(err, "435345c")
		}
	}

	if err := f.Close(); err != nil {
		return nil, errors.Wrap(err, "435435435")
	}

	f, err = os.OpenFile(path, os.O_CREATE, 0600)
	if err != nil {
		return nil, errors.Wrap(err, "asdasd435435")
	}

	return f, err

}

func homeDir() string {
	if h := os.Getenv("HOME"); h != "" {
		return h
	}
	return os.Getenv("USERPROFILE")
}
