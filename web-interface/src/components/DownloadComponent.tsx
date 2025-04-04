import React, { useState, useEffect } from "react";
import { IoMdDownload } from "react-icons/io";
import { useSearchParams } from "react-router-dom";
import { Button, Spinner } from "@heroui/react";

import config from "../configs/config";
import { useFileDownloadContext } from "../contexts/FileDownloadContext.tsx";
import { FileDownloadItem, GetSharingFilesParams } from "../types";
import { useNotificationContext } from "../contexts/NotificationContext";
import api from "../api/api";

import FileDownloadList from './FileDownloadList';

const DownloadComponent: React.FC = () => {
    const [empty, setEmpty] = useState<boolean>(false);
    const {notifyError} = useNotificationContext();
    const { files, setFiles, addFile, downloadZip } = useFileDownloadContext();
    const [searchParams] = useSearchParams();
    const [folderId, setFolderId] = useState<string | null>(null);
    const [loading, setLoading] = useState<boolean>(false);

    useEffect(() => {
        if (!files) {
            setEmpty(true);
        } else {
            setEmpty(false);
        };
    }, [files])

    useEffect(() => {
        setFolderId(searchParams.get("folder-id"));
    }, [searchParams])

    useEffect(() => {
        if (!folderId) return;
        const fetchData = async () => {
            try {    
                setLoading(true);
                const params : GetSharingFilesParams = {
                    folder_id: folderId,
                }
                const response = await api.get(config.GET_SHARING_FILES_URL, {
                    params: params,
                });
                setLoading(false);
                const data = response?.data;

                if (!data) {
                    notifyError("Fetching Files Error", "Received invalid response while fetching file data from the server");
                    return;
                }
                if ( !data?.files ) {
                    setFiles(null);
                    return;
                }
                for (const file of data.files) {
                    const fileNameWoExt = file["file_name"];
                    const fileExtension = file["file_extension"];
                    if (fileExtension === ".zip") continue;
                    const fileSize = file["file_size"];

                    const fileName = `${fileNameWoExt}${fileExtension}`

                    const fileItem : FileDownloadItem = {
                        fileNameWoExt: fileNameWoExt,
                        fileExtension: fileExtension,
                        fileSize: fileSize,
                    }
                    addFile(fileItem, fileName);
                }


            } catch (error) {
                setLoading(false);
                return;
            }
        }
        fetchData();

    }, [folderId])

    return (
        <div className="flex flex-col justify-center items-center w-full h-full">
            
            <div className="flex items-center justify-center w-[87%] mt-3 mb-2 border-dashed border-2 border-gray-300 rounded-lg px-2 pt-2 pb-1">
                {   !empty ?
                        loading ?
                            <Spinner color="default"/>
                        : 
                            <FileDownloadList/>
                    : 
                        <div className="flex w-full h-[12rem] justify-center items-center">
                            <div className="flex flex-col justify-start items-center mt-1">
                                <div className="flex justify-center items-center">
                                    <span className="font-brsonoma font-normal text-gray-950 text-sm mr-1">
                                        Sharing Folder is Empty
                                    </span>
                                </div>
                                <span className="font-brsonoma font-light text-gray-500 text-xs">
                                    Upload your files to get started
                                </span>
                            </div>
                        </div>
                }
            </div>
            <div className="flex w-[90%]">
                <Button 
                    isDisabled={empty || loading}
                    color="primary"
                    className="text-[13px] bg-primary-gradient rounded-md w-full mt-1"
                    size="md"
                    onPress={() => {
                        if (!folderId) return;
                        downloadZip(folderId)
                    }}
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
