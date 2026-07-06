package main

import (
	"fmt"
	"image/color"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

// ---- Duration formatting ------------------------------------------------

func formatClock(d time.Duration) string {
	if d < 0 {
		d = 0
	}
	h := d / time.Hour
	d -= h * time.Hour
	m := d / time.Minute
	d -= m * time.Minute
	s := d / time.Second
	d -= s * time.Second
	cs := d / (10 * time.Millisecond)
	return fmt.Sprintf("%02d:%02d:%02d.%02d", h, m, s, cs)
}

// formatSecondsShort mirrors the source project's formatShortDuration:
// above a minute it's "Xm SSs", below a minute it's "S.CCs" (seconds not
// zero-padded in that branch, matching the original exactly).
func formatSecondsShort(d time.Duration) string {
	if d < 0 {
		d = 0
	}
	totalMs := d.Milliseconds()
	minutes := totalMs / 60000
	seconds := (totalMs % 60000) / 1000
	centis := (totalMs % 1000) / 10
	if minutes > 0 {
		return fmt.Sprintf("%dm %02ds", minutes, seconds)
	}
	return fmt.Sprintf("%d.%02ds", seconds, centis)
}

// ---- Responsive split layout --------------------------------------------

// responsiveSplitLayout arranges [logPanel, divider, timerPanel]. In
// landscape (width >= height) the log panel sits on the left, timer panel
// on the right, divided by a vertical line. In portrait, the timer panel
// sits on top and the log panel on the bottom, divided by a horizontal line.
type responsiveSplitLayout struct{}

const dividerThickness float32 = 2

func (r *responsiveSplitLayout) Layout(objs []fyne.CanvasObject, size fyne.Size) {
	logPanel, divider, timerPanel := objs[0], objs[1], objs[2]

	if size.Width >= size.Height {
		half := size.Width / 2
		logPanel.Move(fyne.NewPos(0, 0))
		logPanel.Resize(fyne.NewSize(half-dividerThickness/2, size.Height))

		divider.Move(fyne.NewPos(half-dividerThickness/2, 0))
		divider.Resize(fyne.NewSize(dividerThickness, size.Height))

		timerPanel.Move(fyne.NewPos(half+dividerThickness/2, 0))
		timerPanel.Resize(fyne.NewSize(size.Width-half-dividerThickness/2, size.Height))
	} else {
		half := size.Height / 2
		timerPanel.Move(fyne.NewPos(0, 0))
		timerPanel.Resize(fyne.NewSize(size.Width, half-dividerThickness/2))

		divider.Move(fyne.NewPos(0, half-dividerThickness/2))
		divider.Resize(fyne.NewSize(size.Width, dividerThickness))

		logPanel.Move(fyne.NewPos(0, half+dividerThickness/2))
		logPanel.Resize(fyne.NewSize(size.Width, size.Height-half-dividerThickness/2))
	}
}

func (r *responsiveSplitLayout) MinSize(objs []fyne.CanvasObject) fyne.Size {
	logPanel, _, timerPanel := objs[0], objs[1], objs[2]
	lm, tm := logPanel.MinSize(), timerPanel.MinSize()
	w := lm.Width + tm.Width
	h := lm.Height
	if tm.Height > h {
		h = tm.Height
	}
	return fyne.NewSize(w, h)
}

// fixedSizeLayout forces its single child to a fixed size regardless of
// content - used for the small square S{n} index badge.
type fixedSizeLayout struct {
	size fyne.Size
}

func (f fixedSizeLayout) Layout(objs []fyne.CanvasObject, _ fyne.Size) {
	for _, o := range objs {
		o.Resize(f.size)
	}
}

func (f fixedSizeLayout) MinSize([]fyne.CanvasObject) fyne.Size {
	return f.size
}

// ---- Entry row -----------------------------------------------------------

type entryRow struct {
	box          *fyne.Container
	indexBadgeBG *canvas.Rectangle
	indexBadge   *canvas.Text
	titleText    *canvas.Text
	badgeText    *canvas.Text
	preSetLabel  *widget.Label
	slackText    *canvas.Text
	pausedLabel  *widget.Label
	border       *canvas.Rectangle
}

func newEntryRow() *entryRow {
	indexBadgeBG := canvas.NewRectangle(colorIndigoStrong)
	indexBadgeBG.CornerRadius = 6
	indexBadge := canvas.NewText("", colorIndigo)
	indexBadge.TextStyle = fyne.TextStyle{Bold: true}
	indexBadge.TextSize = 12
	indexBadge.Alignment = fyne.TextAlignCenter
	badgeStack := container.NewStack(indexBadgeBG, container.NewCenter(indexBadge))
	badgeSized := container.New(fixedSizeLayout{size: fyne.NewSize(32, 32)}, badgeStack)

	title := canvas.NewText("", colorSlate200)
	title.TextStyle = fyne.TextStyle{Bold: true}
	title.TextSize = 12

	badge := canvas.NewText("", colorGreen)
	badge.TextStyle = fyne.TextStyle{Bold: true}
	badge.TextSize = 10
	badge.Hidden = true

	preSet := widget.NewLabel("")
	preSet.Importance = widget.LowImportance

	slack := canvas.NewText("", colorSlate200)
	slack.TextStyle = fyne.TextStyle{Monospace: true, Bold: true}
	slack.TextSize = 18

	paused := widget.NewLabel("")
	paused.Importance = widget.LowImportance
	paused.Alignment = fyne.TextAlignTrailing

	titleRow := container.NewHBox(title, badge)
	textCol := container.NewVBox(titleRow, preSet)
	leftSide := container.NewBorder(nil, nil, badgeSized, nil, textCol)

	bottomRow := container.NewBorder(nil, nil, slack, paused)
	content := container.NewVBox(leftSide, bottomRow)

	border := canvas.NewRectangle(color.Transparent)
	border.StrokeColor = colorCardBorder
	border.StrokeWidth = 1
	border.FillColor = colorCardBG
	border.CornerRadius = 12

	padded := container.NewPadded(content)
	stacked := container.NewStack(border, padded)

	return &entryRow{
		box:          stacked,
		indexBadgeBG: indexBadgeBG,
		indexBadge:   indexBadge,
		titleText:    title,
		badgeText:    badge,
		preSetLabel:  preSet,
		slackText:    slack,
		pausedLabel:  paused,
		border:       border,
	}
}

func (r *entryRow) update(e EntrySnapshot, totalAtPause time.Duration) {
	r.titleText.Text = fmt.Sprintf("SET %d + SLACK TRACK", e.Index)
	r.titleText.Refresh()

	r.indexBadge.Text = fmt.Sprintf("S%d", e.Index)
	r.preSetLabel.SetText(fmt.Sprintf("Pre-Set Rest Duration: %s", formatSecondsShort(e.PreSetRest)))
	r.pausedLabel.SetText(fmt.Sprintf("Paused at %s total", formatSecondsShort(totalAtPause)))
	r.slackText.Text = formatClock(e.Slack)

	if e.InProgress {
		r.slackText.Color = colorGreen
		r.badgeText.Text = "Set In Progress"
		r.badgeText.Hidden = false
		r.border.StrokeColor = colorCardBorderActive
		r.border.FillColor = colorCardBGActive
		r.indexBadge.Color = colorGreen
		r.indexBadgeBG.FillColor = color.NRGBA{R: 0x34, G: 0xD3, B: 0x99, A: 0x1A} // emerald-400/10
	} else {
		r.slackText.Color = colorSlate200
		r.badgeText.Hidden = true
		r.border.StrokeColor = colorCardBorder
		r.border.FillColor = colorCardBG
		r.indexBadge.Color = colorIndigo
		r.indexBadgeBG.FillColor = color.NRGBA{R: 0x63, G: 0x66, B: 0xF1, A: 0x1A} // indigo-500/10
	}
	r.indexBadge.Refresh()
	r.indexBadgeBG.Refresh()
	r.badgeText.Refresh()
	r.slackText.Refresh()
	r.border.Refresh()
}

// ---- Main UI ---------------------------------------------------------------

type UI struct {
	state *AppState
	win   fyne.Window

	statusText  *canvas.Text
	statusBadge *canvas.Rectangle
	totalMain   *canvas.Text // hh:mm:ss part - always light, never recolored
	totalFrac   *canvas.Text // .cs part - colored by state
	toggleBtn   *widget.Button
	resetBtn    *widget.Button

	avgLabel   *widget.Label
	totalLabel *widget.Label

	emptyState *fyne.Container
	listBox    *fyne.Container
	rows       []*entryRow
}

func NewUI(win fyne.Window, state *AppState) *UI {
	u := &UI{state: state, win: win}

	// ---- Timer panel (right in landscape / top in portrait) ----
	u.statusText = canvas.NewText("READY TO REST", colorMuted)
	u.statusText.Alignment = fyne.TextAlignCenter
	u.statusText.TextStyle = fyne.TextStyle{Bold: true}
	u.statusText.TextSize = 13

	u.statusBadge = canvas.NewRectangle(color.NRGBA{R: 0x1A, G: 0x1F, B: 0x30, A: 0xFF})
	u.statusBadge.CornerRadius = 14
	u.statusBadge.StrokeColor = colorCardBorder
	u.statusBadge.StrokeWidth = 1
	statusStack := container.NewStack(u.statusBadge, container.NewPadded(container.NewCenter(u.statusText)))

	u.totalMain = canvas.NewText("00:00:00", colorWhite)
	u.totalMain.TextStyle = fyne.TextStyle{Monospace: true, Bold: true}
	u.totalMain.TextSize = 64

	u.totalFrac = canvas.NewText(".00", colorMuted)
	u.totalFrac.TextStyle = fyne.TextStyle{Monospace: true, Bold: true}
	u.totalFrac.TextSize = 40

	timeRow := container.NewHBox(u.totalMain, u.totalFrac)

	caption := canvas.NewText("TOTAL WORKOUT ACTION TIME", colorMuted)
	caption.Alignment = fyne.TextAlignCenter
	caption.TextSize = 12

	u.resetBtn = widget.NewButton("Reset", func() {
		u.state.Reset()
		u.rebuildList()
		u.refresh()
	})
	u.toggleBtn = widget.NewButton("Start", func() {
		u.state.Toggle(time.Now())
		u.rebuildList()
		u.refresh()
	})
	hint := canvas.NewText("PRESS SPACE TO TOGGLE  •  PRESS R TO RESET", colorMuted)
	hint.Alignment = fyne.TextAlignCenter
	hint.TextSize = 11

	buttonsRow := container.NewCenter(container.NewHBox(u.resetBtn, u.toggleBtn))

	timerPanel := container.NewCenter(container.NewVBox(
		container.NewCenter(statusStack),
		container.NewPadded(container.NewCenter(timeRow)),
		container.NewCenter(caption),
		container.NewPadded(buttonsRow),
		container.NewCenter(hint),
	))

	// ---- Log panel (left in landscape / bottom in portrait) ----
	title := canvas.NewText("EXCESS REST & SET LOG", colorWhite)
	title.TextStyle = fyne.TextStyle{Bold: true}
	title.TextSize = 15

	u.avgLabel = widget.NewLabel("Avg. Slack Time\n0.00s")
	u.avgLabel.Alignment = fyne.TextAlignTrailing
	u.totalLabel = widget.NewLabel("Total Slack\n0.00s")
	u.totalLabel.Alignment = fyne.TextAlignTrailing

	header := container.NewBorder(nil, nil, title, container.NewHBox(u.avgLabel, u.totalLabel))

	u.emptyState = container.NewCenter(container.NewVBox(
		widget.NewLabelWithStyle("No workout sets started yet", fyne.TextAlignCenter, fyne.TextStyle{Bold: true}),
		widget.NewLabelWithStyle(
			"Start the stopwatch at 0:00. This is your rest time.\nWhen you're ready to lift, press Pause to start your set.\nThe timer here tracks your rest and any excess slack.",
			fyne.TextAlignCenter, fyne.TextStyle{},
		),
	))

	u.listBox = container.NewVBox()
	scroll := container.NewVScroll(u.listBox)

	logPanel := container.NewBorder(container.NewPadded(header), nil, nil, nil, container.NewPadded(scroll))

	// ---- Header bar ----
	appTitle := canvas.NewText("Workout Pause Timer", colorWhite)
	appTitle.TextStyle = fyne.TextStyle{Bold: true}
	appTitle.TextSize = 16
	appSubtitle := canvas.NewText("Track rest periods between sets in real-time", colorSlate400)
	appSubtitle.TextSize = 12
	headerText := container.NewVBox(appTitle, appSubtitle)

	headerBG := canvas.NewRectangle(colorHeaderBG)
	header2 := container.NewStack(headerBG, container.NewPadded(headerText))

	// ---- Compose responsive layout ----
	divider := canvas.NewRectangle(colorDivider)
	split := container.New(&responsiveSplitLayout{}, logPanel, divider, timerPanel)
	root := container.NewBorder(header2, nil, nil, nil, split)

	win.SetContent(root)

	// Keyboard shortcuts
	win.Canvas().SetOnTypedKey(func(ev *fyne.KeyEvent) {
		switch ev.Name {
		case fyne.KeySpace:
			u.state.Toggle(time.Now())
			u.rebuildList()
			u.refresh()
		case fyne.KeyR:
			u.state.Reset()
			u.rebuildList()
			u.refresh()
		}
	})

	u.rebuildList()
	u.refresh()

	return u
}

// rebuildList rebuilds the log panel's row widgets to match the current
// number of entries. Called after any Start/Pause/Reset action.
func (u *UI) rebuildList() {
	count := u.state.EntryCount()
	if count == 0 {
		u.listBox.Objects = []fyne.CanvasObject{u.emptyState}
		u.rows = nil
		u.listBox.Refresh()
		return
	}
	if len(u.rows) != count {
		u.rows = make([]*entryRow, count)
		objs := make([]fyne.CanvasObject, count)
		for i := 0; i < count; i++ {
			u.rows[i] = newEntryRow()
			objs[i] = u.rows[i].box
		}
		u.listBox.Objects = objs
		u.listBox.Refresh()
	}
}

// refresh updates all live-changing text: the main timer, status, button
// label, the in-progress log row (if any), and the summary stats. Safe to
// call frequently (e.g. from a ticker) and after any state-changing action.
func (u *UI) refresh() {
	now := time.Now()

	total := u.state.TotalElapsed(now)
	full := formatClock(total)
	u.totalMain.Text = full[:8] // hh:mm:ss
	u.totalFrac.Text = full[8:] // .cs
	running := u.state.IsRunning()
	if running {
		u.totalFrac.Color = colorGreen
	} else if u.state.HasStarted() {
		u.totalFrac.Color = colorOrange
	} else {
		u.totalFrac.Color = colorMuted
	}
	u.totalMain.Refresh()
	u.totalFrac.Refresh()

	setNum := u.state.CurrentSetNumber()
	switch {
	case running:
		// "Resting / preparing for set N" - indigo, matching the source's
		// running-state badge.
		u.statusText.Text = fmt.Sprintf("RESTING / PREPARING FOR SET %d", setNum)
		u.statusText.Color = colorIndigo
		u.statusBadge.FillColor = color.NRGBA{R: 0x63, G: 0x66, B: 0xF1, A: 0x1A}   // indigo-500/10
		u.statusBadge.StrokeColor = color.NRGBA{R: 0x63, G: 0x66, B: 0xF1, A: 0x4D} // indigo-500/30
	case u.state.HasStarted():
		// "Active set N" (currently paused / lifting) - emerald.
		u.statusText.Text = fmt.Sprintf("ACTIVE SET: %d", setNum)
		u.statusText.Color = colorGreen
		u.statusBadge.FillColor = color.NRGBA{R: 0x10, G: 0xB9, B: 0x81, A: 0x1A}   // emerald-500/10
		u.statusBadge.StrokeColor = color.NRGBA{R: 0x10, G: 0xB9, B: 0x81, A: 0x4D} // emerald-500/30
	default:
		u.statusText.Text = "READY TO LIFT"
		u.statusText.Color = colorMuted
		u.statusBadge.FillColor = color.NRGBA{R: 0x1E, G: 0x29, B: 0x3B, A: 0x4D} // slate-800/30
		u.statusBadge.StrokeColor = colorSlate800
	}
	u.statusText.Refresh()
	u.statusBadge.Refresh()

	switch {
	case running:
		u.toggleBtn.SetText("Pause")
		u.toggleBtn.Importance = widget.WarningImportance
	case u.state.HasStarted():
		u.toggleBtn.SetText("Resume")
		u.toggleBtn.Importance = widget.SuccessImportance
	default:
		u.toggleBtn.SetText("Start")
		u.toggleBtn.Importance = widget.SuccessImportance
	}
	u.toggleBtn.Refresh()

	// Live-update the in-progress row, if any.
	if n := u.state.EntryCount(); n > 0 && len(u.rows) == n {
		snap := u.state.Entry(n-1, now)
		if snap.InProgress {
			// "Paused at" = total elapsed at the moment this pause happened,
			// i.e. accumulated total minus nothing extra (segment already
			// folded in) - equivalently, the running total right now since
			// it's frozen while paused.
			u.rows[n-1].update(snap, total)
		}
	}
	// Backfill totals for locked rows once (cheap to just do every refresh
	// for correctness given the small row counts involved).
	for i := 0; i < len(u.rows); i++ {
		snap := u.state.Entry(i, now)
		if !snap.InProgress {
			// total-at-pause for a locked entry = sum of PreSetRest up to
			// and including this entry.
			var cum time.Duration
			for j := 0; j <= i; j++ {
				s := u.state.Entry(j, now)
				cum += s.PreSetRest
			}
			u.rows[i].update(snap, cum)
		}
	}

	avg, tot := u.state.Stats(now)
	u.avgLabel.SetText(fmt.Sprintf("Avg. Slack Time\n%s", formatSecondsShort(avg)))
	u.totalLabel.SetText(fmt.Sprintf("Total Slack\n%s", formatSecondsShort(tot)))
}

// StartTicker runs a background refresh loop for live-updating displays
// (the running timer and the in-progress slack entry).
func (u *UI) StartTicker() {
	ticker := time.NewTicker(30 * time.Millisecond)
	go func() {
		for range ticker.C {
			u.refresh()
		}
	}()
}
