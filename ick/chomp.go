package ick

func ChompNL(line string) string {
	if line[len(line)-1] != '\n' {
		return line
	} else {
		return line[:len(line)-1] // chomp
	}
}
