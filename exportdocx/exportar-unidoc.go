package exportdocx

import (
	"fmt"
	"monorepo/dhistorias"
	"monorepo/ust"
	"strings"

	"github.com/pargomx/gecko/gko"
	"github.com/unidoc/unioffice/common"
	"github.com/unidoc/unioffice/document"
	"github.com/unidoc/unioffice/measurement"
	"github.com/unidoc/unioffice/schema/soo/wml"
)

// ================================================================ //
// ========== unidoc ============================================== //

func ExportarDocx(personaID int, repo dhistorias.Repo, filepath string) error {
	op := gko.Op("ExportarDocx").Ctx("personaID", personaID)

	persona, err := repo.GetPersona(personaID)
	if err != nil {
		return op.Err(err)
	}
	proyecto, err := repo.GetProyecto(persona.ProyectoID)
	if err != nil {
		return op.Err(err)
	}
	historias, err := repo.ListNodoHistoriasByPadreID(persona.PersonaID)
	if err != nil {
		return err
	}

	// Comenzar documento
	doc := document.New()
	defer doc.Close()

	// Título del documento
	addTitulo(doc, 0, persona.Nombre+": "+proyecto.Titulo)

	// Capítulos
	for _, h := range historias {
		historia, err := dhistorias.GetHistoria(h.HistoriaID, dhistorias.GetDescendientes|dhistorias.GetReglas|dhistorias.GetTramos, repo)
		if err != nil {
			return err
		}

		// Título del capítulo.
		addTitulo(doc, 1, historia.Historia.Titulo)

		// Texto normal.
		addTextoCuerpo(doc, historia.Historia.Objetivo)
		addReglas(doc, historia.Reglas)
		addTextoCuerpo(doc, historia.Historia.Descripcion)
		err = addTramos(doc, historia.Tramos)
		if err != nil {
			return err
		}

		// Subcapítulos.
		for _, descendiente := range historia.Descendientes {
			addHistoriaDocx(doc, descendiente, 2)
		}
	}

	return doc.SaveToFile(filepath)
}

type txtEstilo struct {
	texto  string
	bold   bool
	italic bool
}

// Acepta una línea de texto con _cursiva_ y *negrita*.
func separarTextoEnRunsConEstilo(texto string) []txtEstilo {
	var runs []txtEstilo

	txt := txtEstilo{}

	for _, char := range texto {

		if char != '*' && char != '_' {
			txt.texto += string(char)

		} else if char == '*' && txt.texto == "" {
			txt.bold = !txt.bold

		} else if char == '*' && txt.texto != "" {
			if !txt.bold {
				runs = append(runs, txt)
				txt = txtEstilo{bold: true, italic: txt.italic}

			} else if txt.bold {
				runs = append(runs, txt)
				txt = txtEstilo{italic: txt.italic}
			}

		} else if char == '_' && txt.texto == "" {
			txt.italic = !txt.italic

		} else if char == '_' {
			if !txt.italic {
				runs = append(runs, txt)
				txt = txtEstilo{italic: true, bold: txt.bold}

			} else if txt.italic {
				runs = append(runs, txt)
				txt = txtEstilo{bold: txt.bold}
			}
		}
	}
	if txt.texto != "" {
		runs = append(runs, txt)
	}
	return runs
}

// points is constant of the text height, it's 12 points.
const points measurement.Distance = 12

func addTitulo(doc *document.Document, nivel int, txt string) error {
	if txt == "" {
		return nil
	}
	pr := doc.AddParagraph()
	if nivel == 0 {
		pr.SetStyle("Title")
	} else if nivel <= 6 {
		pr.SetStyle(fmt.Sprintf("Heading%d", nivel))
	} else {
		return fmt.Errorf("nivel de título inválido: %d", nivel)
	}
	pr.SetAfterSpacing(1.5 * points)
	putTextConEstilo(&pr, txt, false)
	return nil
}

func putTextConEstilo(pr *document.Paragraph, txt string, newLines bool) {
	if txt == "" {
		return
	}
	for _, segmentoTxt := range separarTextoEnRunsConEstilo(txt) {
		run := pr.AddRun()
		run.Properties().SetFontFamily("Times New Roman")
		if segmentoTxt.bold {
			run.Properties().Bold()
		}
		if segmentoTxt.italic {
			run.Properties().Italic()
		}
		if newLines {
			// Separar líneas dentro del mismo párrafo.
			lineasTxt := strings.Split(segmentoTxt.texto, "\n")
			for j, lineaTxt := range lineasTxt {

				// Si son menos de 10 palabras no justificar
				if len(strings.Fields(lineaTxt)) < 10 {
					pr.SetAlignment(wml.ST_JcLeft)
				}

				run.AddText(lineaTxt)
				// No agregar a la última línea ni si solo hay una.
				if j < len(lineasTxt)-1 && len(lineasTxt) > 1 {
					run.AddBreak()
				}
			}
		} else {
			run.AddText(segmentoTxt.texto)
		}
	}
}

