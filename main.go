package main

import (
	_ "embed"
	"image/color"
	"strconv"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"
)

//go:embed icon.png
var iconData []byte

// Custom navy blue theme
type navyTheme struct{}

func (t navyTheme) Color(name fyne.ThemeColorName, variant fyne.ThemeVariant) color.Color {
	switch name {
	case theme.ColorNameBackground:
		return color.NRGBA{15, 15, 40, 255} // Dark navy
	case theme.ColorNameButton:
		return color.NRGBA{0, 0, 128, 255} // Navy blue
	case theme.ColorNameForeground:
		return color.NRGBA{220, 220, 255, 255} // Light blue-white
	case theme.ColorNamePrimary:
		return color.NRGBA{100, 149, 237, 255} // Cornflower blue
	default:
		return theme.DefaultTheme().Color(name, variant)
	}
}

func (t navyTheme) Font(style fyne.TextStyle) fyne.Resource {
	return theme.DefaultTheme().Font(style)
}

func (t navyTheme) Icon(name fyne.ThemeIconName) fyne.Resource {
	return theme.DefaultTheme().Icon(name)
}

func (t navyTheme) Size(name fyne.ThemeSizeName) float32 {
	return theme.DefaultTheme().Size(name)
}

// CountdownApp holds references to labels for updating
type CountdownApp struct {
	monthLabel *canvas.Text
	weekLabel  *canvas.Text
	dayLabel   *canvas.Text
	hourLabel  *canvas.Text
	ticker     *time.Ticker
}

// createStyledCard creates a custom card with rounded corners and specific colors
func createStyledCard(title, value string) (*canvas.Text, *fyne.Container) {
	// Create background rectangle with rounded corners
	background := canvas.NewRectangle(color.NRGBA{220, 220, 255, 255}) // Light blue-white
	background.SetMinSize(fyne.NewSize(80, 80))
	background.CornerRadius = 5

	// Create title text
	titleText := canvas.NewText(title, color.NRGBA{15, 15, 40, 255}) // Dark navy
	titleText.Alignment = fyne.TextAlignCenter
	titleText.TextStyle.Bold = true
	titleText.TextSize = 12

	// Create value text (this will be updated)
	valueText := canvas.NewText(value, color.NRGBA{15, 15, 40, 255}) // Dark navy
	valueText.Alignment = fyne.TextAlignCenter
	valueText.TextStyle.Bold = true
	valueText.TextSize = 32

	// Create text container with proper layout
	textContainer := container.New(layout.NewCenterLayout(), container.NewVBox(
		container.New(layout.NewCenterLayout(), titleText),
		container.New(layout.NewCenterLayout(), valueText),
	))

	// Stack background and text using border layout
	card := container.NewBorder(nil, nil, nil, nil, background, textContainer)
	card.Resize(fyne.NewSize(150, 100))

	return valueText, card
}

// calculateTimeRemaining computes all time components until end of year
func calculateTimeRemaining() (int, int, int, int) {
	now := time.Now()
	endOfYear := time.Date(now.Year(), 12, 31, 23, 59, 59, 0, now.Location())
	duration := endOfYear.Sub(now)

	// Calculate components
	totalHours := int(duration.Hours())
	totalDays := int(duration.Hours() / 24)

	monthLeft := 12 - int(now.Month())
	weekLeft := totalDays / 7
	dayLeft := totalDays
	hourLeft := totalHours

	return monthLeft, weekLeft, dayLeft, hourLeft
}

// updateCountdown refreshes all the labels with current values
func (app *CountdownApp) updateCountdown() {
	months, weeks, days, hours := calculateTimeRemaining()

	// Wrap UI updates in fyne.Do() for thread safety
	fyne.Do(func() {
		app.monthLabel.Text = strconv.Itoa(months)
		app.weekLabel.Text = strconv.Itoa(weeks)
		app.dayLabel.Text = strconv.Itoa(days)
		app.hourLabel.Text = strconv.Itoa(hours)

		// Refresh the canvas objects to show changes
		app.monthLabel.Refresh()
		app.weekLabel.Refresh()
		app.dayLabel.Refresh()
		app.hourLabel.Refresh()
	})
}

// startTicker begins live updates every minute
func (app *CountdownApp) startTicker() {
	app.ticker = time.NewTicker(1 * time.Minute)
	go func() {
		for range app.ticker.C {
			app.updateCountdown()
		}
	}()
}

// stopTicker properly cleans up the ticker
func (app *CountdownApp) stopTicker() {
	if app.ticker != nil {
		app.ticker.Stop()
	}
}

func main() {
	myApp := app.New()
	myApp.Settings().SetTheme(&navyTheme{})

	// Load icon from root directory
	iconResource := fyne.NewStaticResource("icon", iconData)
	myApp.SetIcon(iconResource)

	myWindow := myApp.NewWindow("Year of Synchronicity")
	myWindow.Resize(fyne.Size{Width: 360, Height: 110})
	myWindow.CenterOnScreen()

	// Create countdown app with custom styled cards
	countdownApp := &CountdownApp{}

	// Create custom styled cards and get references to value labels
	monthValue, monthCard := createStyledCard("Months", "0")
	weekValue, weekCard := createStyledCard("Weeks", "0")
	dayValue, dayCard := createStyledCard("Days", "0")
	hourValue, hourCard := createStyledCard("Hours", "0")

	// Store references for updating
	countdownApp.monthLabel = monthValue
	countdownApp.weekLabel = weekValue
	countdownApp.dayLabel = dayValue
	countdownApp.hourLabel = hourValue

	// Your original layout structure
	content := container.New(
		layout.NewCenterLayout(),
		container.NewBorder(
			nil,
			container.New(layout.NewGridLayout(4),
				monthCard,
				weekCard,
				dayCard,
				hourCard,
			),
			nil,
			nil,
			nil,
		),
	)

	myWindow.SetContent(content)

	// Initial update and start live updates
	countdownApp.updateCountdown()
	countdownApp.startTicker()

	// Cleanup when window closes
	myWindow.SetOnClosed(func() {
		countdownApp.stopTicker()
	})

	myWindow.ShowAndRun()
}
