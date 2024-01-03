package gofilelist

import (
	"fmt"
	"io"
	"os"
	"runtime"
	"strings"
	"time"
)

/*
*
* gofilelist
* File List manager package on Golang
* Created by Jonathan G. Hecl
* https://github.com/jonathanhecl/gofilelist
*
 */

const (
	IsWindows        = runtime.GOOS == "windows"
	CommentStyle     = "//"
	LineBreak        = "\r"
	LineBreakWindows = "\r\n"
)

type Item struct {
	Value   string
	Comment string
}

type FileList struct {
	items        []Item
	lastModified time.Time
	changed      bool
	itemMap      map[string]struct{}
}

func New() *FileList {
	return &FileList{
		lastModified: time.Now(),
		changed:      true,
		itemMap:      make(map[string]struct{}),
	}
}

func (f *FileList) Count() int {
	return len(f.items)
}

func (f *FileList) LastModified() time.Time {
	return f.lastModified
}

func (f *FileList) Changed() bool {
	return f.changed
}

func Load(filename string) (*FileList, error) {
	f := New()
	lines, err := readFile(filename)
	if err != nil {
		return nil, err
	}

	for _, line := range lines {
		if item := validLine(line); item.Value != "" {
			f.items = append(f.items, item)
		}
	}

	f.changed = false
	return f, nil
}

func validLine(line string) (item Item) {
	item = Item{}
	line = strings.TrimSpace(line)
	if len(line) > 0 {
		if strings.HasPrefix(line, CommentStyle) {
			return
		}
		if strings.Contains(line, CommentStyle) {
			split := strings.Split(line, CommentStyle)
			if len(split) > 1 {
				item.Value = strings.TrimSpace(split[0])
				item.Comment = strings.TrimSpace(split[1])
			}
		} else {
			item.Value = line
		}
	}
	return
}

func makeLine(item Item) string {
	if item.Value == "" {
		return ""
	}

	line := item.Value
	if len(item.Comment) > 0 {
		line += fmt.Sprintf("\t"+CommentStyle+"%s", item.Comment)
	}

	return line
}

func (f *FileList) Save(filename string) error {
	if file, err := os.Create(filename); err != nil {
		return err
	} else {
		defer file.Close()

		lineBreak := LineBreak
		if IsWindows {
			lineBreak = LineBreakWindows
		}

		for i := range f.items {
			line := makeLine(f.items[i]) + lineBreak

			if _, err := file.Write([]byte(line)); err != nil {
				return err
			}
		}

		f.changed = false
	}

	return nil
}

func readFile(filename string) ([]string, error) {
	f, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	var (
		buf   []byte = make([]byte, 32*1024)
		lines []string
		line  []byte = []byte{}
	)
	for {
		n, err := f.Read(buf)
		if n > 0 {
			for i := 0; i < n; i++ {
				if buf[i] == 10 ||
					buf[i] == 13 {
					if len(line) > 0 {
						lines = append(lines, string(line))
						line = []byte{}
					}
				} else {
					line = append(line, buf[i])
				}
			}
		}
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, fmt.Errorf("read %d bytes: %v", n, err)
		}
	}
	if len(line) > 0 {
		lines = append(lines, string(line))
	}
	return lines, nil
}

func (f *FileList) SetItems(items []Item) {
	f.items = items
	f.lastModified = time.Now()
	f.changed = true
}

func (f *FileList) GetItems() []Item {
	return f.items
}

func (f *FileList) Get(value string) Item {
	for _, v := range f.items {
		if v.Value == value {
			return v
		}
	}
	return Item{}
}

func (f *FileList) GetComment(value string) string {
	for _, v := range f.items {
		if v.Value == value {
			return v.Comment
		}
	}
	return ""
}

func (f *FileList) GetAllWithComment(comment string) []Item {
	var items []Item
	for _, v := range f.items {
		if v.Comment == comment {
			items = append(items, v)
		}
	}
	return items
}

func (f *FileList) Exists(value string) bool {
	_, exists := f.itemMap[value]
	return exists
}

func (f *FileList) AddOnce(value string, comment string) {
	if _, exists := f.itemMap[value]; exists {
		for i, v := range f.items {
			if v.Value == value {
				if v.Comment != comment {
					f.items[i].Comment = comment
					f.lastModified = time.Now()
					f.changed = true
				}
				return
			}
		}
	} else {
		f.items = append(f.items, Item{Value: value, Comment: comment})
		f.itemMap[value] = struct{}{}
		f.lastModified = time.Now()
		f.changed = true
	}
}

func (f *FileList) Add(value string, comment string) {
	f.items = append(f.items, Item{Value: value, Comment: comment})
	f.itemMap[value] = struct{}{}
	f.lastModified = time.Now()
	f.changed = true
}

func (f *FileList) Remove(item string) {
	for i, v := range f.items {
		if v.Value == item {
			f.items = append(f.items[:i], f.items[i+1:]...)
			delete(f.itemMap, item)
			f.lastModified = time.Now()
			f.changed = true
			return
		}
	}
}
