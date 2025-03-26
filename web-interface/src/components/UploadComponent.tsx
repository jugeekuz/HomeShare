import React from "react";
import { FileProvider } from "../contexts/FileContext.tsx";
import FileLoader from "../components/FileLoader.tsx";
import FileList from "../components/FileList.tsx";
import FileUploader from "../components/FileUploader.tsx";

const UploadComponent: React.FC = () => {
  return (
    <div className="flex flex-col w-full h-full">
        <FileProvider>
            <div className="flex flex-col items-center justify-center w-full h-[14rem]">
                <FileLoader />
            </div>
            <FileList />
            <div className="flex flex-col justify-center items-center w-full mt-1">
                <FileUploader />
            </div>
        </FileProvider>
    </div>
  );
};

export default UploadComponent;
