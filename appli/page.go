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

type FilterOptions struct {
    CreationMin   int
    CreationMax   int
    MembersCounts []int
}

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

    entryCreationMin := widget.NewEntry()
    entryCreationMin.SetPlaceHolder("Année min")
    entryCreationMax := widget.NewEntry()
    entryCreationMax.SetPlaceHolder("Année max")

    chk1 := widget.NewCheck("1 membre", nil)
    chk2 := widget.NewCheck("2 membres", nil)
    chk3 := widget.NewCheck("3 membres", nil)
    chk4 := widget.NewCheck("4 membres", nil)

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
            m = append(m, 4)
        }
        filtresActuels.MembersCounts = m

        AfficherLeMenu()
    })

    btnReset := widget.NewButton("Réinitialiser", func() {
        filtresActuels = FilterOptions{}
        AfficherLeMenu()
    })

    colonneFiltres := container.NewVBox(
        widget.NewLabelWithStyle("Filtres", fyne.TextAlignCenter, fyne.TextStyle{Bold: true}),
        widget.NewSeparator(),
        widget.NewLabel("Date de création"),
        entryCreationMin,
        entryCreationMax,
        widget.NewSeparator(),
        widget.NewLabel("Nombre de membres"),
        chk1,
        chk2,
        chk3,
        chk4,
        widget.NewSeparator(),
        btnAppliquer,
        btnReset,
    )

    grilleArtistes := container.NewGridWithColumns(4)

    artistesAffiches := FiltrerArtistes(tousLesArtistes, filtresActuels)

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

    conteneurDefilant := container.NewVScroll(grilleArtistes)

    split := container.NewHSplit(colonneFiltres, conteneurDefilant)
    split.SetOffset(0.25)

    titrePrincipal := canvas.NewText("Nos Artistes", color.Black)
    titrePrincipal.TextSize = 60
    titrePrincipal.TextStyle = fyne.TextStyle{Bold: true}
    titrePrincipal.Alignment = fyne.TextAlignCenter

    contenuFinal := container.NewBorder(titrePrincipal, nil, nil, nil, split)
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
    grandeImage.SetMinSize(fyne.NewSize(400, 400))

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

    infosContainer := container.NewVBox(
        lblNom,
        widget.NewSeparator(),
        lblDate,
        lblAlbum,
        widget.NewSeparator(),
        lblMembresTitre,
        lblMembresListe,
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
