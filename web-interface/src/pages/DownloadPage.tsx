import React from "react";
import { useLocation } from "react-router-dom";
import config from '../configs/config';
import { FileDownloadIcon } from "../components/FileIcons";


const DownloadPage: React.FC = () => {
    const location = useLocation();
    const queryParams = new URLSearchParams(location.search);
    const file = queryParams.get('file');

    const downloadItems = async () => {
        try {
            if (!file) return;
            const response = await fetch(`${config.DOWNLOAD_URL}?file=${file}`, {
                method: "GET",
            });
    
            if (!response.ok) {
                throw new Error(`HTTP error! status: ${response.status}`);
            }

            const blob = await response.blob();

            const url = window.URL.createObjectURL(blob);

            const a = document.createElement("a");
            a.style.display = "none";
            a.href = url;
            a.download = file;
    
            document.body.appendChild(a);
            a.click();

            window.URL.revokeObjectURL(url);
            document.body.removeChild(a);
        } catch (error) {
            console.error("Download failed:", error);
        }
    };
    return (
        <div className="flex h-full w-full justify-center items-center px-1">
            <div className="flex flex-col justify-center items-center w-[25rem] xl:w-[32rem] bg-white rounded-xl shadow-lg -mt-20">
                <div className="flex flex-col items-center justify-center w-full h-64">
                <div className="flex flex-col items-center justify-center w-full h-full">
                <div className="w-full h-full p-2">
                    <div 
                        className="flex flex-col w-full h-full justify-center items-center rounded-lg border-dashed border-2 border-gray-300 p-10 cursor-pointer"
                        onClick={downloadItems}
                        >
                                <FileDownloadIcon className="w-12 h-12"/>

                                <span className="font-linik  text-gray-800 text-sm mt-2">
                                    Click the icon to <span className="text-blue-500">download</span> files shared
                                </span>
                                <span className="font-linik  text-gray-400 text-xs ">
                                    Files will be downloaded in .zip format
                                </span>
                            </div>
                        </div>
                    </div>
                </div>
            </div>
        </div>
    );
};

export default DownloadPage;
