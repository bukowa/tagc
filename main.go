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
	"strings"
)

// input flags
var (
	InputSearch = flag.Bool("s", false, "perform search with tags")
	InputCmd = flag.String("c", "", "command to store")
	InputTags = flag.String("t", "", "tags to store separated by ','")
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
	Tags []string `json:"tags"`
}

// output
type Results struct {
	Matches map[int][]string
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

	if err = store(f, &data, *InputCmd, *InputTags); err != nil {
		log.Fatal(err)
	}
}

func store(f *os.File, data *Data, cmd, tags string) error {
	t := strings.Split(tags, ",")

	newCmd := &Command{
		Command: cmd,
		Tags:    t,
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

	var probability = make(map[*Command]int)

	for _, cmd := range data.Commands {
		for _, tag := range cmd.Tags {
			for _, searchTag := range tagSlice {
				if strings.ToLower(tag) == strings.ToLower(searchTag) {
					probability[cmd] += 1
				}
			}
		}
	}

	var commands[]*Command
	var number int

	for cmd, i := range probability {
		if i > number {
			number = i
			commands = []*Command{cmd}
			continue
		}
		if i == number {
			commands = append(commands, cmd)
			number = i
			continue
		}
	}

	var result = Results{Matches: make(map[int][]string)}

	for cmd, i := range probability {
		result.Matches[i] = append(result.Matches[i], cmd.Command)
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
