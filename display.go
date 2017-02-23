package htmlPDF

import (
	"github.com/jung-kurt/gofpdf"
)

type DisplayCommand struct {
	command interface{}
}

type SolidColor struct {
	color Color
	rect  Rect
}

func (d DisplayCommand) draw(pdf *gofpdf.Fpdf) {
	switch command := d.command.(type) {
	case SolidColor:
		r := command.rect
		c := command.color
		pdf.SetFillColor(int(c.r), int(c.g), int(c.b))
		pdf.Rect(r.x, r.y, r.width, r.height, "F")
	}
}

func buildDisplayList(layoutRoot *LayoutBox) map[int]DisplayCommand {
	list := map[int]DisplayCommand{}
	renderLayoutBox(layoutRoot, list)
	return list
}

func renderLayoutBox(layoutBox *LayoutBox, list map[int]DisplayCommand) {
	//renderBackground
	backgroundCommand := renderBackground(layoutBox)
	if backgroundCommand != nil {
		list[len(list)] = *backgroundCommand
	}

	//renderBorders
	renderBorders(layoutBox, list)

	//TODO renderText

	//Render child
	for _, child := range layoutBox.children {
		renderLayoutBox(child, list)
	}
}

func renderBackground(layoutBox *LayoutBox) *DisplayCommand {
	colorBackrgound := getColor(layoutBox, "background")
	if colorBackrgound == nil {
		return nil
	}
	return &DisplayCommand{
		command: SolidColor{
			color: *colorBackrgound,
			rect:  layoutBox.dimensions.borderBox(),
		},
	}
}

func renderBorders(layoutBox *LayoutBox, list map[int]DisplayCommand) {
	colorBorder := getColor(layoutBox, "border-color")
	if colorBorder == nil {
		return
	}
	//Return if white
	//TODO change crete Color with nil
	if colorBorder.r == 255 && colorBorder.g == 255 && colorBorder.b == 255 {
		return
	}

	d := layoutBox.dimensions

	borderBox := d.borderBox()

	// Left border
	list[len(list)] = DisplayCommand{
		command: SolidColor{
			color: *colorBorder,
			rect: Rect{
				x:      borderBox.x,
				y:      borderBox.y,
				width:  d.border.left,
				height: borderBox.height,
			},
		},
	}

	// Right border
	list[len(list)] = DisplayCommand{
		command: SolidColor{
			color: *colorBorder,
			rect: Rect{
				x:      borderBox.x + borderBox.width - d.border.right,
				y:      borderBox.y,
				width:  d.border.right,
				height: borderBox.height,
			},
		},
	}

	// Top border
	list[len(list)] = DisplayCommand{
		command: SolidColor{
			color: *colorBorder,
			rect: Rect{
				x:      borderBox.x,
				y:      borderBox.y,
				width:  borderBox.width,
				height: d.border.top,
			},
		},
	}

	// Bottom border
	list[len(list)] = DisplayCommand{
		command: SolidColor{
			color: *colorBorder,
			rect: Rect{
				x:      borderBox.x,
				y:      borderBox.y + borderBox.height - d.border.bottom,
				width:  borderBox.width,
				height: d.border.bottom,
			},
		},
	}
}

//Return the specified color for CSS property name
func getColor(layoutBox *LayoutBox, name string) *Color {
	switch layoutBox.box_type.(type) {
	case BlockNode, InlineNode:
		color := layoutBox.style.value(name).color
		return &color
	case AnonymousBlock:
		return nil
	default:
		return nil
	}
}
