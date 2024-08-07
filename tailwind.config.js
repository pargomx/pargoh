/** @type {import('tailwindcss').Config} */
module.exports = {
  content: [
	"./htmltmpl/*.html",
	"./htmltmpl/**/*.html",
	
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

