/** @type {import('tailwindcss').Config} */
module.exports = {
  content: [
	"./pargoh/plantillas/*.html",
	"./pargoh/plantillas/**/*.html",
	
],
  theme: {
    extend: {
		cursor: {
			'edit': 'url("/assets/img/pen.svg"), pointer',
		},
	},
  },
  plugins: [],
}

