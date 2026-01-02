package printer

type KOTItem struct {
	Name string
	Qty  int
}

type KOT struct {
	Restaurant string
	OrderNo    string
	Table      string
	Items      []KOTItem
	Notes      string
	QR         string
	Logo       []byte
}

func BuildKOT(p PrinterProfile, kot KOT) []byte {
	b := NewBuilder(p)

	// Logo
	if kot.Logo != nil {
		b.Logo(kot.Logo)
	}

	// Header
	b.AlignCenter()
	b.Bold(true)
	b.Double(true)
	b.Text("KOT")
	b.Double(false)
	b.Bold(false)

	b.Text(kot.Restaurant)
	b.Line()

	b.AlignLeft()
	b.Row("Order", kot.OrderNo)
	b.Row("Table", kot.Table)
	b.Line()

	// Items
	b.Bold(true)
	b.Text("ITEM            QTY")
	b.Bold(false)

	for _, item := range kot.Items {
		b.Row(item.Name, string(rune('0'+item.Qty)))
	}

	b.Line()

	if kot.Notes != "" {
		b.Bold(true)
		b.Text("Notes:")
		b.Bold(false)
		b.Text(kot.Notes)
	}

	// QR
	if kot.QR != "" {
		b.QR(kot.QR)
	}

	b.Cut()
	return b.Bytes()
}
