package file

import (
	"bufio"
	"os"
)

func ReadRunes(filepath string) (chan rune, error) {
	ch := make(chan rune)
	file, err := os.Open(filepath)
	if err != nil {
		close(ch)
		return nil, err
	}

	go func() {
		defer file.Close()
		reader := bufio.NewReader(file)
		for {
			char, _, err := reader.ReadRune()
			if err != nil {
				break
			}
			ch <- char
		}
		close(ch)
	}()

	return ch, nil
}
