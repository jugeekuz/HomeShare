import React from 'react'

import AvatarLogo from '../assets/icons/avatar.svg';
import {
    Navbar, 
    NavbarBrand, 
    NavbarContent, 
    NavbarItem, 
    Link, 
    Button, 
    AvatarIcon,
    Avatar, 
    Dropdown, 
    DropdownItem, 
    DropdownMenu, 
    DropdownTrigger
} from "@heroui/react";
const NavBar : React.FC = () => {
    return (
        <Navbar
            shouldHideOnScroll
            // isBordered
            isBlurred={false}
            className='bg-transparent'
            classNames={{
                wrapper: "max-w-full w-full px-[3.3rem]"
            }}
            position='sticky'
        >
            <NavbarBrand as="div">
                <div className="flex justify-start items-center w-64">
                    <span className="font-signatra text-gray-100 text-[3rem] sh">
                        HomeShare
                    </span>
                </div>
            </NavbarBrand>
            
            <NavbarContent as="div" justify='end'>
                <Dropdown placement="bottom">
                    <DropdownTrigger>
                        <Avatar
                        isBordered
                        className='cursor-pointer'
                        color="secondary"
                        size="sm"
                        icon={
                            <img src={AvatarLogo} alt="Avatar" className='' />
                        }
                        />
                    </DropdownTrigger>
                    <DropdownMenu aria-label="Profile Actions" variant="flat">
                        <DropdownItem key="profile" className="h-14 gap-2">
                            <p className="font-semibold">Signed in as</p>
                            <p className="font-semibold">zoey@example.com</p>
                        </DropdownItem>
                        
                        <DropdownItem key="logout" color="danger">
                        Log Out
                        </DropdownItem>
                    </DropdownMenu>
                </Dropdown>
            </NavbarContent>
        </Navbar>
    )
}

export default NavBar