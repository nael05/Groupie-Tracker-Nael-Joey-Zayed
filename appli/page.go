package groupie

import (
	"image/color"
	"sort"
	"strconv"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/storage"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

var maFenetre fyne.Window
var tousLesArtistes []List_artist
var lApplication fyne.App

func OuvertureApp() {
	lApplication = app.New()
	lApplication.Settings().SetTheme(theme.LightTheme())

	maFenetre = lApplication.NewWindow("Groupie Tracker nael-joey-zayed")
	maFenetre.Resize(fyne.NewSize(1000, 800))
	maFenetre.CenterOnScreen()

	dictionnaireArtistes := Api()

	for _, artiste := range dictionnaireArtistes {
		tousLesArtistes = append(tousLesArtistes, artiste)
	}

	sort.Slice(tousLesArtistes, func(i, j int) bool {
		return tousLesArtistes[i].Name < tousLesArtistes[j].Name
	})

	AfficherLeMenu()

	maFenetre.ShowAndRun()
}

func AfficherLeMenu() {
	grilleArtistes := container.NewGridWithColumns(5)

	for _, artisteCourant := range tousLesArtistes {
		var imageDeBase *canvas.Image
		lienImage, erreur := storage.ParseURI(artisteCourant.Image)
		if erreur == nil {
			imageDeBase = canvas.NewImageFromURI(lienImage)
		} else {
			imageDeBase = canvas.NewImageFromResource(theme.MediaMusicIcon())
		}
		imageDeBase.FillMode = canvas.ImageFillContain
		imageDeBase.SetMinSize(fyne.NewSize(150, 150))

		artistePourLeClic := artisteCourant
		imageInteractive := NewImageCliquable(imageDeBase, func() {
			AfficherLesDetails(artistePourLeClic)
		})

		etiquetteNom := widget.NewLabel(artistePourLeClic.Name)
		etiquetteNom.Alignment = fyne.TextAlignCenter
		etiquetteNom.TextStyle = fyne.TextStyle{Bold: true}

		carteArtiste := container.NewVBox(imageInteractive, etiquetteNom)
		grilleArtistes.Add(carteArtiste)
	}

	conteneurDefilant := container.NewVScroll(grilleArtistes)

	titrePrincipal := canvas.NewText("Nos Artistes", color.Black)
	titrePrincipal.TextSize = 60
	titrePrincipal.TextStyle = fyne.TextStyle{Bold: true}
	titrePrincipal.Alignment = fyne.TextAlignCenter

	contenuFinal := container.NewBorder(titrePrincipal, nil, nil, nil, conteneurDefilant)
	maFenetre.SetContent(contenuFinal)
}

func AfficherLesDetails(artiste List_artist) {
	boutonHome := widget.NewButtonWithIcon("Retour", theme.HomeIcon(), func() {
		AfficherLeMenu()
	})

	barreNavigation := container.NewHBox(layout.NewSpacer(), boutonHome)

	var grandeImage *canvas.Image
	uri, err := storage.ParseURI(artiste.Image)
	if err == nil {
		grandeImage = canvas.NewImageFromURI(uri)
	} else {
		grandeImage = canvas.NewImageFromResource(theme.MediaMusicIcon())
	}
	grandeImage.FillMode = canvas.ImageFillContain
	grandeImage.SetMinSize(fyne.NewSize(350, 350))

	lblNom := canvas.NewText(artiste.Name, color.Black)
	lblNom.TextSize = 30
	lblNom.TextStyle = fyne.TextStyle{Bold: true}

	texteDate := "Date de création : " + strconv.Itoa(artiste.CreationDate)
	lblDate := widget.NewLabelWithStyle(texteDate, fyne.TextAlignLeading, fyne.TextStyle{Italic: true})

	lblAlbum := widget.NewLabel("Premier Album : " + artiste.FirstAlbum)
	lblAlbum.TextStyle = fyne.TextStyle{Bold: true}

	listeMembres := strings.Join(artiste.Members, "\n- ")
	lblMembresTitre := widget.NewLabelWithStyle("Membres :", fyne.TextAlignLeading, fyne.TextStyle{Bold: true, Underline: true})
	lblMembresListe := widget.NewLabel("- " + listeMembres)

	lblRelationsTitre := widget.NewLabelWithStyle("Concerts (Lieux & Dates) :", fyne.TextAlignLeading, fyne.TextStyle{Bold: true, Underline: true})

	conteneurRelations := container.NewVBox()
	
	relations, errRel := GetRelations(artiste.RelationsUrl)
	if errRel != nil {
		conteneurRelations.Add(widget.NewLabel("Impossible de charger les infos concerts."))
	} else {
		for lieu, dates := range relations {
			lieuPropre := strings.ReplaceAll(lieu, "_", " ")
			lieuPropre = strings.Title(lieuPropre)
			
			lblLieu := canvas.NewText(lieuPropre, color.NRGBA{R: 0, G: 0, B: 150, A: 255})
			lblLieu.TextStyle = fyne.TextStyle{Bold: true}
			
			datesStr := strings.Join(dates, ", ")
			lblDates := widget.NewLabel(datesStr)
			lblDates.Wrapping = fyne.TextWrapWord

			conteneurRelations.Add(lblLieu)
			conteneurRelations.Add(lblDates)
			conteneurRelations.Add(layout.NewSpacer())
		}
	}

	infosContainer := container.NewVBox(
		lblNom,
		widget.NewSeparator(),
		lblDate,
		lblAlbum,
		widget.NewSeparator(),
		lblMembresTitre,
		lblMembresListe,
		widget.NewSeparator(),
		lblRelationsTitre,
		conteneurRelations,
	)

	split := container.NewHSplit(grandeImage, container.NewVScroll(infosContainer))
	split.SetOffset(0.4)

	pageDetails := container.NewBorder(barreNavigation, nil, nil, nil, split)

	maFenetre.SetContent(pageDetails)
}

type ImageCliquable struct {
	widget.BaseWidget
	Image *canvas.Image
	OnTap func()
}

func NewImageCliquable(img *canvas.Image, onTap func()) *ImageCliquable {
	imageCliquable := &ImageCliquable{
		Image: img,
		OnTap: onTap,
	}
	imageCliquable.ExtendBaseWidget(imageCliquable)
	return imageCliquable
}

func (i *ImageCliquable) Tapped(_ *fyne.PointEvent) {
	if i.OnTap != nil {
		i.OnTap()
	}
}

func (i *ImageCliquable) CreateRenderer() fyne.WidgetRenderer {
	return widget.NewSimpleRenderer(i.Image)
}