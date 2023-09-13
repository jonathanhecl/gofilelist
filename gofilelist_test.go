package gofilelist

import (
	"reflect"
	"testing"
)

func Test_New(t *testing.T) {
	f := New()
	if f == nil {
		t.Error("FileList is nil")
	}

	f.Add("test", "")
	f.AddOnce("test", "")

	if len(f.items) != 1 {
		t.Error("FileList.items is not 1")
	}

	if !f.Changed() {
		t.Error("FileList.Changed() is false")
	}
}

func Test_Save(t *testing.T) {
	f := New()
	f.Add("test", "")
	f.AddOnce("test", "")
	f.Add("test2", "comment")

	f.Add("test3", "")
	f.Remove("test3")

	err := f.Save("test.txt")
	if err != nil {
		t.Error(err)
	}
}

func Test_Load(t *testing.T) {
	f, err := Load("test.txt")
	if err != nil {
		t.Error(err)
	}

	if len(f.items) != 2 {
		t.Error("FileList.items is not 2")
	}

	if f.items[0].Value != "test" {
		t.Errorf("Expected test, got %s", f.items[0].Value)
	}

	if f.items[1].Value != "test2" {
		t.Errorf("Expected test2, got %s", f.items[1].Value)
	}

	if f.items[1].Comment != "comment" {
		t.Errorf("Expected comment, got %s", f.items[1].Comment)
	}
}

func Test_makeLine(t *testing.T) {
	type args struct {
		item Item
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			"Line without comment",
			args{Item{"test", ""}},
			"test",
		},
		{
			"Line with comment",
			args{Item{"test", "comment"}},
			"test\t//comment",
		},
		{
			name: "Ignore comment",
			args: args{Item{"", "ignore it"}},
			want: "",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := makeLine(tt.args.item); got != tt.want {
				t.Errorf("makeLine() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_validLine(t *testing.T) {
	type args struct {
		line string
	}
	tests := []struct {
		name     string
		args     args
		wantItem Item
	}{
		{
			"Line without comment",
			args{"test"},
			Item{"test", ""},
		},
		{
			"Line with comment",
			args{"test\t//comment"},
			Item{"test", "comment"},
		},
		{
			"Comment",
			args{"//comment"},
			Item{"", ""},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotItem := validLine(tt.args.line); !reflect.DeepEqual(gotItem, tt.wantItem) {
				t.Errorf("validLine() = %v, want %v", gotItem, tt.wantItem)
			}
		})
	}
}
