package main

import (
	_ "embed"

	"fyne.io/fyne/v2"
)

//go:embed assets/fonts/JetBrainsMono-Regular.ttf
var jetBrainsMonoRegularBytes []byte

//go:embed assets/fonts/JetBrainsMono-Bold.ttf
var jetBrainsMonoBoldBytes []byte

//go:embed assets/fonts/Inter-Regular.ttf
var interRegularBytes []byte

//go:embed assets/fonts/Inter-Bold.ttf
var interBoldBytes []byte

var (
	fontMonoRegular = fyne.NewStaticResource("JetBrainsMono-Regular.ttf", jetBrainsMonoRegularBytes)
	fontMonoBold    = fyne.NewStaticResource("JetBrainsMono-Bold.ttf", jetBrainsMonoBoldBytes)
	fontSansRegular = fyne.NewStaticResource("Inter-Regular.ttf", interRegularBytes)
	fontSansBold    = fyne.NewStaticResource("Inter-Bold.ttf", interBoldBytes)
)
