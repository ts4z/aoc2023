package argv

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"log"
	"os"

	"github.com/ts4z/aoc2023/ick"
)

func ReadToChannels(lines chan<- string, errs chan<- error) {
	r := bufio.NewReader(Reader(func(filename string, err error) {
		errs <- fmt.Errorf("error reading %q: %w", filename, err)
	}))

	defer close(lines)
	defer close(errs)

	for {
		line, err := r.ReadString('\n')
		if err == nil {
			lines <- line
		} else {
			if errors.Is(err, io.EOF) {
				break
			}
			errs <- err
			break
		}
	}

}

// ReadAll is the moral equivalent of using <ARGV> in a list context in perl.
//
// This should be changed to read _any_ Reader into a []string.
func ReadAll() ([]string, error) {
	ech := make(chan error, 1)

	r := bufio.NewReader(Reader(func(filename string, err error) {
		ech <- fmt.Errorf("error reading %q: %w", filename, err)
	}))

	lines := []string{}
	for {
		select {
		case err := <-ech:
			return nil, err
		default:
			// carry on
		}
		line, err := r.ReadString('\n')
		if err != nil {
			if errors.Is(err, io.EOF) {
				break
			}
			log.Printf("err %v with %d lines in buffer", err, len(lines))
			return nil, err
		}
		if line[len(line)-1] != '\n' {
			return nil, errors.New("line didn't end in LF")
		}
		lines = append(lines, line)
	}

	close(ech)
	if err, ok := <-ech; ok {
		return nil, err
	}

	return lines, nil
}

// ReadChompAll is the moral equivalent of using <ARGV> in a list context,
// followed by chomping the input lines (with the canonical record separator).
func ReadChompAll() ([]string, error) {
	lines, err := ReadAll()
	if err != nil {
		return nil, err
	}
	chomped := ick.MapSlice(ick.ChompNL, lines)
	return chomped, nil
}

// Reader provides a Reader (filehandle) type interface similar to the Perl
// ARGV filehandle.  It inhales all of the data from the command-line named
// files, or if there aren't any, it consumes os.Stdin.
func Reader(onError func(filename string, err error)) io.Reader {
	if len(os.Args) == 1 {
		return os.Stdin
	}

	r, writer := io.Pipe()

	go func() {
		for _, filename := range os.Args[1:] {
			if reader, err := os.Open(filename); err != nil {
				onError(filename, err)
			} else {
				defer reader.Close()
				_, err := io.Copy(writer, reader)
				if err != nil {
					onError(filename, err)
				}
			}
		}
		defer writer.Close()
	}()

	return r
}

func ReaderLoggingErrors() io.Reader {
	return Reader(func(filename string, err error) {
		log.Printf("error handling file %q: %v", filename, err)
	})
}
