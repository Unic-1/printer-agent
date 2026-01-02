package printer

type PrinterProfile struct {
	CharsPerLine int
}

var Profile58 = PrinterProfile{CharsPerLine: 32}
var Profile80 = PrinterProfile{CharsPerLine: 48}
