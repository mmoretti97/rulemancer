TEX=slides.tex
PDF=slides.pdf

.PHONY: pdf clean

pdf:
	latexmk -pdf -interaction=nonstopmode -halt-on-error $(TEX)

clean:
	latexmk -C
	rm -f *.vrb *.nav *.snm $(PDF)
