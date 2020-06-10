package main
/*
This file is part of Alfred
(c) 2020 - 0xSha.io
*/

import (
	"InfoSecAlfred/scrap"
	"fyne.io/fyne"
	"fyne.io/fyne/app"
	"fyne.io/fyne/canvas"
	"fyne.io/fyne/layout"
	"fyne.io/fyne/widget"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	"image/color"
	"log"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"time"
)

func main() {

	db, err := scrap.InitDB()
	defer db.Close()

	var master scrap.Master
	var entity scrap.Entity

	if err := db.First(&master, 1).Error; err != nil {
		log.Println("Master Name not found.")

		// i'm the default master :>
		master.Name = "0XSha"
		db.Create(&master)
	} else {
		db.Table("masters").Select("name").Row().Scan(&master.Name)
		log.Println(master.Name)
	}

	if err != nil {

		log.Println("Can't init DB exiting ...")
		os.Exit(-1)
	}

	myApp := app.New()
	myWindow := myApp.NewWindow("InfoSec Alfred")

	// fun
	greet := canvas.NewText("Hello master "+master.Name, color.RGBA{
		R: 189,
		G: 147,
		B: 249,
		A: 0,
	})
	centered := fyne.NewContainerWithLayout(layout.NewHBoxLayout(),
		layout.NewSpacer(), greet, layout.NewSpacer())

	image := canvas.NewImageFromResource(resourceAlfredLgPng) // NewImageFromFile( "./assets/alfred-lg.png")

	myWindow.Resize(fyne.NewSize(300, 300))

	image.FillMode = canvas.ImageFillOriginal

	progress := widget.NewProgressBar()
	progress.SetValue(0)

	status := widget.NewLabel("Idle...")

	statusContainer := fyne.NewContainerWithLayout(layout.NewHBoxLayout(),
		layout.NewSpacer(), status, layout.NewSpacer())

	exitButton := widget.NewButton("Exit", func() {
		myApp.Quit()
	})

	aboutButton := widget.NewButton("About", func() {
		ShowAbout(myApp)
	})

	reloadButton := widget.NewButton("Update", func() {

		_, err := http.Get("http://google.com")
		if err != nil {
			ShowNotification(myApp, "Sorry Master "+master.Name+" but I need internet connection to do that.", "Sorry .")
		}

		total, err := FetchEverything(db, entity, progress, status)
		totalStr := strconv.Itoa(total)
		log.Println("Total new links : " + totalStr)

		if err != nil {

			ShowNotification(myApp, "Sorry Master"+master.Name+" An Error occurred while fetching data", "Sorry")

		} else {

			if total == 0 {

				ShowNotification(myApp, "Sorry Master "+master.Name+" I've fetched everything but no new links found.", "Sorry")

			} else {
				ShowNotification(myApp, "Alright Master "+master.Name+" I've added :"+totalStr+" links to my library", "Links added")
			}
		}

	})

	imageCentered := fyne.NewContainerWithLayout(layout.NewHBoxLayout(), layout.NewSpacer(), image, layout.NewSpacer())

	buttons := fyne.NewContainerWithLayout(layout.NewHBoxLayout(), layout.NewSpacer(), reloadButton, aboutButton, exitButton, layout.NewSpacer())

	myWindow.SetContent(fyne.NewContainerWithLayout(layout.NewVBoxLayout(), centered, imageCentered, progress, statusContainer, buttons))

	myWindow.ShowAndRun()

}

func ShowAbout(a fyne.App) {

	win := a.NewWindow("About")
	win.Resize(fyne.NewSize(200, 100))
	win.SetFixedSize(true)

	aboutURL := url.URL{
		Scheme: "https",
		Host:   "0xsha.io",
	}

	copyText := widget.NewLabel("CopyRight (c) 2020")

	about := widget.NewHyperlink("By 0XSha.io", &aboutURL)

	container := fyne.NewContainerWithLayout(layout.NewHBoxLayout(),
		layout.NewSpacer(), about, layout.NewSpacer())

	containerC := fyne.NewContainerWithLayout(layout.NewHBoxLayout(),
		layout.NewSpacer(), copyText, layout.NewSpacer())

	win.SetContent(fyne.NewContainerWithLayout(layout.NewVBoxLayout(), container, containerC))

	win.Show()

}

func ShowNotification(a fyne.App, text string, title string) {
	//time.Sleep(time.Second * 5)

	//linkUrl := widget.NewHyperlink("Read",targetUrl)
	//titleLabel := widget.NewHyperlink(title , targetUrl)
	titleLabel := widget.NewLabel(text)

	image := canvas.NewImageFromResource(resourceAlfredSmPng) // NewImageFromFile( "./assets/alfred.png")
	image.FillMode = canvas.ImageFillOriginal

	container := fyne.NewContainerWithLayout(layout.NewHBoxLayout(),
		image, titleLabel, layout.NewSpacer())

	win := a.NewWindow(title)

	win.Resize(fyne.NewSize(300, 100))

	win.SetFixedSize(true)

	win.SetContent(fyne.NewContainerWithLayout(layout.NewVBoxLayout(), container))

	win.Show()

	time.Sleep(time.Second * 7)
	win.Hide()
}

