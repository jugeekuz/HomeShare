import React, {useRef , useState} from 'react'

import { FileProvider } from '../contexts/FileContext'

import FileLoader from './FileLoader'
import FileUploader from './FileUploader'
const ShareBox = () => {
    const fileInputRef = useRef(null);

    return (
        <div className="flex h-full w-full justify-center items-center ">
            <div className="flex flex-col justify-center items-center h-[20rem] w-[25rem] max-w-[85%] max-h-[60%] bg-white rounded-xl shadow-lg -mt-20">
                
                <FileProvider>
                    <div className="flex flex-col items-center justify-center w-full h-3/5">
                        <div className="flex flex-col justify-center items-center text-center mt-10">
                            <FileLoader/>
                            {/* <span className="">
                                Upload Files
                            </span> */}
                        </div>
                    </div>
                    <div className="flex flex-col justify-center gap-3 items-center w-full h-2/5 py-5">
                        <FileUploader/>
                    </div>
                </FileProvider>
            </div>
        </div>
    )
}

export default ShareBox