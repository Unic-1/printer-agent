package printer

import (
	"strings"
)

func (b *Builder) Line() {
	for i := 0; i < b.profile.CharsPerLine; i++ {
		b.buf = append(b.buf, '-')
	}
	b.buf = append(b.buf, '\n')
}

func (b *Builder) Row(left, right string) {
	left = strings.TrimSpace(left)
	right = strings.TrimSpace(right)

	space := b.profile.CharsPerLine - len(left) - len(right)
	if space < 1 {
		space = 1
	}

	b.buf = append(b.buf, []byte(left)...)
	for i := 0; i < space; i++ {
		b.buf = append(b.buf, ' ')
	}
	b.buf = append(b.buf, []byte(right)...)
	b.buf = append(b.buf, '\n')
}

func (b *Builder) Logo(logo []byte) {
	if len(logo) == 0 {
		return
	}

	b.AlignCenter()
	b.buf = append(b.buf, logo...)
	b.Feed(1)
}
