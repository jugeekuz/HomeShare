import React from "react";

import { FileUploadProvider } from "../contexts/FileUploadContext";
import FileLoader from "../components/FileLoader.tsx";
import FileUploadList from "../components/FileUploadList";
import FileUploader from "../components/FileUploader.tsx";

const UploadComponent: React.FC = () => {
  return (
    <div className="flex flex-col w-full h-full">
        <FileUploadProvider>
            <div className="flex flex-col items-center justify-center w-full h-[14rem]">
                <FileLoader />
            </div>
            <FileUploadList />
            <div className="flex flex-col justify-center items-center w-full mt-1">
                <FileUploader />
            </div>
        </FileUploadProvider>
    </div>
  );
};

export default UploadComponent;
