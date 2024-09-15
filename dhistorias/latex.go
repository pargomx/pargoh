package dhistorias

import (
	"fmt"
	"strings"

	"github.com/pargomx/gecko/gko"
	"github.com/rwestlund/gotex"

	_ "embed"
)

// ================================================================ //
// ========== LaTeX =============================================== //

//go:embed latex-header.tex
var latexHeader string

func (tex *latexBuilder) writeHeader(titulo string) {
	// tex.addCommand("\\documentclass[12pt]{book}")
	// tex.addCommand("\\usepackage{graphicx}")
	// tex.addCommand("\\graphicspath{ {/home/andrew/pargodata/imagenes/} }")
	tex.buf.WriteString(latexHeader)
	tex.buf.WriteString("\n\n")

	tex.addCommand("\\begin{document}")

	tex.addCommandf("\\title{%v}", titulo)
	tex.addCommand("\\author{Documento autogenerado}")
	tex.addCommand("\\date{\\today}")

	tex.addCommand("\\maketitle")
	tex.addCommand("\\tableofcontents")
}

func (tex *latexBuilder) writePersona(per PersonaExport) {
	tex.addChapter(per.Persona.Nombre)
	tex.addParrafo(per.Persona.Descripcion)
	for _, his := range per.Historias {
		tex.addSection(his.Historia.Titulo)
		tex.addParrafo(his.Historia.Objetivo)
		for _, tramo := range his.Tramos {
			if tramo.Imagen != "" {
				tex.addImagen(tramo.Imagen, tramo.Texto)
			} else {
				tex.addParrafo(tramo.Texto)
			}
		}

		for _, his2 := range his.Historias {
			tex.addSubsection(his2.Historia.Titulo)
			tex.addParrafo(his2.Historia.Objetivo)
			for _, tramo := range his2.Tramos {
				if tramo.Imagen != "" {
					tex.addImagen(tramo.Imagen, tramo.Texto)
				} else {
					tex.addParrafo(tramo.Texto)
				}
			}

			for _, his3 := range his.Historias {
				tex.addSubSubection(his3.Historia.Titulo)
				tex.addParrafo(his3.Historia.Objetivo)
				for _, tramo := range his3.Tramos {
					if tramo.Imagen != "" {
						tex.addImagen(tramo.Imagen, tramo.Texto)
					} else {
						tex.addParrafo(tramo.Texto)
					}
				}
			}
		}
	}
}

func GetPersonaLaTeX(repo Repo, personaID int) (latexBuilder, error) {
	const op = "ExportarPersonaPDF"
	tex := latexBuilder{}
	per, err := GetPersonaExport(personaID, repo)
	if err != nil {
		return tex, gko.Err(err).Op(op)
	}
	tex.writeHeader(per.Proyecto.Titulo)
	tex.writePersona(*per)
	tex.addCommand("\\end{document}")
	return tex, nil
}

func GetProyectoLaTeX(repo Repo, proyectoID string) (latexBuilder, error) {
	tex := latexBuilder{}
	pry, err := GetProyectoExport(proyectoID, repo)
	if err != nil {
		return tex, gko.Err(err).Op("ExportarTeX")
	}
	tex.writeHeader(pry.Proyecto.Titulo)
	for _, per := range pry.Personas {
		tex.writePersona(per)
	}
	tex.addCommand("\\end{document}")
	return tex, nil
}

// ================================================================ //
// ========== BUILDER ============================================= //

type latexBuilder struct {
	buf strings.Builder
}

// Map of special LaTeX characters and their escaped equivalents
var latexEscapes = map[rune]string{
	'\\': `\textbackslash`,
	'{':  `\{`,
	'}':  `\}`,
	'$':  `\$`,
	'&':  `\&`,
	'#':  `\#`,
	'_':  `\_`,
	'%':  `\%`,
	'~':  `\textasciitilde`,
	'^':  `\textasciicircum`,
}

func latexEscape(input string) string {
	var escaped strings.Builder
	for _, char := range input {
		if escapedChar, ok := latexEscapes[char]; ok {
			escaped.WriteString(escapedChar)
		} else {
			escaped.WriteRune(char)
		}
	}
	return escaped.String()
}

// ================================================================ //
// ========== INPUTS ============================================== //

func (b *latexBuilder) addCommand(cmd string) {
	b.buf.WriteString(cmd)
	b.buf.WriteString("\n\n")
}
func (b *latexBuilder) addCommandf(cmd string, txt string) {
	b.buf.WriteString(fmt.Sprintf(cmd, latexEscape(txt)))
	b.buf.WriteString("\n\n")
}
func (b *latexBuilder) addParrafo(txt string) {
	b.buf.WriteString(latexEscape(txt))
	b.buf.WriteString("\n\n")
}
func (b *latexBuilder) addChapter(txt string) {
	b.buf.WriteString("\\chapter{")
	b.buf.WriteString(latexEscape(txt))
	b.buf.WriteString("}\n\n")
}
func (b *latexBuilder) addSection(txt string) {
	b.buf.WriteString("\\section{")
	b.buf.WriteString(latexEscape(txt))
	b.buf.WriteString("}\n\n")
}
func (b *latexBuilder) addSubsection(txt string) {
	b.buf.WriteString("\\subsection{")
	b.buf.WriteString(latexEscape(txt))
	b.buf.WriteString("}\n\n")
}
func (b *latexBuilder) addSubSubection(txt string) {
	b.buf.WriteString("\\subsection{")
	b.buf.WriteString(latexEscape(txt))
	b.buf.WriteString("}\n\n")
}
func (b *latexBuilder) addImagen(file, label string) {
	b.addParrafo(label)
	b.addCommandf("\\includegraphics[width=8cm]{%v}", file)
	b.buf.WriteString("\n\n")
}

// ================================================================ //
// ========== OUTPUTS ============================================= //

func (b *latexBuilder) String() string {
	return b.buf.String()
}

func (b *latexBuilder) ToPDF() ([]byte, error) {
	pdf, err := gotex.Render(b.buf.String(), gotex.Options{
		Command:   "pdflatex",
		Runs:      1,
		Texinputs: "/home/andrew/pargodata/imagenes",
	})
	if err != nil {
		return nil, gko.Err(err).Op("ExportarPDF")
	}
	return pdf, nil
}
