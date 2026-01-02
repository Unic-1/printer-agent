package printer

type Builder struct {
	profile PrinterProfile
	buf     []byte
}

func NewBuilder(p PrinterProfile) *Builder {
	return &Builder{
		profile: p,
		buf:     append([]byte{}, ESC_INIT...),
	}
}

func (b *Builder) Bytes() []byte {
	return b.buf
}

func (b *Builder) Feed(lines int) {
	b.buf = append(b.buf, 0x1B, 0x64, byte(lines))
}

func (b *Builder) Cut() {
	b.buf = append(b.buf, ESC_CUT...)
}
