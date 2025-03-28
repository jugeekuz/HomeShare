import {heroui} from '@heroui/theme'
/** @type {import('tailwindcss').Config} */
export default {
    content: [
        "./src/**/*.{js,jsx,ts,tsx}",
        "./node_modules/@heroui/theme/dist/**/*.{js,ts,jsx,tsx}",
    ],
    theme: {
        extend: {
            fontFamily: {
                signatra: ['Signatra', 'cursive'],
                eina: ['Eina', 'sans-serif'],
                linik: ["LinikSans", "sans-serif"],
                roboto: ['Roboto', 'sans-serif'],
                coolvetica: ['Coolvetica', 'sans-serif'],
                brsonoma: ['BRSonoma', 'sans-serif'],
            },
            colors: {
                wprimary: '#fcfbfc',
                wsecondary: '#ffffff',
                contrast: '#f3f5f7',
            },       
            backgroundImage: {
                'primary-gradient': 'linear-gradient(to top right, #3b82f6, #1d4ed8)',
                'secondary-gradient': 'linear-gradient(to right, #3b82f6, #9333ea)'
            },
            
        },
    },
    plugins: [
        {
        "tailwindConfig": "./tailwind.config.js"
        },
        heroui(),
    ],
}

