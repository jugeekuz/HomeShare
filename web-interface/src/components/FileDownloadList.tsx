import React, {useEffect, useState} from 'react';

import PerfectScrollbar from 'react-perfect-scrollbar';
import { IoMdDownload } from "react-icons/io";
import { HiXMark } from "react-icons/hi2";
import "react-perfect-scrollbar/dist/css/styles.css";
import { ScrollShadow } from '@heroui/react';

import { FileDownloadItem, FileIconType } from '../types';
import { useSearchParams } from 'react-router-dom';
import { useFileDownloadContext } from '../contexts/FileDownloadContext';
import { useNotificationContext } from '../contexts/NotificationContext.tsx';
import DownloadItem from './DownloadItem.tsx';

const FileBox : React.FC<{fileName: string}> = ({fileName}) => {
    const { files, downloadFile } = useFileDownloadContext();
    const [searchParams] = useSearchParams();
    const { notifyError } = useNotificationContext();
    const [file, setFile] = useState<FileDownloadItem | null>(null);
    const [folderId, setFolderId] = useState<string | null>(null);

    useEffect(() => {
        if (!files) return;
        if (!(fileName in files)) {
            notifyError("File Error", `File name ${fileName} wasn't found in the store`);
            return;
        }
        setFile(files[fileName]);        
    }, [])

    useEffect(() => {
        setFolderId(searchParams.get("folder-id"));
    }, [searchParams])

    const downloadCb = async () => {
        if (!folderId) return;
        downloadFile(fileName, folderId);
    }
    const convertBytes = (bytes: number) : string => {
        if (bytes === 0) return "0 B";

        const mapping = ["B", "kB", "MB", "GB", "TB"];
        const n = Math.floor(Math.log(bytes) / Math.log(1024));
        
        if (n >= mapping.length) throw new Error("Byte string is too large to convert.");

        const value = bytes / (1024 ** n);
        
        const formattedValue = parseFloat(value.toFixed(2));

        return `${formattedValue} ${mapping[n]}`;
    }

    return (
    <div className="flex flex-col justify-center items-center w-full h-[4.5rem] rounded-lg bg-contrast px-2 mb-1 border-1 border-gray-200" key={fileName}>
        <div className="relative flex flex-row justify-start items-center w-full p-2">
            <div className="flex justify-center items-center w-7 h-7">
                <DownloadItem fileName={fileName}  />
            </div>
            <div className="flex flex-col w-full px-3 justify-center overflow-hidden text-nowrap">
                <span className="font-linik text-[13px] text-gray-700 font-bold">
                    {fileName}
                </span>
                <div className="flex flex-row justify-start items-center">
                    <span className="text-xs text-gray-600">
                        {files && convertBytes(file?.fileSize ? file.fileSize : 0)}
                    </span>
                </div>
            </div>
            {/* Here */}
            <div className="absolute flex justify-center items-center rounded-full border bg-wsecondary w-8 h-8 right-3 top-1/2 -translate-y-1/2 cursor-pointer">
                <IoMdDownload
                    onClick={downloadCb}
                    size={15}
                    className='text-primary'
                />
            </div>
        </div>
    </div>)
}

const FileDownloadList : React.FC = () => {
    const { files } = useFileDownloadContext();

    return (
        <div className="flex flex-col items-center justify-center w-full gap-2 max-h-[10.5rem]">
            {files && Object.keys(files).length > 0 && (
                <PerfectScrollbar
                    className="w-full h-full"
                    options={{
                        wheelSpeed: 1,
                        suppressScrollX: true,
                        minScrollbarLength: 30,
                    }}
                >   
                    <ScrollShadow 
                        className='w-full h-full'
                        visibility='bottom'
                        size={12}
                    >
                        <div className="flex flex-col items-start w-full gap-[2px]">
                            {Object.keys(files).map((fileName) => (
                                <FileBox key={fileName} fileName={fileName} />
                            ))}
                        </div>
                    </ScrollShadow>
                </PerfectScrollbar>
            )}
        </div>
    );
}

export default FileDownloadList;