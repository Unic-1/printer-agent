package printer

func (b *Builder) AlignLeft() {
	b.buf = append(b.buf, ESC_ALIGN_LEFT...)
}

func (b *Builder) AlignCenter() {
	b.buf = append(b.buf, ESC_ALIGN_CTR...)
}

func (b *Builder) Bold(on bool) {
	if on {
		b.buf = append(b.buf, ESC_BOLD_ON...)
	} else {
		b.buf = append(b.buf, ESC_BOLD_OFF...)
	}
}

func (b *Builder) Double(on bool) {
	if on {
		b.buf = append(b.buf, ESC_DOUBLE_ON...)
	} else {
		b.buf = append(b.buf, ESC_DOUBLE_OFF...)
	}
}

func (b *Builder) Text(t string) {
	b.buf = append(b.buf, []byte(t)...)
	b.buf = append(b.buf, '\n')
}
