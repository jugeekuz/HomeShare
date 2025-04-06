import React from 'react'
import { useAuth } from '../contexts/AuthContext';
import avatarLogo from '../assets/icons/avatar.svg';
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
import { useNavigate } from 'react-router-dom';
const NavBar : React.FC = () => {
    const { logout, claims } = useAuth();
    const navigate = useNavigate();
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
            <NavbarBrand as="div" className='-ml-6 sm:ml-0'>
                <div
                    onClick={() => navigate("/")} 
                    className="flex justify-start items-center w-64 cursor-pointer">
                    <span className="font-signatra text-gray-100 text-[3rem]">
                        HomeShare
                    </span>
                </div>
            </NavbarBrand>
            
            <NavbarContent as="div" justify='end' className='-mr-6 sm:ml-0'>
                <Dropdown placement="bottom">
                    <DropdownTrigger>
                        <Avatar
                        isBordered
                        className='cursor-pointer'
                        color="secondary"
                        size="sm"
                        icon={
                            <img src={avatarLogo} alt="Avatar" className='' />
                        }
                        />
                    </DropdownTrigger>
                    <DropdownMenu aria-label="Profile Actions" variant="flat">
                        <DropdownItem key="profile" className="h-14 gap-2">
                            <p className="font-semibold">Signed in as</p>
                            <p className="font-semibold">{claims?.user_id}</p>
                        </DropdownItem>
                        
                        <DropdownItem key="logout" color="danger">
                            <div 
                                onClick={logout}
                                className="flex"
                            >    
                            Log Out
                            </div>
                        </DropdownItem>
                    </DropdownMenu>
                </Dropdown>
            </NavbarContent>
        </Navbar>
    )
}

export default NavBar