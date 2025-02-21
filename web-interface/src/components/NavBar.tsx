import React from 'react'

const NavBar : React.FC = () => {
    return (
        <nav className="flex flex-row justify-start items-center bg-transparent h-16 w-full">
            <div className="flex justify-center items-center w-64">
                <span className="font-signatra text-gray-100 text-[3rem] sh">
                    HomeShare
                </span>
            </div>
        </nav>
    )
}

export default NavBar