/** @type {import('tailwindcss').Config} */
module.exports = {
    content: [
        "./app/views/**/*.templ",
        "./app/**/*.go",
    ],
    theme: {
        extend: {
        colors: {
            // Colores personalizados.
            scifi: 'rgba(0, 87, 133, 1)', 
            'scifi-30': 'rgba(0, 87, 133, 0.3)',
            'white-80': 'rgba(255, 255, 255, 0.8)',
            background: 'white',
        },
        fontFamily: {
            // Fuente personalizada.
            ibm: ['"IBM Plex Sans"', 'sans-serif'],
        },
        backgroundImage: {
            'main-bg': "url('/static/img/fondo.png')",
            'search-icon': "url('/static/img/icono_busqueda_blanco.svg')",
        }
        },
    },
    plugins: [],
}