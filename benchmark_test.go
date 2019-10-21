package ahocorasick_test

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"github.com/sepetrov/ahocorasick"
)

var benchmarks = []struct {
	lang      string
	textSizes []int
}{
	{
		"bg",
		[]int{1, 10, 100},
	},
	{
		"de",
		[]int{1, 10, 100},
	},
	{
		"en",
		[]int{1, 10, 100},
	},
	{
		"ru",
		[]int{1, 10, 100},
	},
	{
		"sv",
		[]int{1, 10, 100},
	},
	{
		"zh",
		[]int{1, 10, 100},
	},
}

func Benchmark(b *testing.B) {
	for _, bb := range benchmarks {
		b.Run(bb.lang, func(b *testing.B) {
			var dict []string
			if d, err := readDictionary(filepath.Join("testdata", fmt.Sprintf("dictionary.%s.txt", bb.lang))); err != nil {
				b.Fatal(err)
			} else {
				dict = bytes2Strings(d)
			}

			b.Run("indexing", func(b *testing.B) {
				b.ResetTimer()
				for n := 0; n < b.N; n++ {
					ahocorasick.New(dict)
				}
			})

			for _, size := range bb.textSizes {
				size := size
				b.Run(fmt.Sprintf("searching in %d MB text", size), func(b *testing.B) {
					trie := ahocorasick.New(dict)

					var txt string
					if t, err := readText(bb.lang, size); err != nil {
						b.Fatal(err)
					} else {
						txt = string(t)
					}

					b.ResetTimer()
					for n := 0; n < b.N; n++ {
						trie.Search(txt)
					}
				})
			}
		})
	}
}

// readDictionary reads the dictionary file and returns its words.
func readDictionary(filename string) ([][]byte, error) {
	var dict [][]byte

	f, err := os.OpenFile(filename, os.O_RDONLY, 0660)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	r := bufio.NewReader(f)
	for {
		l, err := r.ReadBytes('\n')
		if err != nil || err == io.EOF {
			break
		}
		l = bytes.TrimSpace(l)
		dict = append(dict, l)
	}

	return dict, nil
}

// readText returns text for lang with size MB.
func readText(lang string, size int) ([]byte, error) {
	f := filepath.Join("testdata", fmt.Sprintf("text.%s.%dMB.txt", lang, size))
	b, err := ioutil.ReadFile(f)
	if err != nil {
		return nil, fmt.Errorf("read text from %q: %s", f, err)
	}
	return b, nil
}

// bytes2Strings converts b to slice of strings.
func bytes2Strings(b [][]byte) []string {
	var ss = make([]string, 0, len(b))
	for _, x := range b {
		ss = append(ss, string(x))
	}
	return ss
}
