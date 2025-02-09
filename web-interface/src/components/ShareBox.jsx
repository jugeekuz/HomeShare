import React from 'react'

import { Button } from '@heroui/button'

import FileUploader from './FileUploader'

const ShareBox = () => {
    return (
        <div className="flex h-full w-full justify-center items-center ">
            <div className="flex flex-col justify-end items-center h-[40rem] w-[32rem] max-w-[85%] max-h-[60%] bg-white rounded-xl shadow-lg -mt-20">
                {/* <div className="flex w-full justify-center items-center h-16">
                    <span className=" text-xl">
                      
                    </span>
                </div> */}

                <FileUploader/>
                <div className="flex justify-center items-center w-[80%] mb-8">
                    <Button color='primary' className='text-md w-[80%]' size='lg'>
                        Transfer Files
                    </Button>
                </div>
            </div>
        </div>
    )
}

export default ShareBox