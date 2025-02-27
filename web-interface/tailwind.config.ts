import {heroui} from '@heroui/theme'
/** @type {import('tailwindcss').Config} */
export default {
    content: [
        "./src/**/*.{js,jsx,ts,tsx}",
        "./node_modules/@heroui/theme/dist/components/(button|divider|progress|ripple|spinner).js"
    ],
    theme: {
        extend: {
        fontFamily: {
            signatra: ['Signatra', 'cursive'],
            linik: ["LinikSans", "sans-serif"],
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

