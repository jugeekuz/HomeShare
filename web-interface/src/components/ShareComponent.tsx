import React, {useState} from 'react'

import { Input, DatePicker } from '@heroui/react'
import { Button } from '@heroui/button';
import { FaLink } from "react-icons/fa6";

const ShareComponent = () => {

    const handleSumbit = () => {
        return undefined
    }

    return (
        <div className="flex flex-col w-full h-full justify-center items-center">
            <form className="flex flex-col items-center justify-center w-full gap-3 my-2">
                <div className="flex flex-col w-4/5 mt-1">
                    <label className="text-sm font-linik text-gray-00 mb-1">Sharing Folder Name</label>
                    <Input
                        className='w-full'
                        placeholder="Enter a name"
                        startContent={
                            <FaLink size={20} className='text-gray-500 mr-2'/>
                        }
                        type="email"
                    />
                </div>
                <div className="flex flex-col w-4/5 mb-1">
                    <label className="text-sm font-normal text-gray-700 mb-1 ">Sharing Folder Expiration</label>
                    <DatePicker 
                        className="w-full"
                        size='md'
                    />
                </div>
                <Button
                    className=''
                    color='primary'
                >
                    <FaLink size={20} className='text-white mr-1'/>Create Sharing Folder
                </Button>
            </form>
        </div>
    )
}

export default ShareComponent