func addImagen(doc *document.Document, imgPath string) error {
	if imgPath == "" {
		return nil
	}
	imgPath = "imagenes/" + imgPath
	img, err := common.ImageFromFile(imgPath)
	if err != nil {
		return err
	}
	imgRef, err := doc.AddImage(img)
	if err != nil {
		return err
	}

	// Agregar en línea
	pr := doc.AddParagraph()
	pr.SetAlignment(wml.ST_JcCenter)
	pr.SetAfterSpacing(1.5 * points)

	run := pr.AddRun()

	inlineImg, err := run.AddDrawingInline(imgRef)
	if err != nil {
		return fmt.Errorf("unable to add inline image: %w", err)
	}

	// Convertir los pixeles del width a inches y limitar el tamaño.
	const maxWidth measurement.Distance = 6.5 * measurement.Inch // 8.5 pulgadas de ancho, menos 2 pulgadas de margen.
	const maxHeigth measurement.Distance = 9 * measurement.Inch  // 11 pulgadas de alto, menos 2 pulgadas de margen.

	// 96 DPI
	imgW := measurement.Distance(imgRef.Size().X) / 96 * measurement.Inch
	imgH := measurement.Distance(imgRef.Size().Y) / 96 * measurement.Inch

	if imgW >= maxWidth {
		imgW = maxWidth
		imgH = imgRef.RelativeHeight(maxWidth)
	}
	if imgH >= maxHeigth {
		imgH = maxHeigth
		imgW = imgRef.RelativeWidth(maxHeigth)
	}

	inlineImg.SetSize(imgW, imgH)

	// run.AddBreak()
	// run.AddText( fmt.Sprintf("Original %v*%vpx, %v*%vin Resized: %v*%vin",
	// 	imgRef.Size().X, imgRef.Size().Y,
	// 	float64(imgRef.Size().X)/96, float64(imgRef.Size().Y)/96,
	// 	imgW/measurement.Inch, imgH/measurement.Inch,
	// ))

	return nil
}

func addTextoCuerpo(doc *document.Document, txt string) {
	if txt == "" {
		return
	}
	// Separar párrafos.
	for _, parrafoTxt := range strings.Split(txt, "\n\n") {
		pr := doc.AddParagraph()
		pr.SetFirstLineIndent(0.5 * measurement.Inch)             // Indentación de primera línea.
		pr.SetLineSpacing(1.5*points, wml.ST_LineSpacingRuleAuto) // Espacio entre líneas
		pr.SetAlignment(wml.ST_JcBoth)                            // Justificar texto.
		pr.SetAfterSpacing(1.5 * points)                          // Espacio después del párrafo.

		putTextConEstilo(&pr, parrafoTxt, true)
	}
}

func addReglas(doc *document.Document, reglas []ust.Regla) {
	if len(reglas) == 0 {
		return
	}
	ndBullet := doc.Numbering.Definitions()[0]
	for _, regla := range reglas {
		pr := doc.AddParagraph()
		pr.SetNumberingDefinition(ndBullet)
		pr.SetNumberingLevel(0)

		pr.SetLineSpacing(1.5*points, wml.ST_LineSpacingRuleAuto)
		pr.SetAfterSpacing(1 * points)
		putTextConEstilo(&pr, regla.Texto, true)
	}
}

func addTramos(doc *document.Document, tramos []ust.Tramo) error {
	if len(tramos) == 0 {
		return nil
	}
	// Create numbering definition.
	nd := doc.Numbering.AddDefinition()

	// Add level to number definition with decimal format.
	lvl := nd.AddLevel()
	lvl.SetFormat(wml.ST_NumberFormatDecimal)
	lvl.SetAlignment(wml.ST_JcLeft)
	lvl.Properties().SetLeftIndent(1.5 * points)

	// Sets the numbering level format.
	lvl.SetText("%1.")

	for _, regla := range tramos {
		pr := doc.AddParagraph()
		pr.SetNumberingDefinition(nd)
		pr.SetNumberingLevel(0)

		pr.SetLineSpacing(1.5*points, wml.ST_LineSpacingRuleAuto)
		pr.SetAfterSpacing(1.5 * points)

		if regla.Imagen == "" {
			putTextConEstilo(&pr, regla.Texto, true)

		} else {
			err := addImagen(doc, regla.Imagen)
			if err != nil {
				return err
			}
			putTextConEstilo(&pr, regla.Texto, false)
		}
	}
	return nil
}

func addHistoriaDocx(doc *document.Document, his dhistorias.HistoriaRecursiva, nivel int) error {
	err := addTitulo(doc, nivel, his.Titulo)
	if err != nil {
		return err
	}
	addTextoCuerpo(doc, his.Objetivo)
	addReglas(doc, his.Reglas)
	addTextoCuerpo(doc, his.Descripcion)
	err = addTramos(doc, his.Tramos)
	if err != nil {
		return err
	}

	for _, descendiente := range his.Descendientes {
		err := addHistoriaDocx(doc, descendiente, nivel+1)
		if err != nil {
			return err
		}
	}
	return nil
}