// pass in DB instance, Entity, progress and label widget
func FetchEverything(db *gorm.DB, entity scrap.Entity, bar *widget.ProgressBar, label *widget.Label) (int, error) {

	totalFound := 0

	// later on need a proper calculation based on modules.
	var step = 0.0125

	//NewsAPI
	db.Raw("SELECT id FROM entities ORDER BY ID DESC LIMIT 1").Scan(&entity)
	newsAPI, err := scrap.FetchNewsAPI()
	label.SetText("Fetching NewsAPI ...")

	if err != nil {
		log.Println("Error fetching NewsAPI Data:", err)
		return 0, err
	}
	newsAPITotal, err := scrap.WriteNewsAPIToDB(newsAPI, entity, db)
	totalFound += newsAPITotal
	bar.SetValue(step)
	step += step

	//Reddit/NetSec
	redditNetSec, err := scrap.FetchNetSecReddit()
	label.SetText("Fetching Reddit NetSec ...")

	if err != nil {
		log.Println("Error fetching NewsAPI Data:", err)
		return 0, err
	}
	redditTotal, err := scrap.WriteRedditNetSecToDB(redditNetSec, entity, db)
	log.Println(redditTotal)

	bar.SetValue(step)
	step += step

	// Exploit-DB
	exploitDB, err := scrap.FetchExploitDB()
	label.SetText("Fetching Exploit-DB ...")

	if err != nil {
		log.Println("Error fetching Exploit-DB Data:", err)
		return 0, err
	}
	expTotal, err := scrap.WriteExploitDBToDB(exploitDB, entity, db)
	log.Println(expTotal)

	bar.SetValue(step)
	step += step

	// InfoSecJobs
	db.Raw("SELECT id FROM entities ORDER BY ID DESC LIMIT 1").Scan(&entity)
	entity.ID++

	infoSecJobs, err := scrap.FetchInfoSecJobs()
	label.SetText("InfoSec-Jobs...")

	if err != nil {
		log.Println("Error fetching InfoSecJobs Data:", err)
		return 0, err
	}

	infoSecJobsTotal, err := scrap.WriteInfoSecJobsToDB(infoSecJobs, entity, db)
	totalFound += infoSecJobsTotal

	bar.SetValue(step)
	step += step

	// Pentesterland
	db.Raw("SELECT id FROM entities ORDER BY ID DESC LIMIT 1").Scan(&entity)
	entity.ID++

	pentesterLand, err := scrap.FetchPentesterLand()
	label.SetText("Fetching PentesterLand...")

	if err != nil {
		log.Println("Error fetching Pentester,Land Data:", err)
		return 0, err
	}

	pentesterLandTotal, err := scrap.WritePentesterLandJobsToDB(pentesterLand, entity, db)
	log.Println(pentesterLandTotal)

	bar.SetValue(step)
	step += step

	// GithubAdvisory
	db.Raw("SELECT id FROM entities ORDER BY ID DESC LIMIT 1").Scan(&entity)
	entity.ID++

	githubAdvisory, err := scrap.FetchGithubAdvisory()
	label.SetText("Fetching Github Advisories..")

	if err != nil {
		log.Println("Error fetching Github Advisory Data:", err)
		return 0, err
	}

	githubAdvisoryTotal, err := scrap.WriteGithubAdvisoryToDB(githubAdvisory, entity, db)
	totalFound += githubAdvisoryTotal

	bar.SetValue(step)
	step += step

	//hackerOne
	db.Raw("SELECT id FROM entities ORDER BY ID DESC LIMIT 1").Scan(&entity)
	label.SetText("Fetching HackerOne ...")
	entity.ID++
	h1, err := scrap.FetchHackerOne()
	if err != nil {
		log.Println("Error fetching H1 Data:", err)
		return 0, err
	}
	h1Total, err := scrap.WriteH1ToDB(h1, entity, db)
	totalFound += h1Total

	bar.SetValue(step)
	step += step

	// TheHackerNews
	db.Raw("SELECT id FROM entities ORDER BY ID DESC LIMIT 1").Scan(&entity)
	label.SetText("Fetching HackerNews ...")

	entity.ID++
	hNews, err := scrap.FetchTheHackerNews()
	if err != nil {
		log.Println("Error fetching hackernews Data:", err)
		return 0, err

	}
	hNewsTotal, err := scrap.WriteHNewsToDB(hNews, entity, db)
	totalFound += hNewsTotal

	bar.SetValue(step)
	step += step

	label.SetText("Idle ...")
	return totalFound, nil

}
