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

var filtresActuels FilterOptions
var rechercheActuelle string

type FilterOptions struct {
	CreationMin   int
	CreationMax   int
	MembersCounts []int
}

func OuvertureApp() {
	lApplication = app.New()
	lApplication.Settings().SetTheme(theme.DarkTheme())

	maFenetre = lApplication.NewWindow("Groupie Tracker nael-joey-zayed")
	maFenetre.Resize(fyne.NewSize(1200, 800))
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
	entryRecherche := widget.NewEntry()
	entryRecherche.SetPlaceHolder("Rechercher un artiste, membre, année...")
	entryRecherche.SetText(rechercheActuelle)

	grilleArtistes := container.NewGridWithColumns(4)
	conteneurDefilant := container.NewVScroll(grilleArtistes)

	rafraichirGrille := func() {
		grilleArtistes.Objects = nil

		artistesAffiches := FiltrerArtistes(tousLesArtistes, filtresActuels)
		artistesAffiches = RechercherArtistes(artistesAffiches, rechercheActuelle)

		for _, artisteCourant := range artistesAffiches {
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
		grilleArtistes.Refresh()
	}

	entryRecherche.OnChanged = func(texte string) {
		rechercheActuelle = texte
		rafraichirGrille()
	}

	barreRecherche := container.NewBorder(nil, nil,
		widget.NewIcon(theme.SearchIcon()),
		nil,
		entryRecherche,
	)

	entryCreationMin := widget.NewEntry()
	entryCreationMin.SetPlaceHolder("Année min")
	entryCreationMax := widget.NewEntry()
	entryCreationMax.SetPlaceHolder("Année max")

	chk1 := widget.NewCheck("1 membre", nil)
	chk2 := widget.NewCheck("2 membres", nil)
	chk3 := widget.NewCheck("3 membres", nil)
	chk4 := widget.NewCheck("4+ membres", nil)

	btnAppliquer := widget.NewButton("Appliquer les filtres", func() {
		filtresActuels = FilterOptions{}

		if entryCreationMin.Text != "" {
			if v, err := strconv.Atoi(entryCreationMin.Text); err == nil {
				filtresActuels.CreationMin = v
			}
		}
		if entryCreationMax.Text != "" {
			if v, err := strconv.Atoi(entryCreationMax.Text); err == nil {
				filtresActuels.CreationMax = v
			}
		}

		var m []int
		if chk1.Checked {
			m = append(m, 1)
		}
		if chk2.Checked {
			m = append(m, 2)
		}
		if chk3.Checked {
			m = append(m, 3)
		}
		if chk4.Checked {
			m = append(m, 4, 5, 6, 7, 8)
		}
		filtresActuels.MembersCounts = m

		rafraichirGrille()
	})

	btnReset := widget.NewButton("Réinitialiser", func() {
		filtresActuels = FilterOptions{}
		rechercheActuelle = ""
		AfficherLeMenu()
	})

	colonneFiltres := container.NewVBox(
		widget.NewLabelWithStyle("Filtres", fyne.TextAlignCenter, fyne.TextStyle{Bold: true}),
		widget.NewSeparator(),
		widget.NewLabel("Date de création :"),
		entryCreationMin,
		entryCreationMax,
		widget.NewSeparator(),
		widget.NewLabel("Membres :"),
		chk1, chk2, chk3, chk4,
		widget.NewSeparator(),
		btnAppliquer,
		btnReset,
	)

	rafraichirGrille()

	split := container.NewHSplit(container.NewVScroll(colonneFiltres), conteneurDefilant)
	split.SetOffset(0.20)

	titrePrincipal := canvas.NewText("Artistes", color.White)
	titrePrincipal.TextSize = 40
	titrePrincipal.TextStyle = fyne.TextStyle{Bold: true}
	titrePrincipal.Alignment = fyne.TextAlignCenter

	contenuCentre := container.NewBorder(
		container.NewVBox(titrePrincipal, barreRecherche),
		nil, nil, nil,
		split,
	)

	maFenetre.SetContent(contenuCentre)
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

func FiltrerArtistes(artistes []List_artist, opts FilterOptions) []List_artist {
	var result []List_artist
	for _, a := range artistes {
		if !matchCreationDateList(a, opts) {
			continue
		}
		if !matchMembersCountList(a, opts) {
			continue
		}
		result = append(result, a)
	}
	return result
}

func matchCreationDateList(a List_artist, opts FilterOptions) bool {
	if opts.CreationMin != 0 && a.CreationDate < opts.CreationMin {
		return false
	}
	if opts.CreationMax != 0 && a.CreationDate > opts.CreationMax {
		return false
	}
	return true
}

func matchMembersCountList(a List_artist, opts FilterOptions) bool {
	if len(opts.MembersCounts) == 0 {
		return true
	}
	n := len(a.Members)
	for _, v := range opts.MembersCounts {
		if n == v {
			return true
		}
	}
	return false
}

func RechercherArtistes(artistes []List_artist, recherche string) []List_artist {
	if recherche == "" {
		return artistes
	}

	rechercheLower := strings.ToLower(recherche)
	var result []List_artist

	for _, a := range artistes {
		if strings.Contains(strings.ToLower(a.Name), rechercheLower) {
			result = append(result, a)
			continue
		}

		membreTrouve := false
		for _, membre := range a.Members {
			if strings.Contains(strings.ToLower(membre), rechercheLower) {
				membreTrouve = true
				break
			}
		}
		if membreTrouve {
			result = append(result, a)
			continue
		}

		if strings.Contains(strconv.Itoa(a.CreationDate), recherche) {
			result = append(result, a)
			continue
		}

		if strings.Contains(strings.ToLower(a.FirstAlbum), rechercheLower) {
			result = append(result, a)
			continue
		}
	}

	return result
}
