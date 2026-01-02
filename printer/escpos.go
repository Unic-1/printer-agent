package printer

func BuildEscPos(text string, cut bool) []byte {
	buf := []byte{}

	// Initialize
	buf = append(buf, 0x1B, 0x40)

	// Text
	buf = append(buf, []byte(text)...)
	buf = append(buf, '\n')

	// Feed
	buf = append(buf, 0x1B, 0x64, 0x03)

	// Cut
	if cut {
		buf = append(buf, 0x1D, 0x56, 0x00)
	}

	return buf
}
