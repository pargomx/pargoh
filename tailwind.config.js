/** @type {import('tailwindcss').Config} */
module.exports = {
  content: [
	"./htmltmpl/*.html",
	"./htmltmpl/**/*.html",
	"./assets/js/main.js"
	
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

