package newbutton

import (
	"image"
	"image/color"
	"time"

	"gioui.org/gesture"
	"gioui.org/io/event"
	"gioui.org/layout"
	"gioui.org/op/clip"
	"gioui.org/op/paint"
	"gioui.org/unit"
	"gioui.org/widget/material"
)

type (
	C = layout.Context
	D = layout.Dimensions
)

// --------------------------------------------------------------------------------
type CircledButton struct {
	click      gesture.Click
	hover      gesture.Hover
	Label      string
	TxtColor   color.NRGBA
	BkColor    color.NRGBA
	HoverColor color.NRGBA
	PressColor color.NRGBA
	ClickColor color.NRGBA
	Size       unit.Dp
	ClickedAt  time.Time
	NumClicks  int
	// Cached state variables
	IsClicked   bool
	IsPressed   bool
	IsHovered   bool
	IsCancelled bool
}

func (cb *CircledButton) update(gtx C) {
	cb.IsClicked = false
	cb.IsCancelled = false
	// Treat the clicks here
	for {
		ev, ok := cb.click.Update(gtx.Source)
		if !ok {
			break
		}
		switch ev.Kind {
		case gesture.KindClick:
			cb.IsClicked = true
			cb.NumClicks = ev.NumClicks
		case gesture.KindCancel:
			cb.IsCancelled = true
		}
	}

	// Treat Press and Hover here
	cb.IsPressed = cb.click.Pressed()
	cb.IsHovered = cb.hover.Update(gtx.Source)
}

func (cb *CircledButton) Layout(gtx C, th *material.Theme) D {
	cb.update(gtx)

	sz := image.Pt(gtx.Dp(cb.Size), gtx.Dp(cb.Size))
	gtx.Constraints.Min = sz
	gtx.Constraints.Max = sz

	// Change the color when the state changes
	bk := cb.BkColor
	if cb.IsPressed {
		cb.BkColor = cb.PressColor
	} else if cb.IsHovered {
		cb.BkColor = cb.HoverColor
	} else if cb.IsClicked {
		cb.BkColor = cb.ClickColor
	}

	// Draw the circle and make it clickable
	layout.Stack{Alignment: layout.Center}.Layout(gtx,
		layout.Stacked(func(gtx C) D {
			area := clip.Ellipse{Max: sz}.Push(gtx.Ops)
			event.Op(gtx.Ops, &cb.click)
			paint.Fill(gtx.Ops, bk)
			cb.click.Add(gtx.Ops)
			cb.hover.Add(gtx.Ops)
			area.Pop()
			return D{Size: sz}
		}),
		layout.Expanded(func(gtx C) D {
			return layout.Center.Layout(gtx, func(gtx C) D {
				return material.Label(th, 30, cb.Label).Layout(gtx)
			})
		}),
	)
	return D{Size: sz}
}

func (cb *CircledButton) Clicked(gtx C) bool {
	cb.update(gtx)
	return cb.IsClicked
}

func (cb *CircledButton) Pressed() bool {
	return cb.IsPressed
}

func (cb *CircledButton) Hovered() bool {
	return cb.IsHovered
}

func (cb *CircledButton) Cancelled() bool {
	return cb.IsCancelled
}

// --------------------------------------------------------------------------------
