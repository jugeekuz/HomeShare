import React, { useState, useEffect } from "react";
import { useLocation } from "react-router-dom";
import config from '../configs/config';
import FileList from './FileList';
import { Button } from "@heroui/react";
import { FileContext, useFileContext } from "../contexts/FileContext";
import { FileStore, FileItem, FileMeta } from "../types";
import { IoMdDownload } from "react-icons/io";

const files: FileStore = {
    "someImage.jpg" : {
        file : new File([], "someImage.jpg", { type: "text/plain" }),
        fileMeta: {
            fileId:         "someFileId",
            fileName:       "someImage",
            fileExtension:  ".jpg",
            md5Hash:        ""
        }
    },
    "someImageYes.jpg" : {
        file : new File([], "someImageYes.jpg", { type: "text/plain" }),
        fileMeta: {
            fileId:         "someFileId",
            fileName:       "someImageYes",
            fileExtension:  ".jpg",
            md5Hash:        ""
        }
    },
    "someImageNo.jpg" : {
        file : new File([], "someImageNo.jpg", { type: "text/plain" }),
        fileMeta: {
            fileId:         "someFileId",
            fileName:       "someImageNo",
            fileExtension:  ".jpg",
            md5Hash:        ""
        }
    },
}

const DownloadComponent: React.FC = () => {
    const [folderName, setFolderName] = useState<string>("Vytina")
    const { setFiles } = useFileContext();
    
    useEffect(() => {
        setFiles(files)
    }, [])

    return (
        <div className="flex flex-col justify-center items-center w-full h-full">
            <div className="flex flex-col justify-start items-center mt-1">
                <div className="flex justify-center items-center">
                    <span className="font-brsonoma font-normal text-gray-950 text-md mr-1">
                        Sharing Folder <span className="bg-secondary-gradient bg-clip-text text-transparent font-bold">{folderName}</span>
                    </span>
                </div>
                <span className="font-brsonoma font-light text-gray-500 text-xs">
                    Explore what others shared or upload your files
                </span>
            </div>
            <div className="flex w-[87%] mt-3 mb-2 border-dashed border-2 border-gray-300 rounded-lg px-2 pt-2 pb-1">
                <FileList/>
            </div>
            <div className="flex w-[90%]">
                <Button 
                    color="primary"
                    className="text-[13px] bg-primary-gradient rounded-md w-full"
                    size="md"
                >
                    <IoMdDownload 
                        size={15}
                    />Download all in .zip format
                </Button>
            </div>
        </div>
    );
};

export default DownloadComponent;
